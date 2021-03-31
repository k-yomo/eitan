package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/k-yomo/eitan/src/account_service/infra"
	"github.com/k-yomo/eitan/src/account_service/internal/sessionmanager"
	"github.com/k-yomo/eitan/src/pkg/clock"
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
	webAppURL      string
}

func NewAuthHandler(sessionManager sessionmanager.SessionManager, db *sqlx.DB, webAppURL string) *AuthHandler {
	return &AuthHandler{sessionManager: sessionManager, db: db, webAppURL: webAppURL}
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

	http.Redirect(w, r, fmt.Sprintf("%s/account_settings", a.webAppURL), http.StatusFound)
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
