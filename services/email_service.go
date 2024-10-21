package services

import (
	"net/smtp"
	"os"
)

var smtpSendMail = smtp.SendMail

type IEmailService interface {
	SendMail(receiver_address, receiver_name, email_content string) error
}

type EmailService struct {
}

func NewEmailService() *EmailService {
	return &EmailService{}
}

func (s *EmailService) SendMail(receiver_address, receiver_name, email_content string) error {
	from := os.Getenv("SMTP_USER")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	auth := smtp.PlainAuth("", from, os.Getenv("SMTP_PASSWORD"), smtpHost)

	return smtpSendMail(smtpHost+":"+smtpPort, auth, from, []string{receiver_address}, []byte(email_content))
}
