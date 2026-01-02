package server

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	rmodel "github.com/dormitory-life/core/internal/server/request_models"
	core "github.com/dormitory-life/core/internal/service"

	"github.com/dormitory-life/core/internal/constants"
)

func (s *Server) pingHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

func writeErrorResponse(w http.ResponseWriter, err error, code int, details ...string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	response := rmodel.ErrorResponse{
		Error:   err.Error(),
		Details: details,
	}

	_ = json.NewEncoder(w).Encode(response)
}

func (s *Server) handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, core.ErrBadRequest):
		writeErrorResponse(w, constants.ErrBadRequest, http.StatusBadRequest, err.Error())
	case errors.Is(err, core.ErrNotFound):
		writeErrorResponse(w, constants.ErrNotFound, http.StatusNotFound)
	case errors.Is(err, core.ErrConflict):
		writeErrorResponse(w, constants.ErrConflict, http.StatusConflict, err.Error())
	case errors.Is(err, core.ErrUnauthorized):
		writeErrorResponse(w, constants.ErrUnauthorized, http.StatusUnauthorized)
	case errors.Is(err, core.ErrInternal):
		writeErrorResponse(w, constants.ErrInternalServerError, http.StatusInternalServerError)
	default:
		s.logger.Error("Unhandled core error", slog.String("error", err.Error()))
		writeErrorResponse(w, constants.ErrInternalServerError, http.StatusInternalServerError)
	}
}
