package server

import (
	"encoding/json"
	"log/slog"
	"net/http"

	rmodel "github.com/dormitory-life/core/internal/server/request_models"
	"github.com/gorilla/mux"
)

// @Summary Получение чата
// @Description Получение сообщений чата общежития
// @Tags Chat
// @Produce json
// @Params dormitory_id path string true "ID общежития"
// @Params page query int false "номер страницы"
// @Success 200 {object} rmodel.GetChatMessagesResponse "Оценки"
// @Failure 400 {object} rmodel.ErrorResponse "Неверные данные / параметры запроса"
// @Failure 500 {object} rmodel.ErrorResponse "Внутренняя ошибка сервера"
// @Router /core/dormitories/{dormitory_id}/chat [get]
func (s *Server) getDormitoryChatHandler(w http.ResponseWriter, r *http.Request) {
	const handlerName = "getDormitoryChatHandler"

	var (
		vars        = mux.Vars(r)
		dormitoryId = vars["dormitory_id"]
	)

	req, err := new(rmodel.GetChatMessagesRequest).FromUrlQuery(r.URL.Query())
	if err != nil {
		s.handleError(w, err)
		s.logger.Error("error",
			slog.String("error", err.Error()),
			slog.String("handler", handlerName),
		)

		return
	}

	req.DormitoryID = dormitoryId

	resp, err := s.coreService.GetChat(r.Context(), req)
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

// @Summary Создание cообщения
// @Description Создание сообщения в чате общежития
// @Tags Chat
// @Accept json
// @Produce json
// @Params dormitory_id path string true "ID общежития"
// @Params request body true rmodel.CreateChatMessageRequest "Информация о сообщении"
// @Success 200 {object} rmodel.CreateChatMessageResponse "Сообщение создано"
// @Failure 400 {object} rmodel.ErrorResponse "Неверные данные / параметры запроса"
// @Failure 403 {object} rmodel.ErrorResponse "Нет прав на действие"
// @Failure 500 {object} rmodel.ErrorResponse "Внутренняя ошибка сервера"
// @Router /core/dormitories/{dormitory_id}/chat [post]
func (s *Server) createChatMessageHandler(w http.ResponseWriter, r *http.Request) {
	const handlerName = "createChatMessageHandler"

	var (
		vars        = mux.Vars(r)
		dormitoryId = vars["dormitory_id"]
	)

	var req rmodel.CreateChatMessageRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorResponse(w, err, http.StatusBadRequest)
		s.logger.Error("error decoding request",
			slog.String("error", err.Error()),
			slog.String("handler", handlerName),
		)

		return
	}

	req.DormitoryID = dormitoryId

	resp, err := s.coreService.CreateChatMessage(r.Context(), &req)
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
