package core

import (
	"context"
	"fmt"
	"log/slog"

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

	return res, nil
}

func (s *CoreService) GetDormitoryById(
	ctx context.Context,
	request *rmodel.GetDormitoryByIdRequest,
) (*rmodel.GetDormitoryByIdResponse, error) {
	if request == nil {
		return nil, fmt.Errorf("%w: request is nil", ErrBadRequest)
	}

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

	return res, nil
}

func (s *CoreService) DeleteDormitory(
	ctx context.Context,
	request *rmodel.DeleteDormitoryRequest,
) (*rmodel.DeleteDormitoryResponse, error) {
	return nil, ErrUnimplemented
}
