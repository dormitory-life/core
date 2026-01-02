package core

import (
	"errors"

	"github.com/dormitory-life/core/internal/database"
	dberrors "github.com/dormitory-life/core/internal/database/errors"
)

type CoreServiceConfig struct {
	Repository database.Repository
}
type CoreService struct {
	repository database.Repository
}

type CoreServiceClient interface {
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
