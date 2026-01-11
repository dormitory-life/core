package broker

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQBroker struct {
	conn *amqp.Connection
	cfg  RabbitMQBrokerConfig
}

type RabbitMQBrokerConfig struct {
	Host     string
	Port     int
	User     string
	Password string
}

func newRabbitMQBroker(cfg RabbitMQBrokerConfig) *RabbitMQBroker {
	return &RabbitMQBroker{
		cfg: cfg,
	}
}

func (r *RabbitMQBroker) Connect() error {
	conn, err := amqp.Dial(
		fmt.Sprintf(
			"amqp://%s:%s@%s:%d",
			r.cfg.User,
			r.cfg.Password,
			r.cfg.Host,
			r.cfg.Port,
		))

	if err != nil {
		return err
	}

	r.conn = conn

	return nil
}

func (r *RabbitMQBroker) Disconnect() error {
	return r.conn.Close()
}

func (r *RabbitMQBroker) PublishSupportMessage(
	ctx context.Context,
	message []byte,
) error {
	ch, err := r.conn.Channel()
	if err != nil {
		return fmt.Errorf("error open rmq.Connection.Channel: %w", err)
	}

	queue, err := ch.QueueDeclare(
		supportQueueName,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return fmt.Errorf("error declaring support queue: %w", err)
	}

	if err := ch.PublishWithContext(ctx,
		"",
		queue.Name,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			Body:         message,
		},
	); err != nil {
		return fmt.Errorf("error publishing support message: %w", err)
	}

	return nil
}

func (r *RabbitMQBroker) ConsumeSupportMessages(
	ctx context.Context,
) (<-chan amqp.Delivery, error) {
	ch, err := r.conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("error open rmq.Connection.Channel: %w", err)
	}

	queue, err := ch.QueueDeclare(
		supportQueueName,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("error declaring support queue: %w", err)
	}

	messages, err := ch.Consume(
		queue.Name, // queue name
		"",         // consumer
		false,      // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		return nil, fmt.Errorf("error consuming support messages: %w", err)
	}

	return messages, nil
}
