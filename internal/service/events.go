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

func (s *CoreService) GetDormitoryEvents(
	ctx context.Context,
	request *rmodel.GetDormitoryEventsRequest,
) (*rmodel.GetDormitoryEventsResponse, error) {
	if request == nil {
		return nil, fmt.Errorf("%w: request is nil", ErrBadRequest)
	}

	resp, err := s.repository.GetDormitoryEvents(ctx, &dbtypes.GetDormitoryEventsRequest{
		DormitoryId: request.DormitoryId,
		Page:        request.Page,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: error getting events: %v", s.handleDBError(err), err)
	}

	res := new(rmodel.GetDormitoryEventsResponse).From(resp)

	for i, event := range res.Events {
		eventPhotos, err := s.s3Client.GetEntityFiles(ctx, &storage.GetEntityFilesRequest{
			Category:    constants.CategoryEventPhotos,
			EntityId:    request.DormitoryId,
			SubEntityId: event.EventId,
		})
		if err != nil {
			s.logger.Warn("error getting dormitory event photos",
				slog.String("error", err.Error()),
				slog.String("dormId", request.DormitoryId),
				slog.String("eventId", event.EventId))
		}
		res.Events[i].EventPhotos = rmodel.ConvertFileInfos(eventPhotos)
	}

	return res, nil
}

func (s *CoreService) CreateDormitoryEvent(
	ctx context.Context,
	request *rmodel.CreateDormitoryEventRequest,
) (*rmodel.CreateDormitoryEventResponse, error) {
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

	createResp, err := s.repository.CreateDormitoryEvent(ctx, &dbtypes.CreateDormitoryEventRequest{
		DormitoryId: dormitoryId,
		Title:       request.Title,
		Description: request.Description,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: error creating event: %v", s.handleDBError(err), err)
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
			Category:    constants.CategoryEventPhotos,
			EntityId:    dormitoryId,
			SubEntityId: createResp.EventId,
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

			_, deleteErr := s.repository.DeleteDormitoryEvent(ctx, &dbtypes.DeleteDormitoryEventRequest{EventId: createResp.EventId})
			if deleteErr != nil {
				return nil, fmt.Errorf("%w: error deleting event photos after upload error: %v", ErrInternal, errors.Join(err, deleteErr))
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

	return &rmodel.CreateDormitoryEventResponse{
		EventId:              createResp.EventId,
		DormitoryId:          dormitoryId,
		CreatePhotoResponses: uploadedPhotos,
		Title:                request.Title,
		Description:          request.Description,
	}, nil
}

func (s *CoreService) DeleteDormitoryEvent(
	ctx context.Context,
	request *rmodel.DeleteDormitoryEventRequest,
) (*rmodel.DeleteDormitoryEventResponse, error) {
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

	_, err = s.repository.DeleteDormitoryEvent(ctx, &dbtypes.DeleteDormitoryEventRequest{
		EventId: request.EventId,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: error deleting event: %v", s.handleDBError(err), err)
	}

	if err := s.s3Client.DeleteAll(ctx, &storage.DeleteAllRequest{
		Category:    constants.CategoryEventPhotos,
		EntityId:    dormitoryId,
		SubEntityId: request.EventId,
	}); err != nil {
		return nil, fmt.Errorf("%w: error deleting event photos: %v", ErrInternal, err)
	}

	return &rmodel.DeleteDormitoryEventResponse{}, nil
}
