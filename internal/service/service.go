package core

import (
	"context"
	"errors"

	"github.com/dormitory-life/core/internal/database"
	dberrors "github.com/dormitory-life/core/internal/database/errors"
	rmodel "github.com/dormitory-life/core/internal/server/request_models"
)

type CoreServiceConfig struct {
	Repository database.Repository
}
type CoreService struct {
	repository database.Repository
}

type CoreServiceClient interface {
	GetDormitories(ctx context.Context, request *rmodel.GetDormitoriesRequest) (*rmodel.GetDormitoriesResponse, error)
	GetDormitoryById(ctx context.Context, request *rmodel.GetDormitoryByIdRequest) (*rmodel.GetDormitoryByIdResponse, error)
	CreateDormitory(ctx context.Context, request *rmodel.CreateDormitoryRequest) (*rmodel.CreateDormitoryResponse, error)
	UpdateDormitory(ctx context.Context, request *rmodel.UpdateDormitoryRequest) (*rmodel.UpdateDormitoryResponse, error)
	DeleteDormitory(ctx context.Context, request *rmodel.DeleteDormitoryRequest) (*rmodel.DeleteDormitoryResponse, error)
}

func New(cfg CoreServiceConfig) CoreServiceClient {
	return &CoreService{
		repository: cfg.Repository,
	}
}

func (s *CoreService) handleDBError(err error) error {
	switch {
	case errors.Is(err, dberrors.ErrBadRequest):
		return ErrBadRequest
	case errors.Is(err, dberrors.ErrNotFound):
		return ErrNotFound
	case errors.Is(err, dberrors.ErrInternal):
		return ErrInternal
	case errors.Is(err, dberrors.ErrConflict):
		return ErrConflict
	default:
		return ErrInternal
	}
}
