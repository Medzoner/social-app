package notifier

import (
	"fmt"
	"log"

	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
	"social-app/internal/config"
)

type SMSNotifier interface {
	SendPhoneVerification(to string, code string) error
}

type SMS struct {
	Client     *twilio.RestClient
	FromNumber string
}

func NewSMS(cfg config.SMS) *SMS {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: cfg.AccountID,
		Password: cfg.AuthToken,
	})
	return &SMS{
		Client:     client,
		FromNumber: cfg.From,
	}
}

func (s *SMS) SendPhoneVerification(to, code string) error {
	msg := fmt.Sprintf("Votre code de vérification est : %s", code)
	params := &openapi.CreateMessageParams{}
	params.SetTo(to)
	params.SetFrom(s.FromNumber)
	params.SetBody(msg)

	_, err := s.Client.Api.CreateMessage(params)
	if err != nil {
		log.Printf("[Twilio] Erreur d'envoi SMS: %v", err)
	}

	log.Printf("[Twilio] SMS envoyé à %s: %s", to, msg)

	return nil
}
