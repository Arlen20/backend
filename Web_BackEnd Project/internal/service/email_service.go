package service

import (
	"fmt"
	"io"

	"gopkg.in/gomail.v2"
)

// EmailServiceInterface определяет интерфейс для Email сервиса
type EmailServiceInterface interface {
	SendEmail(to, subject, body string, attachment io.Reader, filename string) error
}

type emailService struct {
	smtpHost     string
	smtpPort     int
	smtpUsername string
	smtpPassword string
}

// NewEmailService создает новый экземпляр email сервиса
func NewEmailService() EmailServiceInterface {
	return &emailService{
		smtpHost:     "smtp.gmail.com",
		smtpPort:     587,
		smtpUsername: "nurlybaynurbol@gmail.com", // Лучше получать из конфигурации
		smtpPassword: "rdhk amua afhc mivw",      // Лучше получать из конфигурации
	}
}

// SendEmail отправляет email с опциональным вложением
func (s *emailService) SendEmail(to, subject, body string, attachment io.Reader, filename string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.smtpUsername)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	if attachment != nil && filename != "" {
		m.Attach(filename, gomail.SetCopyFunc(func(w io.Writer) error {
			_, err := io.Copy(w, attachment)
			return err
		}))
	}

	d := gomail.NewDialer(s.smtpHost, s.smtpPort, s.smtpUsername, s.smtpPassword)

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
