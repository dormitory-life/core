package core

import (
	"context"
	"fmt"

	dbtypes "github.com/dormitory-life/core/internal/database/types"
	rmodel "github.com/dormitory-life/core/internal/server/request_models"
	"github.com/dormitory-life/core/internal/support"
)

func (s *CoreService) CreateSupportRequest(
	ctx context.Context,
	request *rmodel.CreateSupportRequest,
) (*rmodel.CreateSupportResponse, error) {
	if request == nil {
		return nil, fmt.Errorf("%w: request is nil", ErrBadRequest)
	}

	userId, dormitoryId, err := s.extractIdsFromRequestContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: error getting ids from context: %v", ErrInternal, err)
	}

	resp, err := s.repository.GetEmailsForSupport(ctx, &dbtypes.GetEmailsForSupportRequest{
		UserId:      userId,
		DormitoryId: dormitoryId,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: error getting emails for support request: %v", s.handleDBError(err), err)
	}

	if err := s.supportClient.PublishSupportMessage(ctx, &support.SupportMessage{
		UserEmail:    resp.UserEmail,
		SupportEmail: resp.SupportEmail,
		Title:        request.Title,
		Description:  request.Description,
	}); err != nil {
		return nil, fmt.Errorf("%w: error publish message: %v", ErrInternal, err)
	}

	return &rmodel.CreateSupportResponse{
		Message: fmt.Sprintf("Sent message\nfrom user: %s\nto support: %s", resp.UserEmail, resp.SupportEmail),
	}, nil
}
