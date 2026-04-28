package emailer

import (
	"context"
	"fmt"
	"html"
	"log/slog"
	"mime"
	"net/smtp"
	"strings"
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

	body := fmt.Sprintf(
		emailTemplate,
		html.EscapeString(req.UserEmail),
		html.EscapeString(req.Title),
		html.EscapeString(req.Description),
	)

	message := buildHTMLMessage(e.email, req.SupportEmail, supportHeader, body)

	if err := smtp.SendMail(
		fmt.Sprintf("%s:%d", e.host, e.port),
		auth,
		e.email,
		[]string{req.SupportEmail},
		[]byte(message),
	); err != nil {
		e.logger.Error("error sending mail", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%w: error sending email: %v", ErrInternal, err)
	}

	return &SendMessageResponse{
		UserEmail:    req.UserEmail,
		SupportEmail: req.SupportEmail,
	}, nil
}

func buildHTMLMessage(from, to, subject, body string) string {
	encodedSubject := mime.QEncoding.Encode("utf-8", subject)

	headers := []string{
		fmt.Sprintf("From: Dormitory Life Support <%s>", from),
		fmt.Sprintf("To: %s", to),
		fmt.Sprintf("Subject: %s", encodedSubject),
		"MIME-Version: 1.0",
		`Content-Type: text/html; charset="UTF-8"`,
		"Content-Transfer-Encoding: 8bit",
	}

	return strings.Join(headers, "\r\n") + "\r\n\r\n" + body
}
