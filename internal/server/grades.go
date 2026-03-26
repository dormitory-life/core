package server

import (
	"encoding/json"
	"log/slog"
	"net/http"

	rmodel "github.com/dormitory-life/core/internal/server/request_models"
	"github.com/gorilla/mux"
)

// @Summary Получение средних оценок общежитий
// @Description Получение средних оценок общежитий
// @Tags Grades
// @Produce json
// @Success 200 {object} rmodel.GetDormitoriesAvgGradesResponse "Оценки"
// @Failure 400 {object} rmodel.ErrorResponse "Неверные данные / параметры запроса"
// @Failure 500 {object} rmodel.ErrorResponse "Внутренняя ошибка сервера"
// @Router /core/dormitories/grades [get]
func (s *Server) getDormitoriesAvgGradesHandler(w http.ResponseWriter, r *http.Request) {
	const handlerName = "getDormitoriesAvgGradesHandler"

	resp, err := s.coreService.GetDormitoriesAvgGrades(r.Context(), &rmodel.GetDormitoriesAvgGradesRequest{})
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

// @Summary Получение средних оценок общежития
// @Description Получение средних оценок общежития
// @Tags Grades
// @Produce json
// @Params dormitory_id path string true "ID общежития"
// @Success 200 {object} rmodel.GetDormitoryAvgGradesResponse "Оценки"
// @Failure 400 {object} rmodel.ErrorResponse "Неверные данные / параметры запроса"
// @Failure 500 {object} rmodel.ErrorResponse "Внутренняя ошибка сервера"
// @Router /core/dormitories/{dormitory_id}/grades [get]
func (s *Server) getDormitoryAvgGradesHandler(w http.ResponseWriter, r *http.Request) {
	const handlerName = "getDormitoryAvgGradesHandler"

	var (
		vars        = mux.Vars(r)
		dormitoryId = vars["dormitory_id"]
	)

	resp, err := s.coreService.GetDormitoryAvgGrades(r.Context(), &rmodel.GetDormitoryAvgGradesRequest{
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

// @Summary Создание оценки общежития
// @Description Создание оценки общежития по критериям
// @Tags Grades
// @Accept json
// @Produce json
// @Params dormitory_id path string true "ID общежития"
// @Params request body true rmodel.CreateDormitoryGradeRequest "Оценки по критериям"
// @Success 200 {object} rmodel.CreateDormitoryGradeResponse "Оценка создана"
// @Failure 400 {object} rmodel.ErrorResponse "Неверные данные / параметры запроса"
// @Failure 403 {object} rmodel.ErrorResponse "Нет прав на оценивание"
// @Failure 409 {object} rmodel.ErrorResponse "Уже оценивали в этом месяце"
// @Failure 500 {object} rmodel.ErrorResponse "Внутренняя ошибка сервера"
// @Router /core/dormitories/{dormitory_id}/grades [post]
func (s *Server) createDormitoryGradeHandler(w http.ResponseWriter, r *http.Request) {
	const handlerName = "createDormitoryGradeHandler"

	var (
		vars        = mux.Vars(r)
		dormitoryId = vars["dormitory_id"]
	)

	var req rmodel.CreateDormitoryGradeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorResponse(w, err, http.StatusBadRequest)
		s.logger.Error("error decoding request",
			slog.String("error", err.Error()),
			slog.String("handler", handlerName),
		)

		return
	}

	req.DormitoryId = dormitoryId

	resp, err := s.coreService.CreateDormitoryGrade(r.Context(), &req)
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
