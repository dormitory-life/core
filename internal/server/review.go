package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	rmodel "github.com/dormitory-life/core/internal/server/request_models"

	"github.com/gorilla/mux"
)

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

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		writeErrorResponse(w, err, http.StatusInternalServerError)
		s.logger.Error("error encoding response",
			slog.String("error", err.Error()),
			slog.String("handler", handlerName),
		)
	}
}

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

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		writeErrorResponse(w, err, http.StatusInternalServerError)
		s.logger.Error("error encoding response",
			slog.String("error", err.Error()),
			slog.String("handler", handlerName),
		)
	}
}

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
