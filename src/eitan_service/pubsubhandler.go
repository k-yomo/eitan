package main

import (
	"cloud.google.com/go/pubsub"
	"context"
	"database/sql"
	"github.com/golang/protobuf/proto"
	"github.com/jmoiron/sqlx"
	"github.com/k-yomo/eitan/src/eitan_service/infra"
	"github.com/k-yomo/eitan/src/internal/pb/eitan"
	"github.com/k-yomo/eitan/src/pkg/clock"
	"github.com/k-yomo/eitan/src/pkg/tx"
	"github.com/k-yomo/eitan/src/pkg/uuid"
	"github.com/pkg/errors"
)

type PubSubHandler struct {
	db        *sqlx.DB
	txManager tx.Manager
}

func NewPubSubHandler(db *sqlx.DB) *PubSubHandler {
	return &PubSubHandler{
		db:        db,
		txManager: tx.NewManager(db),
	}
}

func (p *PubSubHandler) HandleUserRegistration(ctx context.Context, m *pubsub.Message) error {
	userRegistrationEvent := eitan.UserRegistrationEvent{}
	if err := proto.Unmarshal(m.Data, &userRegistrationEvent); err != nil {
		return errors.Wrap(err, "proto.Unmarshal")
	}

	err := p.txManager.RunInTx(ctx, func(ctx context.Context) error {
		_, err := infra.PlayerByUserID(ctx, p.db, userRegistrationEvent.UserId)
		if err == sql.ErrNoRows {
			now := clock.Now()
			player := &infra.Player{
				ID:     uuid.Generate(),
				UserID: userRegistrationEvent.UserId,
				CreatedAt: now,
				UpdatedAt: now,
			}
			return player.Insert(ctx, p.db)
		} else if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
