package core

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/dormitory-life/core/internal/constants"
	dbtypes "github.com/dormitory-life/core/internal/database/types"

	rmodel "github.com/dormitory-life/core/internal/server/request_models"
	"github.com/dormitory-life/core/internal/storage"
	"github.com/google/uuid"
)

func (s *CoreService) GetReviews(
	ctx context.Context,
	request *rmodel.GetDormitoryReviewsRequest,
) (*rmodel.GetDormitoryReviewsResponse, error) {
	if request == nil {
		return nil, fmt.Errorf("%w: request is nil", ErrBadRequest)
	}

	resp, err := s.repository.GetReviews(ctx, &dbtypes.GetDormitoryReviewsRequest{
		DormitoryId: request.DormitoryId,
		Page:        request.Page,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: error getting reviews: %v", s.handleDBError(err), err)
	}

	res := new(rmodel.GetDormitoryReviewsResponse).From(resp)

	for i, review := range res.Reviews {
		reviewPhotos, err := s.s3Client.GetEntityFiles(ctx, &storage.GetEntityFilesRequest{
			Category:    constants.CategoryReviewPhotos,
			EntityId:    request.DormitoryId,
			SubEntityId: review.ReviewId,
		})
		if err != nil {
			s.logger.Warn("error getting dormitory review photos",
				slog.String("error", err.Error()),
				slog.String("dormId", request.DormitoryId),
				slog.String("reviewId", review.ReviewId))
		}
		res.Reviews[i].ReviewPhotos = rmodel.ConvertFileInfos(reviewPhotos)
	}

	return res, nil
}

func (s *CoreService) CreateReview(
	ctx context.Context,
	request *rmodel.CreateReviewRequest,
) (*rmodel.CreateReviewResponse, error) {
	if request == nil {
		return nil, fmt.Errorf("%w: request is nil", ErrBadRequest)
	}

	userId, dormitoryId, err := s.extractIdsFromRequestContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: error getting ids from context: %v", ErrInternal, err)
	}

	if err := s.checkAccess(
		ctx,
		&rmodel.CheckAccessRequest{
			UserId:       userId,
			DormitoryId:  request.DormitoryId,
			RoleRequired: false,
		},
	); err != nil {
		return nil, err
	}

	createResp, err := s.repository.CreateReview(ctx, &dbtypes.CreateReviewRequest{
		OwnerId:     userId,
		DormitoryId: dormitoryId,
		Title:       request.Title,
		Description: request.Description,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: error creating review: %v", s.handleDBError(err), err)
	}

	var uploadedPhotos []rmodel.CreatePhotoResponse

	for _, photoFileHeader := range request.PhotoFilesHeaders {
		file, err := photoFileHeader.Open()
		if err != nil {
			s.logger.Warn("failed to open file",
				slog.String("filename", photoFileHeader.Filename),
				slog.String("error", err.Error()),
			)
			continue
		}
		defer file.Close()

		fileId := uuid.New().String()

		uploadResult, err := s.s3Client.Upload(ctx, &storage.UploadRequest{
			Category:    constants.CategoryReviewPhotos,
			EntityId:    dormitoryId,
			SubEntityId: createResp.ReviewId,
			PhotoId:     fileId,
			FileName:    photoFileHeader.Filename,
			Reader:      file,
			Size:        photoFileHeader.Size,
			MimeType:    s.s3Client.GetMimeType(photoFileHeader.Filename),
		})
		if err != nil {
			for _, uploaded := range uploadedPhotos {
				s.s3Client.Delete(ctx, &storage.DeleteFileRequest{
					Path: &uploaded.FilePath,
				})
			}

			_, deleteErr := s.repository.DeleteReview(ctx, &dbtypes.DeleteReviewRequest{ReviewId: createResp.ReviewId})
			if deleteErr != nil {
				return nil, fmt.Errorf("%w: error deleting review photos after upload error: %v", ErrInternal, errors.Join(err, deleteErr))
			}

			return nil, fmt.Errorf("%w: upload failed: %v", ErrInternal, err)
		}

		uploadedPhotos = append(uploadedPhotos, rmodel.CreatePhotoResponse{
			URL:      uploadResult.URL,
			FilePath: uploadResult.FilePath,
			FileName: photoFileHeader.Filename,
			Size:     uploadResult.Size,
		})
	}

	return &rmodel.CreateReviewResponse{
		ReviewId:             createResp.ReviewId,
		OwnerId:              userId,
		DormitoryId:          dormitoryId,
		CreatePhotoResponses: uploadedPhotos,
		Title:                request.Title,
		Description:          request.Description,
	}, nil
}

func (s *CoreService) DeleteReview(
	ctx context.Context,
	request *rmodel.DeleteReviewRequest,
) (*rmodel.DeleteReviewResponse, error) {
	if request == nil {
		return nil, fmt.Errorf("%w: request is nil", ErrBadRequest)
	}

	userId, dormitoryId, err := s.extractIdsFromRequestContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: error getting ids from context: %v", ErrInternal, err)
	}

	if err := s.checkAccess(
		ctx,
		&rmodel.CheckAccessRequest{
			UserId:       userId,
			DormitoryId:  request.DormitoryId,
			RoleRequired: false,
		},
	); err != nil {
		return nil, err
	}

	_, err = s.repository.DeleteReview(ctx, &dbtypes.DeleteReviewRequest{
		ReviewId: request.ReviewId,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: error deleting review: %v", s.handleDBError(err), err)
	}

	if err := s.s3Client.DeleteAll(ctx, &storage.DeleteAllRequest{
		Category:    constants.CategoryReviewPhotos,
		EntityId:    dormitoryId,
		SubEntityId: request.ReviewId,
	}); err != nil {
		return nil, fmt.Errorf("%w: error deleting review photos: %v", ErrInternal, err)
	}

	return &rmodel.DeleteReviewResponse{}, nil
}
