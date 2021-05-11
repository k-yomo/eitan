package email

import (
	"context"
	"fmt"
	"github.com/k-yomo/eitan/src/pkg/logging"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"go.uber.org/zap"
)

type Client interface {
	Send(ctx context.Context, email *mail.SGMailV3) error
}

type sendgridEmailClient struct {
	sendgridClient *sendgrid.Client
}

func NewSendgridEmailClient(sendgridClient *sendgrid.Client) *sendgridEmailClient {
	return &sendgridEmailClient{sendgridClient: sendgridClient}
}

func (s *sendgridEmailClient) Send(ctx context.Context, email *mail.SGMailV3) error {
	resp, err := s.sendgridClient.SendWithContext(ctx, email)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("send message failed with status: %v, body: %v", resp.StatusCode, resp.Body)
	}
	return err
}

type noopEmailClient struct {
}

func NewNoopEmailClient() *noopEmailClient {
	return &noopEmailClient{}
}

func (n *noopEmailClient) Send(ctx context.Context, email *mail.SGMailV3) error {
	logging.Logger(ctx).Info("email sent", zap.Any("email", email))
	return nil
}
