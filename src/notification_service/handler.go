package main

import (
	"cloud.google.com/go/datastore"
	"cloud.google.com/go/pubsub"
	"context"
	"github.com/golang/protobuf/proto"
	"github.com/k-yomo/eitan/src/internal/pb/eitan"
	"github.com/k-yomo/eitan/src/notification_service/internal/email"
	"github.com/pkg/errors"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type NotificationHandler struct {
	dsClient    *datastore.Client
	emailClient email.Client
}

func NewNotificationHandler(dsClient *datastore.Client, emailClient email.Client) *NotificationHandler {
	return &NotificationHandler{
		dsClient:    dsClient,
		emailClient: emailClient,
	}
}

func (n *NotificationHandler) HandleUserRegistration(ctx context.Context, m *pubsub.Message) error {
	userRegistrationEvent := eitan.UserRegistrationEvent{}
	if err := proto.Unmarshal(m.Data, &userRegistrationEvent); err != nil {
		return errors.Wrap(err, "proto.Unmarshal")
	}

	sgmail := mail.NewSingleEmail(
		&mail.Email{Name: userRegistrationEvent.DisplayName, Address: userRegistrationEvent.Email},
		"Welcome to Eitan!",
		&mail.Email{Name: userRegistrationEvent.DisplayName, Address: userRegistrationEvent.Email},
		// TODO: fix body
		`
Dear ` + userRegistrationEvent.DisplayName + `,

Thank you for signing up for Eitan.

Best wishes,
Eitan Team
`,
		"",
	)

	key := datastore.NameKey("UserRegistrationEvent", userRegistrationEvent.UserId, nil)
	_, err := n.dsClient.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		if err := tx.Get(key, &eitan.UserRegistrationEvent{}); err == nil {
			// entity exists means email already sent
			return nil
		} else {
			if err != datastore.ErrNoSuchEntity {
				return err
			}
		}

		if _, err := tx.Put(key, &userRegistrationEvent); err != nil {
			return err
		}
		if err := n.emailClient.Send(ctx, sgmail); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return errors.Wrap(err, "process UserRegistrationEvent")
	}

	return nil
}
