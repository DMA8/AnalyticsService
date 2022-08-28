package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"gitlab.com/g6834/team31/analytics/internal/config"
)

func (s *Server) analyticsHandlers(cfg *config.Config) http.Handler {
	h := chi.NewMux()
	h.Use(s.ValidateToken)
	h.Route("/", func(r chi.Router) {
		h.Get("/approved_tasks", s.approvedTasks)
		h.Get("/declined_tasks", s.declinedTasks)
		h.Get("/summary_time", s.summaryTime)
	})
	return h
}

// approvedTasks godoc
// @Summary get count of approved tasks
// @Description endpoint return count of approved tasks
// @Produce json
// @Success 200 {object} models.Counter
// @Router /approved_tasks [get]
func (s *Server) approvedTasks(w http.ResponseWriter, r *http.Request) {
	initHeaders(w)
	response, err := s.analyticsService.ApprovedTasks(r.Context())
	if err != nil {
		writeAnswer(w, http.StatusBadRequest, err.Error())
		s.logger.Debug().Msgf("s.approvedTasks err: %v", err)
		return
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		s.logger.Warn().Msg(err.Error())
	}
}

// declinedTasks godoc
// @Summary get count of declined tasks
// @Description endpoint return count of declined tasks
// @Produce json
// @Success 200 {object} models.Counter
// @Router /declined_tasks [get]
func (s *Server) declinedTasks(w http.ResponseWriter, r *http.Request) {
	initHeaders(w)
	response, err := s.analyticsService.DeclinedTasks(r.Context())
	if err != nil {
		writeAnswer(w, http.StatusBadRequest, err.Error())
		s.logger.Debug().Msgf("s.declinedTasks err: %v", err)
		return
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		s.logger.Warn().Msg(err.Error())
	}
}

// summaryTime godoc
// @Summary Get summary time for each task
// @Description Return task id and summary time of decision in seconds
// @Produce json
// @Success 200 {array} models.SummaryTime
// @Router /summary_time [get]
func (s *Server) summaryTime(w http.ResponseWriter, r *http.Request) {
	initHeaders(w)
	// TODO возможно нужна пагинация
	response, err := s.analyticsService.SummaryTime(r.Context())
	if err != nil {
		writeAnswer(w, http.StatusBadRequest, err.Error())
		s.logger.Debug().Msgf("s.summaryTime err: %v", err)
		return
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		s.logger.Warn().Err(err)
	}
}
