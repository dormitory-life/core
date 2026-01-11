package server

import (
	"encoding/json"
	"log/slog"
	"net/http"

	rmodel "github.com/dormitory-life/core/internal/server/request_models"
)

func (s *Server) createSupportRequestHandler(w http.ResponseWriter, r *http.Request) {
	const handlerName = "getDormitoryByIdHandler"

	var req rmodel.CreateSupportRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorResponse(w, err, http.StatusBadRequest)
		s.logger.Error("error decoding request",
			slog.String("error", err.Error()),
			slog.String("handler", handlerName),
		)

		return
	}

	resp, err := s.coreService.CreateSupportRequest(r.Context(), &rmodel.CreateSupportRequest{
		Title:       req.Title,
		Description: req.Description,
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
