package main

import (
	"cloud.google.com/go/pubsub"
	"context"
	"database/sql"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/jmoiron/sqlx"
	"github.com/k-yomo/eitan/src/account_service/infra"
	"github.com/k-yomo/eitan/src/account_service/internal/sessionmanager"
	"github.com/k-yomo/eitan/src/internal/pb/eitan"
	"github.com/k-yomo/eitan/src/pkg/clock"
	"github.com/k-yomo/eitan/src/pkg/event"
	"github.com/k-yomo/eitan/src/pkg/logging"
	"github.com/k-yomo/eitan/src/pkg/tx"
	"github.com/k-yomo/eitan/src/pkg/uuid"
	"github.com/markbates/goth/gothic"
	"go.uber.org/zap"
	"net/http"
	"net/url"
)

type AuthHandler struct {
	sessionManager sessionmanager.SessionManager
	db             *sqlx.DB
	txManager      tx.Manager
	pubsubClient   *pubsub.Client
	webAppURL      string
}

func NewAuthHandler(sessionManager sessionmanager.SessionManager, db *sqlx.DB, pubsubClient *pubsub.Client, webAppURL string) *AuthHandler {
	return &AuthHandler{sessionManager: sessionManager, db: db, pubsubClient: pubsubClient, webAppURL: webAppURL}
}

func (a *AuthHandler) HandleOAuth(w http.ResponseWriter, r *http.Request) {
	gothic.BeginAuthHandler(w, r)
}

func (a *AuthHandler) HandleOAuthCallback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		handleServerError(ctx, err, w)
		return
	}
	// Logout from gothic since we manage session on our own
	if err := gothic.Logout(w, r); err != nil {
		handleServerError(ctx, err, w)
		return
	}

	account, err := infra.AccountByEmail(ctx, a.db, user.Email)
	if err != nil && err != sql.ErrNoRows {
		handleServerError(ctx, err, w)
		return
	}

	// Create account if not exist
	if err == sql.ErrNoRows {
		now := clock.Now()
		var screenImgURL sql.NullString
		if avatarURL, err := url.Parse(user.AvatarURL); err == nil {
			screenImgURL = sql.NullString{String: avatarURL.String(), Valid: true}
		}

		account = &infra.Account{
			ID:           uuid.Generate(),
			Provider:     user.Provider,
			Email:        user.Email,
			DisplayName:  user.Name,
			ScreenImgURL: screenImgURL,
			CreatedAt:    now,
			UpdatedAt:    now,
		}

		if err := a.createAccount(ctx, account); err != nil {
			handleServerError(ctx, err, w)
		}
	}

	if err := a.sessionManager.Login(w, r, account.ID); err != nil {
		handleServerError(ctx, err, w)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("%s/account_settings", a.webAppURL), http.StatusFound)
}

func (a *AuthHandler) createAccount(ctx context.Context, account *infra.Account) error {
	logger := logging.Logger(ctx)

	if err := account.Insert(ctx, a.db); err != nil {
		return err
	}

	m := eitan.AccountRegistrationEvent{
		AccountId:   account.ID,
		Provider:    account.Provider,
		Email:       account.Email,
		DisplayName: account.DisplayName,
	}
	mBytes, err := proto.Marshal(&m)
	if err != nil {
		logger.Error("marshal AccountRegistrationEvent failed", zap.Error(err))
	} else {
		// TODO: retry publishing
		t := a.pubsubClient.Topic(event.AccountRegistrationTopicName)
		if _, err := t.Publish(ctx, &pubsub.Message{Data: mBytes}).Get(ctx); err != nil {
			logger.Error("publish AccountRegistrationEvent failed", zap.Error(err), zap.String("AccountRegistrationEvent", m.String()))
		}
		logger.Debug("published AccountRegistrationEvent", zap.Error(err))
	}

	return nil
}

func (a *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if err := a.sessionManager.Logout(w, r); err != nil {
		handleServerError(r.Context(), err, w)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("%s/login", a.webAppURL), http.StatusFound)
}

func handleServerError(ctx context.Context, err error, w http.ResponseWriter) {
	logging.Logger(ctx).Error(err.Error(), zap.Error(err))
	w.WriteHeader(500)
	w.Write([]byte(`{"status":"500"}`))
}
