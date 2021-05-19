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
	"github.com/k-yomo/eitan/src/internal/pubsubevent"
	"github.com/k-yomo/eitan/src/pkg/clock"
	"github.com/k-yomo/eitan/src/pkg/logging"
	"github.com/k-yomo/eitan/src/pkg/tx"
	"github.com/k-yomo/eitan/src/pkg/uuid"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"go.uber.org/zap"
	"net/http"
)

type AuthHandler struct {
	sessionManager sessionmanager.SessionManager
	db             *sqlx.DB
	txManager      tx.Manager
	pubsubClient   *pubsub.Client
	webAppURL      string
}

func NewAuthHandler(sessionManager sessionmanager.SessionManager, db *sqlx.DB, pubsubClient *pubsub.Client, webAppURL string) *AuthHandler {
	return &AuthHandler{sessionManager: sessionManager, db: db, txManager: tx.NewManager(db), pubsubClient: pubsubClient, webAppURL: webAppURL}
}

func (a *AuthHandler) HandleOAuth(w http.ResponseWriter, r *http.Request) {
	gothic.BeginAuthHandler(w, r)
}

func (a *AuthHandler) HandleOAuthCallback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	gothUser, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		handleServerError(ctx, err, w)
		return
	}
	var userID string
	switch gothUser.Provider {
	case "google":
		googleAuth, err := infra.GoogleAuthByGoogleID(ctx, a.db, gothUser.UserID)
		if googleAuth != nil {
			userID = googleAuth.UserID
			break
		}

		if err != nil && err != sql.ErrNoRows {
			handleServerError(ctx, err, w)
			return
		}
		// Create account if not exist
		user, err := a.createOauthUser(ctx, gothUser)
		if err != nil {
			handleServerError(ctx, err, w)
			return
		}
		userID = user.ID
	}

	if err := a.sessionManager.Login(w, r, userID); err != nil {
		handleServerError(ctx, err, w)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("%s/user_settings", a.webAppURL), http.StatusFound)
}

func (a *AuthHandler) createOauthUser(ctx context.Context, gothUser goth.User) (*infra.User, error) {
	logger := logging.Logger(ctx)

	now := clock.Now()

	user := &infra.User{
		ID:        uuid.Generate(),
		CreatedAt: now,
		UpdatedAt: now,
	}

	var screenImgURL sql.NullString
	if gothUser.AvatarURL != "" {
		screenImgURL = sql.NullString{String: gothUser.AvatarURL, Valid: true}
	}
	userProfile := &infra.UserProfile{
		UserID:       user.ID,
		Email:        gothUser.Email,
		DisplayName:  gothUser.Name,
		ScreenImgURL: screenImgURL,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	userRegisteredEventData, err := proto.Marshal(&eitan.UserRegisteredEvent{
		UserId:      user.ID,
		Provider:    gothUser.Provider,
		Email:       userProfile.Email,
		DisplayName: userProfile.Email,
	})
	if err != nil {
		logger.Error("marshal UserRegisteredEvent failed", zap.Error(err))
	}
	userRegisteredEvent := infra.PubsubEvent{
		ID: uuid.Generate(),
		DeduplicateKey: sql.NullString{
			String: pubsubevent.NewDeduplicateKey(pubsubevent.UserRegisteredTopicName, user.ID),
			Valid:  true,
		},
		Topic:     pubsubevent.UserRegisteredTopicName.String(),
		Data:      string(userRegisteredEventData),
		CreatedAt: now,
	}

	err = a.txManager.RunInTx(ctx, func(ctx context.Context) error {
		if err := user.Insert(ctx, a.db); err != nil {
			return err
		}
		if err := userProfile.Insert(ctx, a.db); err != nil {
			return err
		}
		if err := userRegisteredEvent.Insert(ctx, a.db); err != nil {
			return err
		}
		switch gothUser.Provider {
		case "google":
			googleAuth := infra.GoogleAuth{
				UserID:    user.ID,
				GoogleID:  gothUser.UserID,
				CreatedAt: now,
				UpdatedAt: now,
			}
			if err := googleAuth.Insert(ctx, a.db); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	t := a.pubsubClient.Topic(userRegisteredEvent.Topic)
	m := pubsub.Message{
		Data: []byte(userRegisteredEvent.Data),
		Attributes: pubsubevent.SetDeduplicateKey(
			map[string]string{},
			userRegisteredEvent.DeduplicateKey.String,
		),
	}
	if _, err := t.Publish(ctx, &m).Get(ctx); err != nil {
		logger.Error("publish UserRegisteredEvent failed", zap.Error(err), zap.Any("UserRegisteredEvent", userRegisteredEvent))
	}

	return user, nil
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
