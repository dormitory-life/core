package auth

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	rmodel "github.com/dormitory-life/core/internal/server/request_models"
	pb "github.com/dormitory-life/core/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthClient struct {
	conn   *grpc.ClientConn
	client pb.AuthProtoServiceClient
	logger slog.Logger
}

type AuthClientConfig struct {
	GRPCAuthServerAddress string
	Timeout               time.Duration
	Logger                slog.Logger
}

func New(cfg AuthClientConfig) (*AuthClient, error) {
	conn, err := grpc.NewClient(
		cfg.GRPCAuthServerAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithTimeout(cfg.Timeout),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to auth service: %w", err)
	}

	return &AuthClient{
		conn:   conn,
		client: pb.NewAuthProtoServiceClient(conn),
		logger: cfg.Logger,
	}, nil
}

func (c *AuthClient) CheckAccess(
	ctx context.Context,
	req *rmodel.CheckAccessRequest,
) (*rmodel.CheckAccessResponse, error) {
	c.logger.Debug("checking access", slog.Any("request", req))

	callCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	resp, err := c.client.CheckAccess(callCtx, &pb.CheckAccessRequest{
		UserId:       req.UserId,
		DormitoryId:  req.DormitoryId,
		RoleRequired: req.RoleRequired,
	})

	if err != nil {
		switch err {
		case context.DeadlineExceeded:
			c.logger.Debug("checking access timeout error", slog.Any("request", req), slog.String("error", err.Error()))
			return &rmodel.CheckAccessResponse{
				Allowed: false,
				Reason:  "Auth service timeout",
				Error:   err,
			}, err
		default:
			c.logger.Debug("checking access error", slog.Any("request", req), slog.String("error", err.Error()))
			return &rmodel.CheckAccessResponse{
				Allowed: false,
				Reason:  "Internal auth service error",
				Error:   err,
			}, err
		}
	}

	c.logger.Debug("checking access success", slog.Any("request", req))

	return &rmodel.CheckAccessResponse{
		Allowed: resp.GetAllowed(),
		Reason:  resp.GetReason(),
		Error:   nil,
	}, nil
}
