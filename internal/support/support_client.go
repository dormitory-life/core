package support

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/dormitory-life/core/internal/broker"
	"github.com/dormitory-life/core/internal/emailer"
)

type SupportClient interface {
	PublishSupportMessage(ctx context.Context, msg *SupportMessage) error
	ProcessSupportMessage(ctx context.Context, supportJob *SupportJob) error
}

type SupportClientConfig struct {
	Broker  broker.BrokerClient
	Emailer *emailer.Emailer
	Logger  slog.Logger
}

type SupportSvc struct {
	broker   broker.BrokerClient
	producer SupportProducer
	emailer  *emailer.Emailer
	logger   slog.Logger
}

func New(cfg *SupportClientConfig) SupportClient {
	producer := NewSupportProducer(SupportProducerConfig{
		Broker: cfg.Broker,
		Logger: cfg.Logger,
	})

	return &SupportSvc{
		broker:   cfg.Broker,
		producer: producer,
		emailer:  cfg.Emailer,
		logger:   cfg.Logger,
	}
}

func (s *SupportSvc) PublishSupportMessage(
	ctx context.Context,
	msg *SupportMessage,
) error {
	return s.producer.PublishSupportMessage(ctx, msg)
}

func (s *SupportSvc) ProcessSupportMessage(
	ctx context.Context,
	supportJob *SupportJob,
) error {
	if supportJob == nil {
		return fmt.Errorf("%w: message is nil", ErrBadRequest)
	}

	s.logger.Debug("processing support message", slog.Any("msg", supportJob.msg))

	resp, err := s.emailer.SendMessage(ctx, &emailer.SendMessageRequest{
		UserEmail:    supportJob.msg.UserEmail,
		SupportEmail: supportJob.msg.SupportEmail,
		Title:        supportJob.msg.Title,
		Description:  supportJob.msg.Description,
	})
	if err != nil {
		s.logger.Error("error sending support message", slog.Any("msg", supportJob.msg))

		return fmt.Errorf("%w: error sending support message: %v", ErrInternal, err)
	}

	s.logger.Debug("sent support mail", slog.String("from student", resp.UserEmail), slog.String("to support", resp.SupportEmail))

	return nil
}
