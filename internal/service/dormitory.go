package core

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/dormitory-life/core/internal/cache"
	"github.com/dormitory-life/core/internal/constants"
	dbtypes "github.com/dormitory-life/core/internal/database/types"
	rmodel "github.com/dormitory-life/core/internal/server/request_models"
	"github.com/dormitory-life/core/internal/storage"
)

func (s *CoreService) GetDormitories(
	ctx context.Context,
	request *rmodel.GetDormitoriesRequest,
) (*rmodel.GetDormitoriesResponse, error) {
	if request == nil {
		return nil, fmt.Errorf("%w: request is nil", ErrBadRequest)
	}

	if res, err := s.getDormitoriesFromCache(ctx); err == nil {
		return res, nil
	}

	s.logger.Debug("dormitories cache miss")

	resp, err := s.repository.GetDormitories(ctx, &dbtypes.GetDormitoriesRequest{})
	if err != nil {
		return nil, fmt.Errorf("%w: error getting dormitories: %v", s.handleDBError(err), err)
	}

	res := new(rmodel.GetDormitoriesResponse).From(resp)

	for i, dorm := range res.Dormitories {
		photos, err := s.s3Client.GetEntityFiles(ctx, &storage.GetEntityFilesRequest{
			Category: constants.CategoryDormitoryPhotos,
			EntityId: dorm.Id,
			Amount:   constants.GetDormitoriesDefaultAmount,
		})
		if err != nil {
			s.logger.Warn("error getting dormitory photos", slog.String("error", err.Error()), slog.String("dormId", dorm.Id))
		}

		dormPhotos := rmodel.ConvertFileInfos(photos)

		res.Dormitories[i].Photos = dormPhotos
	}

	s.setDormitoriesToCache(ctx, res)

	return res, nil
}

func (s *CoreService) GetDormitoryById(
	ctx context.Context,
	request *rmodel.GetDormitoryByIdRequest,
) (*rmodel.GetDormitoryByIdResponse, error) {
	if request == nil {
		return nil, fmt.Errorf("%w: request is nil", ErrBadRequest)
	}

	if res, err := s.getDormitoryByIdFromCache(ctx, request.DormitoryId); err == nil {
		return res, nil
	}

	s.logger.Debug("dormitory cache miss", slog.String("dormitoryId", request.DormitoryId))

	resp, err := s.repository.GetDormitoryById(ctx, &dbtypes.GetDormitoryByIdRequest{
		DormitoryId: request.DormitoryId,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: error getting dormitory: %v", s.handleDBError(err), err)
	}

	res := new(rmodel.GetDormitoryByIdResponse).From(resp)

	photos, err := s.s3Client.GetEntityFiles(ctx, &storage.GetEntityFilesRequest{
		Category: constants.CategoryDormitoryPhotos,
		EntityId: res.Dormitory.Id,
	})
	if err != nil {
		s.logger.Warn("error getting dormitory photos", slog.String("error", err.Error()), slog.String("dormId", res.Dormitory.Id))
	}

	res.Dormitory.Photos = rmodel.ConvertFileInfos(photos)

	s.setDormitoryByIDToCache(ctx, res)

	return res, nil
}

func (s *CoreService) CreateDormitory(
	ctx context.Context,
	request *rmodel.CreateDormitoryRequest,
) (*rmodel.CreateDormitoryResponse, error) {
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
			DormitoryId:  dormitoryId,
			RoleRequired: true,
		},
	); err != nil {
		return nil, err
	}

	resp, err := s.repository.CreateDormitory(ctx, &dbtypes.CreateDormitoryRequest{
		DormitoryId:  request.DormitoryId,
		Name:         request.Name,
		Address:      request.Address,
		SupportEmail: request.SupportEmail,
		Description:  request.Description,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: error creating dormitory: %v", s.handleDBError(err), err)
	}

	res := new(rmodel.CreateDormitoryResponse).From(resp)

	return res, nil
}

func (s *CoreService) UpdateDormitory(
	ctx context.Context,
	request *rmodel.UpdateDormitoryRequest,
) (*rmodel.UpdateDormitoryResponse, error) {
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
			DormitoryId:  dormitoryId,
			RoleRequired: true,
		},
	); err != nil {
		return nil, err
	}

	resp, err := s.repository.UpdateDormitory(ctx, &dbtypes.UpdateDormitoryRequest{
		DormitoryId:  request.DormitoryId,
		Name:         request.Name,
		Address:      request.Address,
		SupportEmail: request.SupportEmail,
		Description:  request.Description,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: error updating dormitory: %v", s.handleDBError(err), err)
	}

	res := new(rmodel.UpdateDormitoryResponse).From(resp)

	s.invalidateDormitoryCache(ctx, request.DormitoryId)

	return res, nil
}

func (s *CoreService) DeleteDormitory(
	ctx context.Context,
	request *rmodel.DeleteDormitoryRequest,
) (*rmodel.DeleteDormitoryResponse, error) {
	return nil, ErrUnimplemented
}

func (s *CoreService) getDormitoriesFromCache(
	ctx context.Context,
) (*rmodel.GetDormitoriesResponse, error) {

	res, err := s.cacheClient.Get(
		ctx,
		constants.CacheDormitoriesKey,
		cache.CategoryDormitoryList,
	)
	if err == nil {
		s.logger.Debug("dormitories cache hit")

		var resp rmodel.GetDormitoriesResponse

		if err := json.Unmarshal([]byte(res), &resp); err != nil {
			s.logger.Warn("error unmarshalling cache response", slog.String("error", err.Error()))

			if err := s.cacheClient.Delete(ctx, constants.CacheDormitoriesKey, cache.CategoryDormitoryList); err != nil {
				s.logger.Warn("error invalidating dormitory list cache", slog.String("error", err.Error()))
			}

			return nil, fmt.Errorf("%w: error unmarshalling cache response: %v", ErrInternal, err)
		}

		return &resp, nil
	}

	s.logger.Error("error getting dormitories from cache", slog.String("error", err.Error()))

	return nil, fmt.Errorf("%w: error getting dormitories from cache: %v", ErrInternal, err)
}

func (s *CoreService) setDormitoriesToCache(
	ctx context.Context,
	resp *rmodel.GetDormitoriesResponse,
) error {
	respBytes, err := json.Marshal(resp)
	if err != nil {
		return fmt.Errorf("%w: error marshalling response: %v", ErrInternal, err)
	}

	s.logger.Debug("setting dormitories to cache", slog.String("response", string(respBytes)))

	if err := s.cacheClient.Set(
		ctx,
		constants.CacheDormitoriesKey,
		cache.CategoryDormitoryList,
		string(respBytes),
		constants.DefaultDormitoryListTTL,
	); err != nil {
		s.logger.Warn("error setting dormitories to cache",
			slog.String("response", string(respBytes)),
			slog.String("error", err.Error()))

		return fmt.Errorf("%w: error setting cache: %v", ErrInternal, err)
	}

	return nil
}

func (s *CoreService) getDormitoryByIdFromCache(
	ctx context.Context,
	dormitoryID string,
) (*rmodel.GetDormitoryByIdResponse, error) {

	res, err := s.cacheClient.Get(
		ctx,
		dormitoryID,
		cache.CategoryDormitory,
	)
	if err == nil {
		s.logger.Debug("dormitory cache hit")

		var resp rmodel.GetDormitoryByIdResponse

		if err := json.Unmarshal([]byte(res), &resp); err != nil {
			s.logger.Warn("error unmarshalling cache response", slog.String("error", err.Error()))

			if err := s.cacheClient.Delete(ctx, dormitoryID, cache.CategoryDormitory); err != nil {
				s.logger.Warn("error invalidating dormitory cache", slog.String("error", err.Error()))
			}

			return nil, fmt.Errorf("%w: error unmarshalling cache response: %v", ErrInternal, err)
		}

		return &resp, nil
	}

	s.logger.Error("error getting dormitory from cache", slog.String("error", err.Error()))

	return nil, fmt.Errorf("%w: error getting dormitory from cache: %v", ErrInternal, err)
}

func (s *CoreService) setDormitoryByIDToCache(
	ctx context.Context,
	resp *rmodel.GetDormitoryByIdResponse,
) error {
	respBytes, err := json.Marshal(resp)
	if err != nil {
		return fmt.Errorf("%w: error marshalling response: %v", ErrInternal, err)
	}

	s.logger.Debug("setting dormitory to cache", slog.String("response", string(respBytes)))

	if err := s.cacheClient.Set(
		ctx,
		resp.Dormitory.Id,
		cache.CategoryDormitory,
		string(respBytes),
		constants.DefaultDormitoryTTL,
	); err != nil {
		s.logger.Warn("error setting dormitory to cache",
			slog.String("response", string(respBytes)),
			slog.String("error", err.Error()))

		return fmt.Errorf("%w: error setting cache: %v", ErrInternal, err)
	}

	return nil
}

func (s *CoreService) invalidateDormitoryCache(
	ctx context.Context,
	dormitoryID string,
) error {
	s.logger.Debug("invalidating dormitory cache", slog.String("dormitoryId", dormitoryID))

	if err := s.cacheClient.Delete(
		ctx,
		dormitoryID,
		cache.CategoryDormitory); err != nil {
		s.logger.Warn("error invalidating dormitory cache", slog.String("error", err.Error()))
		return fmt.Errorf("%w: error invalidating dormitory cache: %v", ErrInternal, err)
	}

	return nil
}
