package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	rmodel "github.com/dormitory-life/core/internal/server/request_models"
	"github.com/gorilla/mux"
)

// @Summary Получение ленты
// @Description Получение списка событий общежития
// @Tags Feed
// @Produce json
// @Params dormitory_id path string true "ID общежития"
// @Params page query int false "номер страницы"
// @Success 200 {object} rmodel.GetDormitoryEventsResponse "Лента"
// @Failure 400 {object} rmodel.ErrorResponse "Неверные данные / параметры запроса"
// @Failure 500 {object} rmodel.ErrorResponse "Внутренняя ошибка сервера"
// @Router /core/dormitories/{dormitory_id}/events [get]
func (s *Server) getDormitoryEventsHandler(w http.ResponseWriter, r *http.Request) {
	const handlerName = "getDormitoryEventsHandler"

	var (
		vars        = mux.Vars(r)
		dormitoryId = vars["dormitory_id"]
	)

	req, err := new(rmodel.GetDormitoryEventsRequest).FromUrlQuery(r.URL.Query())
	if err != nil {
		s.handleError(w, err)
		s.logger.Error("error",
			slog.String("error", err.Error()),
			slog.String("handler", handlerName),
		)

		return
	}

	req.DormitoryId = dormitoryId

	resp, err := s.coreService.GetDormitoryEvents(r.Context(), req)
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

// @Summary Создать событие общежития
// @Description Создает новое событие для общежития с заголовком, описанием и фотографиями
// @Tags Feed
// @Accept multipart/form-data
// @Produce json
// @Param dormitory_id path string true "ID общежития"
// @Param title formData string true "Заголовок события"
// @Param description formData string true "Описание события"
// @Param photos formData file true "Фотографии события"
// @Success 201 {object} rmodel.CreateDormitoryEventResponse "Событие создано"
// @Failure 400 {object} rmodel.ErrorResponse "Некорректные данные формы или отсутствуют обязательные поля"
// @Failure 401 {object} rmodel.ErrorResponse "Пользователь не авторизован"
// @Failure 403 {object} rmodel.ErrorResponse "Нет прав на действие"
// @Failure 500 {object} rmodel.ErrorResponse "Внутренняя ошибка сервера"
// @Security BearerAuth
// @Router /core/dormitories/{dormitory_id}/events [post]
func (s *Server) createDormitoryEventHandler(w http.ResponseWriter, r *http.Request) {
	const handlerName = "createDormitoryEventHandler"

	req, err := s.parseCreateEventRequest(w, r)
	if err != nil {
		return
	}

	resp, err := s.coreService.CreateDormitoryEvent(r.Context(), req)
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

// @Summary Удалить событие общежития
// @Description Удаляет событие из ленты общежития
// @Tags Feed
// @Param dormitory_id path string true "ID общежития"
// @Success 204
// @Failure 400 {object} rmodel.ErrorResponse "Неверные данные / параметры запроса"
// @Failure 401 {object} rmodel.ErrorResponse "Пользователь не авторизован"
// @Failure 403 {object} rmodel.ErrorResponse "Нет прав на действие"
// @Failure 500 {object} rmodel.ErrorResponse "Внутренняя ошибка сервера"
// @Security BearerAuth
// @Router /core/dormitories/{dormitory_id}/events/{event_id} [delete]
func (s *Server) deleteDormitoryEventHandler(w http.ResponseWriter, r *http.Request) {
	const handlerName = "deleteDormitoryEventHandler"

	var (
		vars        = mux.Vars(r)
		dormitoryId = vars["dormitory_id"]
		eventId     = vars["event_id"]
	)

	req := &rmodel.DeleteDormitoryEventRequest{
		DormitoryId: dormitoryId,
		EventId:     eventId,
	}

	w.WriteHeader(http.StatusNoContent)

	_, err := s.coreService.DeleteDormitoryEvent(r.Context(), req)
	if err != nil {
		s.handleError(w, err)
		s.logger.Error("error",
			slog.String("error", err.Error()),
			slog.String("handler", handlerName),
		)

		return
	}
}

func (s *Server) parseCreateEventRequest(w http.ResponseWriter, r *http.Request) (*rmodel.CreateDormitoryEventRequest, error) {
	var (
		vars        = mux.Vars(r)
		dormitoryId = vars["dormitory_id"]
	)

	req := &rmodel.CreateDormitoryEventRequest{
		DormitoryId: dormitoryId,
	}

	err := r.ParseMultipartForm(50 << 20)
	if err != nil {
		writeErrorResponse(w, fmt.Errorf("failed to parse form: %v", err), http.StatusBadRequest)
		s.logger.Error("parse form error", slog.String("error", err.Error()))
		return nil, err
	}

	req.PhotoFilesHeaders = r.MultipartForm.File["photos"]
	if len(req.PhotoFilesHeaders) == 0 {
		writeErrorResponse(w, fmt.Errorf("no photos provided"), http.StatusBadRequest)
		return nil, err
	}

	s.logger.Info("uploading review photos",
		slog.String("dormitory_id", dormitoryId),
		slog.Int("count", len(req.PhotoFilesHeaders)),
	)

	req.Title = r.FormValue("title")
	if len(req.Title) == 0 {
		writeErrorResponse(w, fmt.Errorf("no title"), http.StatusBadRequest)
		return nil, err
	}

	req.Description = r.FormValue("description")
	if len(req.Title) == 0 {
		writeErrorResponse(w, fmt.Errorf("no title"), http.StatusBadRequest)
		return nil, err
	}

	return req, nil
}
