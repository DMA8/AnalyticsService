package http

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	"gitlab.com/g6834/team31/analytics/internal/config"
	_ "gitlab.com/g6834/team31/analytics/docs"
	"gitlab.com/g6834/team31/analytics/internal/ports"
	"gitlab.com/g6834/team31/auth/pkg/logging"
)

type Server struct {
	analyticsService ports.Analytics
	server           *http.Server
	Addr             string
	AuthClient       ports.ClientAuth
	logger           *logging.Logger
	cfg              *config.Config
}

func New(analytics ports.Analytics, client ports.ClientAuth, cfg *config.Config, logger *logging.Logger) *Server {
	var s Server
	s.analyticsService = analytics
	s.AuthClient = client
	s.cfg = cfg
	s.logger = logger
	server := &http.Server{
		Handler: s.routes(cfg),
		Addr:    cfg.HTTP.Port,
	}
	s.server = server
	return &s
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

// @title Swagger Analytics API
// @version 1.0
// @description This is analytics server.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host petstore.swagger.io
// @BasePath /analytics/v1
func (s *Server) routes(cfg *config.Config) http.Handler {
	r := chi.NewMux()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(Logger(s.logger))
	// TODO добавить урлы свагера через конфиг
	r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(s.cfg.HTTP.Port+"/swagger/doc.json")))
	r.Get("/healthz", s.healthzHandler)
	r.Mount(cfg.HTTP.ApiVersion, s.analyticsHandlers(cfg))

	return r
}

func (s *Server) healthzHandler(w http.ResponseWriter, r *http.Request) {
	writeAnswer(w, http.StatusOK, "OK")
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *Server) Start(ctx context.Context) chan error {
	chanErr := make(chan error)
	go func(){
		chanErr <- s.server.ListenAndServe()
	}()
	return chanErr 
}
