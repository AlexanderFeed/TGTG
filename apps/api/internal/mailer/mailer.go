package mailer

import (
	"context"
	"crypto/tls"
	"fmt"
	"mime"
	"net"
	"net/mail"
	"net/smtp"
	"strings"
	"time"
)

type Sender interface {
	SendVerificationCode(ctx context.Context, recipient, code, purpose string) error
}

type DevelopmentSender struct{}

func (DevelopmentSender) SendVerificationCode(context.Context, string, string, string) error {
	return nil
}

type SMTPSender struct {
	host     string
	port     string
	username string
	password string
	from     mail.Address
	timeout  time.Duration
}

func NewSMTP(host, port, username, password, from string) (*SMTPSender, error) {
	parsedFrom, err := mail.ParseAddress(from)
	if err != nil {
		return nil, fmt.Errorf("parse SMTP_FROM: %w", err)
	}
	return &SMTPSender{
		host: host, port: port, username: username, password: password,
		from: *parsedFrom, timeout: 12 * time.Second,
	}, nil
}

func (s *SMTPSender) SendVerificationCode(ctx context.Context, recipient, code, purpose string) error {
	to, err := mail.ParseAddress(recipient)
	if err != nil {
		return fmt.Errorf("parse recipient: %w", err)
	}
	address := net.JoinHostPort(s.host, s.port)
	conn, err := (&net.Dialer{Timeout: s.timeout}).DialContext(ctx, "tcp", address)
	if err != nil {
		return fmt.Errorf("connect smtp: %w", err)
	}
	defer conn.Close()
	_ = conn.SetDeadline(time.Now().Add(s.timeout))

	client, err := smtp.NewClient(conn, s.host)
	if err != nil {
		return fmt.Errorf("open smtp client: %w", err)
	}
	defer client.Close()

	if ok, _ := client.Extension("STARTTLS"); !ok {
		return fmt.Errorf("smtp server does not advertise STARTTLS")
	}
	if err := client.StartTLS(&tls.Config{ServerName: s.host, MinVersion: tls.VersionTLS12}); err != nil {
		return fmt.Errorf("start smtp tls: %w", err)
	}
	if s.username != "" {
		if ok, _ := client.Extension("AUTH"); !ok {
			return fmt.Errorf("smtp server does not advertise AUTH")
		}
		if err := client.Auth(smtp.PlainAuth("", s.username, s.password, s.host)); err != nil {
			return fmt.Errorf("authenticate smtp: %w", err)
		}
	}
	if err := client.Mail(s.from.Address); err != nil {
		return fmt.Errorf("set smtp sender: %w", err)
	}
	if err := client.Rcpt(to.Address); err != nil {
		return fmt.Errorf("set smtp recipient: %w", err)
	}
	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("open smtp data: %w", err)
	}

	action := "входа"
	if purpose == "register" {
		action = "регистрации"
	}
	subject := mime.QEncoding.Encode("UTF-8", "Код подтверждения — ЕщёЕсть")
	body := fmt.Sprintf(
		"Ваш код для %s в ЕщёЕсть: %s\r\n\r\nКод действует ограниченное время. Если вы не запрашивали его, просто проигнорируйте письмо.\r\n",
		action, code,
	)
	message := strings.Join([]string{
		"From: " + s.from.String(),
		"To: " + to.String(),
		"Subject: " + subject,
		"Date: " + time.Now().Format(time.RFC1123Z),
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=UTF-8",
		"Content-Transfer-Encoding: 8bit",
		"",
		body,
	}, "\r\n")
	if _, err := writer.Write([]byte(message)); err != nil {
		writer.Close()
		return fmt.Errorf("write smtp data: %w", err)
	}
	if err := writer.Close(); err != nil {
		return fmt.Errorf("close smtp data: %w", err)
	}
	if err := client.Quit(); err != nil {
		return fmt.Errorf("quit smtp: %w", err)
	}
	return nil
}
