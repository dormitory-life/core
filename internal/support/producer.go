package support

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/dormitory-life/core/internal/broker"
)

type SupportProducer struct {
	broker broker.BrokerClient
	logger slog.Logger
}

type SupportProducerConfig struct {
	Broker broker.BrokerClient
	Logger slog.Logger
}

func NewSupportProducer(cfg SupportProducerConfig) SupportProducer {
	return SupportProducer{
		broker: cfg.Broker,
		logger: cfg.Logger,
	}
}

func (s *SupportProducer) PublishSupportMessage(
	ctx context.Context,
	msg *SupportMessage,
) error {
	jsonMessage, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("error marshalling support message: %w", err)
	}

	s.logger.Debug("publishing support message", slog.String("title", msg.Title), slog.String("desc", msg.Description))

	return s.broker.PublishSupportMessage(ctx, jsonMessage)
}
