package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	rmodel "github.com/dormitory-life/core/internal/server/request_models"
	"github.com/gorilla/mux"
)

// @Summary Создать фото общежития
// @Description Создает новое фото для общежития
// @Tags Dormitories
// @Accept multipart/form-data
// @Produce json
// @Param dormitory_id path string true "ID общежития"
// @Param photos formData []file true "Фотографии общежития" collectionFormat(multi)
// @Success 201 {object} rmodel.CreateDormitoryEventResponse "Событие создано"
// @Failure 400 {object} rmodel.ErrorResponse "Некорректные данные формы или отсутствуют обязательные поля"
// @Failure 401 {object} rmodel.ErrorResponse "Пользователь не авторизован"
// @Failure 403 {object} rmodel.ErrorResponse "Нет прав на действие"
// @Failure 500 {object} rmodel.ErrorResponse "Внутренняя ошибка сервера"
// @Security BearerAuth
// @Router /core/dormitories/{dormitory_id}/photos [post]
func (s *Server) createDormitoryPhotosHandler(w http.ResponseWriter, r *http.Request) {
	const handlerName = "createDormitoryPhotosHandler"

	var (
		vars        = mux.Vars(r)
		dormitoryId = vars["dormitory_id"]
	)

	req := &rmodel.CreateDormitoryPhotosRequest{
		DormitoryId: dormitoryId,
	}

	err := r.ParseMultipartForm(50 << 20)
	if err != nil {
		writeErrorResponse(w, fmt.Errorf("failed to parse form: %v", err), http.StatusBadRequest)
		s.logger.Error("parse form error", slog.String("error", err.Error()))
		return
	}

	req.PhotoFilesHeaders = r.MultipartForm.File["photos"]
	if len(req.PhotoFilesHeaders) == 0 {
		writeErrorResponse(w, fmt.Errorf("no photos provided"), http.StatusBadRequest)
		return
	}

	s.logger.Info("uploading photos",
		slog.String("dormitory_id", dormitoryId),
		slog.Int("count", len(req.PhotoFilesHeaders)),
	)

	resp, err := s.coreService.CreateDormitoryPhotos(r.Context(), req)
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

// @Summary Удалить фотографии общежития
// @Description Удаляет фотографии общежития
// @Tags Dormitories
// @Param dormitory_id path string true "ID общежития"
// @Success 200 {object} rmodel.DeleteDormitoryPhotosResponse "Фото удалены"
// @Failure 400 {object} rmodel.ErrorResponse "Неверные данные / параметры запроса"
// @Failure 401 {object} rmodel.ErrorResponse "Пользователь не авторизован"
// @Failure 403 {object} rmodel.ErrorResponse "Нет прав на действие"
// @Failure 500 {object} rmodel.ErrorResponse "Внутренняя ошибка сервера"
// @Security BearerAuth
// @Router /core/dormitories/{dormitory_id}/photos [delete]
func (s *Server) deleteDormitoryPhotosHandler(w http.ResponseWriter, r *http.Request) {
	const handlerName = "deleteDormitoryPhotosHandler"

	var (
		vars        = mux.Vars(r)
		dormitoryId = vars["dormitory_id"]
	)

	req := &rmodel.DeleteDormitoryPhotosRequest{
		DormitoryId: dormitoryId,
	}

	resp, err := s.coreService.DeleteDormitoryPhotos(r.Context(), req)
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
