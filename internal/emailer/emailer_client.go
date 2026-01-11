package emailer

import (
	"context"
	"fmt"
	"log/slog"
	"net/smtp"
)

type EmailerConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Email    string
	Logger   slog.Logger
}

type Emailer struct {
	host     string
	port     int
	user     string
	password string
	email    string
	logger   slog.Logger
}

func New(cfg *EmailerConfig) *Emailer {
	return &Emailer{
		host:     cfg.Host,
		port:     cfg.Port,
		user:     cfg.User,
		password: cfg.Password,
		email:    cfg.Email,
		logger:   cfg.Logger,
	}
}

func (e *Emailer) SendMessage(
	ctx context.Context,
	req *SendMessageRequest,
) (*SendMessageResponse, error) {
	slog.Info("sending mail", slog.Any("mail", req))
	if req == nil {
		return nil, fmt.Errorf("%w: request is nil", ErrBadRequest)
	}

	auth := smtp.PlainAuth("", e.user, e.password, e.host)

	mail := fmt.Sprintf(email_template, req.UserEmail, req.Title, req.Description)

	if err := smtp.SendMail(
		fmt.Sprintf("%s:%d", e.host, e.port),
		auth,
		e.email,
		[]string{req.SupportEmail},
		[]byte(mail),
	); err != nil {
		e.logger.Error("error sending mail", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%w: error sending email: %v", ErrInternal, err)
	}

	return &SendMessageResponse{
		UserEmail:    req.UserEmail,
		SupportEmail: req.SupportEmail,
	}, nil
}
