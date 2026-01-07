package core

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/dormitory-life/core/internal/constants"
	rmodel "github.com/dormitory-life/core/internal/server/request_models"
	"github.com/dormitory-life/core/internal/storage"
	"github.com/google/uuid"
)

func (s *CoreService) CreateDormitoryPhotos(
	ctx context.Context,
	request *rmodel.CreateDormitoryPhotosRequest,
) (*rmodel.CreateDormitoryPhotosResponse, error) {
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
			RoleRequired: true,
		},
	); err != nil {
		return nil, err
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
			Category: constants.CategoryDormitoryPhotos,
			EntityId: dormitoryId,
			PhotoId:  fileId,
			FileName: photoFileHeader.Filename,
			Reader:   file,
			Size:     photoFileHeader.Size,
			MimeType: s.s3Client.GetMimeType(photoFileHeader.Filename),
		})
		if err != nil {
			for _, uploaded := range uploadedPhotos {
				s.s3Client.Delete(ctx, uploaded.FilePath)
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

	return &rmodel.CreateDormitoryPhotosResponse{
		CreatePhotoResponses: uploadedPhotos,
	}, nil
}

func (s *CoreService) DeleteDormitoryPhotos(
	ctx context.Context,
	request *rmodel.DeleteDormitoryPhotosRequest,
) (*rmodel.DeleteDormitoryPhotosResponse, error) {
	err := s.s3Client.DeleteAll(ctx, constants.CategoryDormitoryPhotos, request.DormitoryId)
	if err != nil {
		return nil, fmt.Errorf("%w: error deleting dormitory photos: %v", ErrInternal, err)
	}

	return &rmodel.DeleteDormitoryPhotosResponse{
		DormitoryId: request.DormitoryId,
	}, nil
}
