package main

import (
	"cloud.google.com/go/pubsub"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/k-yomo/eitan/src/account_service/infra"
	"github.com/k-yomo/eitan/src/account_service/internal/sessionmanager"
	"github.com/k-yomo/eitan/src/pkg/clock"
	"github.com/k-yomo/eitan/src/pkg/tx"
	"github.com/k-yomo/eitan/src/pkg/uuid"
	"github.com/markbates/goth/gothic"
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
	account, err := infra.AccountByEmail(ctx, a.db, user.Email)
	if err != nil && err != sql.ErrNoRows {
		handleServerError(ctx, err, w)
		return
	}

	// Create account
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
		if err := account.Insert(ctx, a.db); err != nil {
			handleServerError(ctx, err, w)
			return
		}
	}

	if err := a.sessionManager.Login(w, r, account.ID); err != nil {
		handleServerError(ctx, err, w)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("%s/profile", a.webAppURL), http.StatusFound)
}
