package core

import (
	"context"
	"fmt"

	dbtypes "github.com/dormitory-life/core/internal/database/types"
	rmodel "github.com/dormitory-life/core/internal/server/request_models"
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

	return res, nil
}
