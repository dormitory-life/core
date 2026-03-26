package server

import (
	"encoding/json"
	"log/slog"
	"net/http"

	rmodel "github.com/dormitory-life/core/internal/server/request_models"
	"github.com/gorilla/mux"
)

// @Summary Получение списка общежитий
// @Description Получение списка краткой информации о всех общежитиях
// @Tags Dormitories
// @Produce json
// @Success 200 {object} rmodel.GetDormitoriesResponse "Общежития"
// @Failure 400 {object} rmodel.ErrorResponse "Неверные данные / параметры запроса"
// @Failure 403 {object} rmodel.ErrorResponse "Нет прав на действие"
// @Failure 500 {object} rmodel.ErrorResponse "Внутренняя ошибка сервера"
// @Router /core/dormitories [get]
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

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		writeErrorResponse(w, err, http.StatusInternalServerError)
		s.logger.Error("error encoding response",
			slog.String("error", err.Error()),
			slog.String("handler", handlerName),
		)
	}
}

// @Summary Получение общежития
// @Description Получение подробной информации об общежитии
// @Tags Dormitories
// @Produce json
// @Params dormitory_id path string true "ID общежития"
// @Success 200 {object} rmodel.GetDormitoryByIdResponse "Общежитие"
// @Failure 400 {object} rmodel.ErrorResponse "Неверные данные / параметры запроса"
// @Failure 403 {object} rmodel.ErrorResponse "Нет прав на действие"
// @Failure 500 {object} rmodel.ErrorResponse "Внутренняя ошибка сервера"
// @Router /core/dormitories/{dormitory_id} [get]
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

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		writeErrorResponse(w, err, http.StatusInternalServerError)
		s.logger.Error("error encoding response",
			slog.String("error", err.Error()),
			slog.String("handler", handlerName),
		)
	}
}

// @Summary Создание общежития
// @Description Создает общежитие в системе
// @Tags Dormitories
// @Accept json
// @Produce json
// @Params dormitory_id path string true "ID общежития"
// @Params request body rmodel.CreateDormitoryRequest true "Информация об общежитии"
// @Success 200 {object} rmodel.CreateDormitoryResponse "Общежитие создано"
// @Failure 400 {object} rmodel.ErrorResponse "Неверные данные / параметры запроса"
// @Failure 403 {object} rmodel.ErrorResponse "Нет прав на действие"
// @Failure 500 {object} rmodel.ErrorResponse "Внутренняя ошибка сервера"
// @Router /core/dormitories [post]
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

	resp, err := s.coreService.CreateDormitory(r.Context(), &req)
	if err != nil {
		s.handleError(w, err)
		s.logger.Error("error",
			slog.String("error", err.Error()),
			slog.String("handler", handlerName),
		)

		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		writeErrorResponse(w, err, http.StatusInternalServerError)
		s.logger.Error("error encoding response",
			slog.String("error", err.Error()),
			slog.String("handler", handlerName),
		)
	}
}

// @Summary Обновление информации об общежитии
// @Description Обновляет информацию об общежитии
// @Tags Dormitories
// @Accept json
// @Produce json
// @Params dormitory_id path string true "ID общежития"
// @Params request body rmodel.UpdateDormitoryRequest true "Информация для обновления"
// @Success 200 {object} rmodel.UpdateDormitoryResponse "Общежитие обновлено"
// @Failure 400 {object} rmodel.ErrorResponse "Неверные данные / параметры запроса"
// @Failure 403 {object} rmodel.ErrorResponse "Нет прав на действие"
// @Failure 500 {object} rmodel.ErrorResponse "Внутренняя ошибка сервера"
// @Router /core/dormitories/{dormitory_id} [put]
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

	resp, err := s.coreService.UpdateDormitory(r.Context(), &req)
	if err != nil {
		s.handleError(w, err)
		s.logger.Error("error",
			slog.String("error", err.Error()),
			slog.String("handler", handlerName),
		)

		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		writeErrorResponse(w, err, http.StatusInternalServerError)
		s.logger.Error("error encoding response",
			slog.String("error", err.Error()),
			slog.String("handler", handlerName),
		)
	}
}

// @Summary Удаление общежития
// @Description Удаляет общежитие из системы
// @Tags Dormitories
// @Produce json
// @Params dormitory_id path string true "ID общежития"
// @Success 200 {object} rmodel.DeleteDormitoryResponse "Общежитие удалено"
// @Failure 400 {object} rmodel.ErrorResponse "Неверные данные / параметры запроса"
// @Failure 403 {object} rmodel.ErrorResponse "Нет прав на действие"
// @Failure 500 {object} rmodel.ErrorResponse "Внутренняя ошибка сервера"
// @Router /core/dormitories/{dormitory_id} [delete]
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

	resp, err := s.coreService.DeleteDormitory(r.Context(), &req)
	if err != nil {
		s.handleError(w, err)
		s.logger.Error("error",
			slog.String("error", err.Error()),
			slog.String("handler", handlerName),
		)

		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		writeErrorResponse(w, err, http.StatusInternalServerError)
		s.logger.Error("error encoding response",
			slog.String("error", err.Error()),
			slog.String("handler", handlerName),
		)
	}
}
