// Package mailer — отправка писем. ConsoleSender используется когда SMTP не
// настроен (например в локальной разработке) — просто пишет письмо в лог,
// чтобы можно было скопировать ссылку подтверждения вручную.
package mailer

import (
	"fmt"
	"myproject/internal/logger"
	"net/smtp"
)

type SMTPSender struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

func (s *SMTPSender) Send(to, subject, body string) error {
	addr := fmt.Sprintf("%s:%s", s.Host, s.Port)
	auth := smtp.PlainAuth("", s.Username, s.Password, s.Host)
	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
		s.From, to, subject, body)
	return smtp.SendMail(addr, auth, s.From, []string{to}, []byte(msg))
}

type ConsoleSender struct{}

func (ConsoleSender) Send(to, subject, body string) error {
	logger.Log.Info().Str("to", to).Str("subject", subject).Msg("SMTP не настроен — письмо выведено в лог: " + body)
	return nil
}
