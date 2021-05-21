package main

import (
	"cloud.google.com/go/datastore"
	"cloud.google.com/go/pubsub"
	"context"
	"github.com/k-yomo/eitan/src/internal/pb/eitan"
	"github.com/k-yomo/eitan/src/notification_service/internal/email"
	"github.com/pkg/errors"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"google.golang.org/protobuf/proto"
)

type PubSubHandler struct {
	dsClient    *datastore.Client
	emailClient email.Client
}

func NewPubSubHandler(dsClient *datastore.Client, emailClient email.Client) *PubSubHandler {
	return &PubSubHandler{
		dsClient:    dsClient,
		emailClient: emailClient,
	}
}

func (p *PubSubHandler) HandleUserRegisteredEvent(ctx context.Context, m *pubsub.Message) error {
	userRegisteredEvent := eitan.UserRegisteredEvent{}
	if err := proto.Unmarshal(m.Data, &userRegisteredEvent); err != nil {
		return errors.Wrap(err, "proto.Unmarshal")
	}

	sgmail := mail.NewSingleEmail(
		&mail.Email{Name: userRegisteredEvent.DisplayName, Address: userRegisteredEvent.Email},
		"Welcome to Eitan!",
		&mail.Email{Name: userRegisteredEvent.DisplayName, Address: userRegisteredEvent.Email},
		// TODO: fix body
		`
Dear `+userRegisteredEvent.DisplayName+`,

Thank you for signing up for Eitan.

Best wishes,
Eitan Team
`,
		"",
	)

	if err := p.emailClient.Send(ctx, sgmail); err != nil {
		return err
	}

	return nil
}
