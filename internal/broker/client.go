package broker

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

type BrokerClient interface {
	Connect() error
	Disconnect() error

	SupportClient
}

type SupportClient interface {
	PublishSupportMessage(ctx context.Context, message []byte) error
	ConsumeSupportMessages(ctx context.Context) (<-chan amqp.Delivery, error)
}

func New(cfg RabbitMQBrokerConfig) BrokerClient {
	return newRabbitMQBroker(cfg)
}

type QueueConfig struct {
	SupportQueueName string
}

func ConfigureQueues(cfg QueueConfig) {
	supportQueueName = cfg.SupportQueueName
}
