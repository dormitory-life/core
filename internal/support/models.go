package support

import amqp "github.com/rabbitmq/amqp091-go"

type SupportMessage struct {
	UserEmail    string `json:"user_email"`
	SupportEmail string `json:"support_email"`
	Title        string `json:"title"`
	Description  string `json:"description"`
}

type SupportJob struct {
	msg SupportMessage
	job amqp.Delivery
}
