package core

import (
	"context"
	"fmt"

	dbtypes "github.com/dormitory-life/core/internal/database/types"
	rmodel "github.com/dormitory-life/core/internal/server/request_models"
)

func (s *CoreService) GetChat(
	ctx context.Context,
	request *rmodel.GetChatMessagesRequest,
) (*rmodel.GetChatMessagesResponse, error) {
	if request == nil {
		return nil, fmt.Errorf("%w: request is nil", ErrBadRequest)
	}

	resp, err := s.repository.GetChatMessages(ctx, &dbtypes.GetChatMessagesRequest{
		DormitoryID: request.DormitoryID,
		Page:        request.Page,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: error getting chat: %v", s.handleDBError(err), err)
	}

	res := new(rmodel.GetChatMessagesResponse).From(resp)

	return res, nil
}

func (s *CoreService) CreateChatMessage(
	ctx context.Context,
	request *rmodel.CreateChatMessageRequest,
) (*rmodel.CreateChatMessageResponse, error) {
	if request == nil {
		return nil, fmt.Errorf("%w: request is nil", ErrBadRequest)
	}

	userId, _, err := s.extractIdsFromRequestContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: error getting ids from context: %v", ErrInternal, err)
	}

	if err := s.checkAccess(
		ctx,
		&rmodel.CheckAccessRequest{
			UserId:       userId,
			DormitoryId:  request.DormitoryID,
			RoleRequired: false,
		},
	); err != nil {
		return nil, err
	}

	resp, err := s.repository.CreateChatMessage(ctx, &dbtypes.CreateChatMessageRequest{
		DormitoryID: request.DormitoryID,
		UserID:      userId,
		Text:        request.Text,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: error creating message: %v", s.handleDBError(err), err)
	}

	res := new(rmodel.CreateChatMessageResponse).From(resp)

	return res, nil
}
