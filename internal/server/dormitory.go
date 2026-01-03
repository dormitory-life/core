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

func (s *Server) createDormitoryHandler(w http.ResponseWriter, r *http.Request) {
	const handlerName = "createDormitoryHandler"

	var req rmodel.CreateDormitoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorResponse(w, err, http.StatusBadRequest)
		s.logger.Error("error decoding request",
			slog.String("error", err.Error()),
			slog.String("handler", handlerName),
		)

		return
	}

	resp, err := s.coreService.CreateDormitory(r.Context(), &rmodel.CreateDormitoryRequest{
		DormitoryId:  req.DormitoryId,
		Name:         req.Name,
		Address:      req.Address,
		SupportEmail: req.SupportEmail,
		Description:  req.Description,
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

func (s *Server) updateDormitoryHandler(w http.ResponseWriter, r *http.Request) {
	const handlerName = "updateDormitoryHandler"

	var (
		vars        = mux.Vars(r)
		dormitoryId = vars["dormitory_id"]
	)

	var req rmodel.UpdateDormitoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorResponse(w, err, http.StatusBadRequest)
		s.logger.Error("error decoding request",
			slog.String("error", err.Error()),
			slog.String("handler", handlerName),
		)

		return
	}

	req.DormitoryId = dormitoryId

	resp, err := s.coreService.UpdateDormitory(r.Context(), &rmodel.UpdateDormitoryRequest{
		DormitoryId:  req.DormitoryId,
		Name:         req.Name,
		Address:      req.Address,
		SupportEmail: req.SupportEmail,
		Description:  req.Description,
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

func (s *Server) deleteDormitoryHandler(w http.ResponseWriter, r *http.Request) {
	const handlerName = "deleteDormitoryHandler"

	var (
		vars        = mux.Vars(r)
		dormitoryId = vars["dormitory_id"]
	)

	var req rmodel.DeleteDormitoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorResponse(w, err, http.StatusBadRequest)
		s.logger.Error("error decoding request",
			slog.String("error", err.Error()),
			slog.String("handler", handlerName),
		)

		return
	}

	req.DormitoryId = dormitoryId

	resp, err := s.coreService.DeleteDormitory(r.Context(), &rmodel.DeleteDormitoryRequest{
		DormitoryId: req.DormitoryId,
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
