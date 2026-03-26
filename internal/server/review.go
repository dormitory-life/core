package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	rmodel "github.com/dormitory-life/core/internal/server/request_models"

	"github.com/gorilla/mux"
)

// @Summary Получение отзывов
// @Description Получение отзывов общежития
// @Tags Reviews
// @Produce json
// @Success 200 {object} rmodel.GetDormitoryReviewsResponse "Отзывы"
// @Failure 400 {object} rmodel.ErrorResponse "Неверные данные / параметры запроса"
// @Failure 500 {object} rmodel.ErrorResponse "Внутренняя ошибка сервера"
// @Router /core/dormitories/{dormitory_id}/reviews [get]
func (s *Server) getReviewsHandler(w http.ResponseWriter, r *http.Request) {
	const handlerName = "getReviewsHandler"

	var (
		vars        = mux.Vars(r)
		dormitoryId = vars["dormitory_id"]
	)

	req, err := new(rmodel.GetDormitoryReviewsRequest).FromUrlQuery(r.URL.Query())
	if err != nil {
		writeErrorResponse(w, err, http.StatusBadRequest)
		s.logger.Error("error encoding response",
			slog.String("error", err.Error()),
			slog.String("handler", handlerName),
		)
		return
	}

	req.DormitoryId = dormitoryId

	resp, err := s.coreService.GetReviews(r.Context(), req)
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

// @Summary Создать отзыв
// @Description Создание отзыва на общежитие
// @Tags Reviews
// @Accept multipart/form-data
// @Produce json
// @Param dormitory_id path string true "ID общежития"
// @Param title formData string true "Заголовок отзыва"
// @Param description formData string true "Описание отзыва"
// @Param photos formData []file true "Фотографии отзыва" collectionFormat(multi)
// @Success 201 {object} rmodel.CreateReviewResponse "Отзыв создано"
// @Failure 400 {object} rmodel.ErrorResponse "Некорректные данные формы или отсутствуют обязательные поля"
// @Failure 401 {object} rmodel.ErrorResponse "Пользователь не авторизован"
// @Failure 403 {object} rmodel.ErrorResponse "Нет прав на действие"
// @Failure 500 {object} rmodel.ErrorResponse "Внутренняя ошибка сервера"
// @Security BearerAuth
// @Router /core/dormitories/{dormitory_id}/reviews [post]
func (s *Server) createReviewHandler(w http.ResponseWriter, r *http.Request) {
	const handlerName = "createReviewHandler"

	req, err := s.parseCreateReviewRequest(w, r)
	if err != nil {
		return
	}

	resp, err := s.coreService.CreateReview(r.Context(), req)
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

// @Summary Удалить отзыв
// @Description Удаляет отзыв на общежитие
// @Tags Reviews
// @Param dormitory_id path string true "ID общежития"
// @Success 204
// @Failure 400 {object} rmodel.ErrorResponse "Неверные данные / параметры запроса"
// @Failure 401 {object} rmodel.ErrorResponse "Пользователь не авторизован"
// @Failure 403 {object} rmodel.ErrorResponse "Нет прав на действие"
// @Failure 500 {object} rmodel.ErrorResponse "Внутренняя ошибка сервера"
// @Security BearerAuth
// @Router /core/dormitories/{dormitory_id}/reviews/{event_id} [delete]
func (s *Server) deleteReviewHandler(w http.ResponseWriter, r *http.Request) {
	const handlerName = "deleteReviewHandler"

	var (
		vars        = mux.Vars(r)
		dormitoryId = vars["dormitory_id"]
		reviewId    = vars["review_id"]
	)

	req := &rmodel.DeleteReviewRequest{
		DormitoryId: dormitoryId,
		ReviewId:    reviewId,
	}

	w.WriteHeader(http.StatusNoContent)

	_, err := s.coreService.DeleteReview(r.Context(), req)
	if err != nil {
		s.handleError(w, err)
		s.logger.Error("error",
			slog.String("error", err.Error()),
			slog.String("handler", handlerName),
		)

		return
	}
}

func (s *Server) parseCreateReviewRequest(w http.ResponseWriter, r *http.Request) (*rmodel.CreateReviewRequest, error) {
	var (
		vars        = mux.Vars(r)
		dormitoryId = vars["dormitory_id"]
	)

	req := &rmodel.CreateReviewRequest{
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
