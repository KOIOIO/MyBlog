// Package email 提供邮件发送能力（构造注入，不依赖 global）。
package email

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"

	"server/config"

	"github.com/jordan-wright/email"
)

// Sender 封装邮件发送所需的配置。
type Sender struct {
	cfg config.Email
}

// New 根据邮件配置构造发送器。
func New(cfg config.Email) *Sender {
	return &Sender{cfg: cfg}
}

// Send 向逗号分隔的收件人列表发送 HTML 邮件。
func (s *Sender) Send(to, subject, body string) error {
	recipients := strings.Split(to, ",")
	return s.send(recipients, subject, body)
}

// send 执行邮件发送操作。
func (s *Sender) send(to []string, subject string, body string) error {
	cfg := s.cfg
	auth := smtp.PlainAuth("", cfg.From, cfg.Secret, cfg.Host)

	e := email.NewEmail()
	if cfg.Nickname != "" {
		e.From = fmt.Sprintf("%s <%s>", cfg.Nickname, cfg.From)
	} else {
		e.From = cfg.From
	}
	e.To = to
	e.Subject = subject
	e.HTML = []byte(body)

	hostAddr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	if cfg.IsSSL {
		return e.SendWithTLS(hostAddr, auth, &tls.Config{ServerName: cfg.Host})
	}
	return e.Send(hostAddr, auth)
}
