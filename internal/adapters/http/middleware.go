package http

import (
	"gitlab.com/g6834/team31/analytics/internal/domain/errors"
	"gitlab.com/g6834/team31/analytics/internal/domain/models"
	"gitlab.com/g6834/team31/auth/pkg/logging"

	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
)

type ctxKey int

const ridKey ctxKey = ctxKey(0)

type userLogin string

// Logger .
func Logger(l *logging.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(rw http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(rw, r.ProtoMajor)
			start := time.Now()
			defer func() {
				// var logEntry map[string]interface{}
				Entry := &logging.Entry{
					Service:      "analytics",
					Method:       r.Method,
					Url:          r.URL.Path,
					Query:        r.URL.RawQuery,
					RemoteIP:     r.RemoteAddr,
					Status:       ww.Status(),
					Size:         ww.BytesWritten(),
					ReceivedTime: start,
					Duration:     time.Since(start),
					ServerIP:     r.Host,
					UserAgent:    r.Header.Get("User-Agent"),
					//RequestId:    GetReqID(r.Context()),
				}
				l.Info().Msgf("%+v", Entry)
			}()
			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}

func GetReqID(ctx context.Context) string {
	return ctx.Value(ridKey).(string)
}

func (s *Server) ValidateToken(next http.Handler) http.Handler {
	fn := func(rw http.ResponseWriter, r *http.Request) {
		access, err := r.Cookie("accessToken")
		if err != nil {
			s.logger.Debug().Msgf("middleware ValidateToken COULDN'T validate accessToken %v", err)
			writeAnswer(rw, http.StatusForbidden, errors.ErrCookie.Error())
			return
		}
		refresh, err := r.Cookie("refreshToken")
		if err != nil {
			s.logger.Debug().Err(err).Msgf("middleware ValidateToken COULDN'T validate refreshToken")
			writeAnswer(rw, http.StatusForbidden, errors.ErrCookie.Error())
			return
		}
		ctx := r.Context()
		credential, err := s.AuthClient.Validate(ctx, models.JWTTokens{
			Access:  access.Value,
			Refresh: refresh.Value,
		})
		if err != nil {
			s.logger.Debug().Err(err).Msgf("s.Server middleware ValidateToken s.AuthClient.Validate error")
			writeAnswer(rw, http.StatusForbidden, errors.ErrBadCredential.Error())
			return
		}
		ctx = context.WithValue(ctx, userLogin("userLogin"), credential.Login)
		if credential.IsUpdate {
			rw.Header().Add("Set-Cookie", "accessToken="+credential.AccessToken)
			rw.Header().Add("Set-Cookie", "refreshToken="+credential.RefreshToken)
		}
		next.ServeHTTP(rw, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
