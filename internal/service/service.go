package core

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/dormitory-life/core/internal/auth"
	"github.com/dormitory-life/core/internal/database"
	dberrors "github.com/dormitory-life/core/internal/database/errors"
	rmodel "github.com/dormitory-life/core/internal/server/request_models"
)

type CoreServiceConfig struct {
	Repository database.Repository
	AuthClient *auth.AuthClient
	Logger     slog.Logger
}
type CoreService struct {
	repository database.Repository
	authClient *auth.AuthClient
	logger     slog.Logger
}

type CoreServiceClient interface {
	GetDormitories(ctx context.Context, request *rmodel.GetDormitoriesRequest) (*rmodel.GetDormitoriesResponse, error)
	GetDormitoryById(ctx context.Context, request *rmodel.GetDormitoryByIdRequest) (*rmodel.GetDormitoryByIdResponse, error)
	CreateDormitory(ctx context.Context, request *rmodel.CreateDormitoryRequest) (*rmodel.CreateDormitoryResponse, error)
	UpdateDormitory(ctx context.Context, request *rmodel.UpdateDormitoryRequest) (*rmodel.UpdateDormitoryResponse, error)
	DeleteDormitory(ctx context.Context, request *rmodel.DeleteDormitoryRequest) (*rmodel.DeleteDormitoryResponse, error)

	GetDormitoriesAvgGrades(ctx context.Context, request *rmodel.GetDormitoriesAvgGradesRequest) (*rmodel.GetDormitoriesAvgGradesResponse, error)
	GetDormitoryAvgGrades(ctx context.Context, request *rmodel.GetDormitoryAvgGradesRequest) (*rmodel.GetDormitoryAvgGradesResponse, error)
	CreateDormitoryGrade(ctx context.Context, request *rmodel.CreateDormitoryGradeRequest) (*rmodel.CreateDormitoryGradeResponse, error)
}

func New(cfg CoreServiceConfig) CoreServiceClient {
	return &CoreService{
		repository: cfg.Repository,
		authClient: cfg.AuthClient,
		logger:     cfg.Logger,
	}
}

func (s *CoreService) checkAccess(
	ctx context.Context,
	req *rmodel.CheckAccessRequest,
) error {
	if req == nil {
		return fmt.Errorf("%w: request is nil", ErrBadRequest)
	}

	resp, err := s.authClient.CheckAccess(ctx, req)
	if err != nil {
		return fmt.Errorf(
			"%w: error checking access for user %s dormitory %s: %v",
			ErrInternal,
			req.UserId,
			req.DormitoryId,
			err,
		)
	}

	if resp == nil {
		return fmt.Errorf("%w: empty response from auth service", ErrInternal)
	}

	if !resp.Allowed {
		return fmt.Errorf(
			"%w: access denied for user %s dormitory %s: %s",
			ErrForbidden,
			req.UserId,
			req.DormitoryId,
			resp.Reason,
		)
	}

	return nil
}

func (s *CoreService) extractIdsFromRequestContext(ctx context.Context) (string, string, error) {
	userId := ctx.Value("userId")
	dormitoryId := ctx.Value("dormitoryId")
	if fmt.Sprintf("%s", userId) == "" {
		return "", "", fmt.Errorf("%w: empty userId in context", ErrBadRequest)
	}

	if fmt.Sprintf("%s", dormitoryId) == "" {
		return "", "", fmt.Errorf("%w: empty userId in context", ErrBadRequest)
	}

	userIdStr := fmt.Sprintf("%s", userId)
	dormitoryIdStr := fmt.Sprintf("%s", dormitoryId)

	s.logger.Debug("extracted from context", slog.String("userId", userIdStr), slog.String("dormitoryId", dormitoryIdStr))
	return userIdStr, dormitoryIdStr, nil
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
