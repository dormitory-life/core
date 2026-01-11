package support

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/dormitory-life/core/internal/broker"
	amqp "github.com/rabbitmq/amqp091-go"
)

type SupportConsumerConfig struct {
	Broker        broker.BrokerClient
	SupportClient SupportClient
	Logger        slog.Logger
}

type SupportConsumer struct {
	broker        broker.BrokerClient
	supportClient SupportClient
	logger        slog.Logger
}

func NewSupportConsumer(cfg *SupportConsumerConfig) SupportConsumer {
	return SupportConsumer{
		broker:        cfg.Broker,
		supportClient: cfg.SupportClient,
		logger:        cfg.Logger,
	}
}

func (c *SupportConsumer) Start(ctx context.Context) error {
	msgStream, err := c.broker.ConsumeSupportMessages(ctx)
	if err != nil {
		return fmt.Errorf("error starting support consumer: %w", err)
	}

	c.logger.Debug("support consumer started", slog.Any("msgStream", msgStream))

	go c.ProcessSupportMessages(ctx, msgStream)

	return nil
}

func (c *SupportConsumer) ProcessSupportMessages(ctx context.Context, msgStream <-chan amqp.Delivery) {
	c.logger.Debug("processing support messages")

	for {
		select {
		case <-ctx.Done():
			c.logger.Info("Context cancelled - stopping support message processing", slog.String("error", ctx.Err().Error()))
			return
		case delivery, ok := <-msgStream:
			if !ok {
				c.logger.Error("Message stream channel closed - broker connection lost")
				return
			}

			c.handleDelivery(ctx, delivery)
		}
	}
}

func (c *SupportConsumer) handleDelivery(ctx context.Context, delivery amqp.Delivery) {
	c.logger.Info("handling delivery", slog.Any("delivery", delivery))

	var supportMessage SupportMessage
	if err := json.Unmarshal(delivery.Body, &supportMessage); err != nil {
		c.logger.Error("failed to unmarshal support message - sending NACK",
			slog.Any("message", delivery.Body),
			"error", err)

		// parsing error - set nack to broker message
		if nackErr := delivery.Nack(false, true); nackErr != nil {
			c.logger.Error("failed to send NACK", slog.String("error", nackErr.Error()))
		}

		return
	}

	supportJob := &SupportJob{
		msg: supportMessage,
		job: delivery,
	}

	c.supportClient.ProcessSupportMessage(ctx, supportJob)
}
