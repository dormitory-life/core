package server

import (
	"encoding/json"
	"log/slog"
	"net/http"

	rmodel "github.com/dormitory-life/core/internal/server/request_models"
	"github.com/gorilla/mux"
)

func (s *Server) getDormitoriesHandler(w http.ResponseWriter, r *http.Request) {
	const handlerName = "getDormitoriesHandler"

	resp, err := s.coreService.GetDormitories(r.Context(), &rmodel.GetDormitoriesRequest{})
	if err != nil {
		s.handleError(w, err)
		s.logger.Error("error",
			slog.String("error", err.Error()),
			slog.String("handler", handlerName),
		)

		return
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		writeErrorResponse(w, err, http.StatusInternalServerError)
		s.logger.Error("error encoding response",
			slog.String("error", err.Error()),
			slog.String("handler", handlerName),
		)
	}
}

func (s *Server) getDormitoryByIdHandler(w http.ResponseWriter, r *http.Request) {
	const handlerName = "getDormitoryByIdHandler"

	var (
		vars        = mux.Vars(r)
		dormitoryId = vars["dormitory_id"]
	)

	resp, err := s.coreService.GetDormitoryById(r.Context(), &rmodel.GetDormitoryByIdRequest{
		DormitoryId: dormitoryId,
	})
	if err != nil {
		s.handleError(w, err)
		s.logger.Error("error",
			slog.String("error", err.Error()),
			slog.String("handler", handlerName),
		)

		return
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		writeErrorResponse(w, err, http.StatusInternalServerError)
		s.logger.Error("error encoding response",
			slog.String("error", err.Error()),
			slog.String("handler", handlerName),
		)
	}
}
