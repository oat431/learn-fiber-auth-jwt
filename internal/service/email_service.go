package service

import (
	"fmt"
	"net/smtp"
	"oat431/learn-fiber-auth-jwt/internal/config"
)

type SMTPService struct {
	config *config.Config
}

func NewSMTPService(cfg *config.Config) *SMTPService {
	return &SMTPService{config: cfg}
}

func (s *SMTPService) SendMail(to string) error {
	subject := "Subject: Go Email Integration \n"
	body := "This is a test email sent from a Go application using SMTP."
	message := []byte(subject + "\n" + body)

	auth := smtp.PlainAuth(
		"",
		s.config.SMTPUser,
		s.config.SMTPPassword,
		s.config.SMTPHost,
	)

	err := smtp.SendMail(
		fmt.Sprintf("%s:%d", s.config.SMTPHost, s.config.SMTPPort),
		auth,
		s.config.SMTPUser,
		[]string{to},
		message,
	)
	if err != nil {
		return err
	}

	return nil
}
