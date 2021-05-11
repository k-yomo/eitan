package main

import (
	"cloud.google.com/go/datastore"
	"cloud.google.com/go/pubsub"
	"context"
	"github.com/golang/protobuf/proto"
	"github.com/k-yomo/eitan/src/internal/pb/eitan"
	"github.com/k-yomo/eitan/src/notification_service/internal/email"
	"github.com/k-yomo/eitan/src/pkg/logging"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"go.uber.org/zap"
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

func (n *NotificationHandler) HandleAccountRegistration(ctx context.Context, m *pubsub.Message) error {
	logger := logging.Logger(ctx)

	accountRegistrationEvent := eitan.AccountRegistrationEvent{}
	if err := proto.Unmarshal(m.Data, &accountRegistrationEvent); err != nil {
		logger.Error("unmarshal AccountRegistrationEvent failed", zap.Error(err))
		return err
	}

	sgmail := mail.NewSingleEmail(
		&mail.Email{Name: accountRegistrationEvent.DisplayName, Address: accountRegistrationEvent.Email},
		"Welcome to Eitan!",
		&mail.Email{Name: accountRegistrationEvent.DisplayName, Address: accountRegistrationEvent.Email},
		// TODO: fix body
		`
Dear ` + accountRegistrationEvent.DisplayName + `,

Thank you for signing up for Eitan.

Best wishes,
Eitan Team
`,
		"",
	)

	key := datastore.NameKey("AccountRegistrationEvent", accountRegistrationEvent.AccountId, nil)
	_, err := n.dsClient.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		if err := tx.Get(key, &eitan.AccountRegistrationEvent{}); err == nil {
			// entity exists means email already sent
			return nil
		} else {
			if err != datastore.ErrNoSuchEntity {
				return err
			}
		}

		if _, err := tx.Put(key, &accountRegistrationEvent); err != nil {
			return err
		}
		if err := n.emailClient.Send(ctx, sgmail); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		logger.Error("process AccountRegistrationEvent failed", zap.Error(err), zap.String("event", accountRegistrationEvent.String()))
		return err
	}

	return nil
}
