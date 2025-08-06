package notifier

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/mailgun/mailgun-go/v4"
	gomail "gopkg.in/mail.v2"
	"social-app/internal/config"
)

type Mailerx interface {
	SendEmailVerification(to string, code string) error
}

type MailTrap struct {
	Host     string
	AuthUser string
	AuthPass string
	From     string
	Port     int
}

func NewMailTrap(cfg config.Mailtrap) *MailTrap {
	return &MailTrap{
		Host:     cfg.Host,
		Port:     cfg.Port,
		AuthUser: cfg.AuthUser,
		AuthPass: cfg.AuthPass,
		From:     cfg.From,
	}
}

func (m *MailTrap) SendEmailVerification(to, code string) error {
	message := gomail.NewMessage()

	// Set email headers
	message.SetHeader("From", "youremail@email.com")
	message.SetHeader("To", to)
	message.SetHeader("Subject", "Hello from the Mailtrap team")

	message.SetBody("text/plain", "Verification code: "+code)

	dialer := gomail.NewDialer(m.Host, m.Port, m.AuthUser, m.AuthPass)

	if err := dialer.DialAndSend(message); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

type Mailgun struct {
	Domain string
	APIKey string
	From   string
}

func NewMailgun(cfg config.Mailgun) *Mailgun {
	return &Mailgun{Domain: cfg.Host, APIKey: cfg.ApiKey, From: cfg.From}
}

func (m *Mailgun) SendEmailVerification(to, code string) error {
	mg := mailgun.NewMailgun(m.Domain, m.APIKey)

	subject := "Vérification de votre email"
	body := fmt.Sprintf("Voici votre code de vérification : %s", code)

	message := mailgun.NewMessage(m.From, subject, body, to)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, _, err := mg.Send(ctx, message)
	if err != nil {
		log.Printf("[Mailgun] Erreur d'envoi: %v", err)
	}

	log.Printf("[Mailgun] Email sent successfully to %s, with code: %v", to, code)

	return nil
}
