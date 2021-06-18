package main

import (
	"cloud.google.com/go/pubsub"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/k-yomo/eitan/src/account_service/internal/infra"
	"github.com/k-yomo/eitan/src/account_service/internal/sessionmanager"
	"github.com/k-yomo/eitan/src/internal/pb/eitan"
	"github.com/k-yomo/eitan/src/internal/pubsubevent"
	"github.com/k-yomo/eitan/src/pkg/clock"
	"github.com/k-yomo/eitan/src/pkg/logging"
	"github.com/k-yomo/eitan/src/pkg/randnum"
	"github.com/k-yomo/eitan/src/pkg/sqlutil"
	"github.com/k-yomo/eitan/src/pkg/tx"
	"github.com/k-yomo/eitan/src/pkg/uuid"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/proto"
	"net/http"
)

type authProvider string

const (
	authProviderEmail  authProvider = "email"
	authProviderGoogle authProvider = "google"
)

type AuthHandler struct {
	sessionManager sessionmanager.SessionManager
	db             *sqlx.DB
	txManager      tx.Manager
	pubsubClient   *pubsub.Client
	webAppURL      string
}

func NewAuthHandler(sessionManager sessionmanager.SessionManager, db *sqlx.DB, pubsubClient *pubsub.Client, webAppURL string) *AuthHandler {
	return &AuthHandler{
		sessionManager: sessionManager,
		db:             db,
		txManager:      tx.NewManager(db),
		pubsubClient:   pubsubClient,
		webAppURL:      webAppURL,
	}
}

func (a *AuthHandler) CreateEmailConfirmation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	params := struct {
		Email string `json:"email"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		handleClientError(ctx, err, w)
		return
	}

	now := clock.Now()
	emailConfirmation := infra.EmailConfirmation{
		Email:            params.Email,
		ConfirmationCode: randnum.RandNumString(6),
		CreatedAt:        now,
	}
	emailConfirmationCreatedEvent, m, err := newEmailConfirmationCreatedEvent(&emailConfirmation)
	if err != nil {
		handleServerError(ctx, err, w)
	}

	err = a.txManager.RunInTx(ctx, func(ctx context.Context) error {
		userProfile, err := infra.GetUserProfileByEmail(ctx, a.db, params.Email)
		if err != nil && err != sql.ErrNoRows {
			return err
		}
		if userProfile != nil {
			return errors.New("The email address is already used.")
		}

		if err := emailConfirmation.InsertOrUpdate(ctx, a.db); err != nil {
			return err
		}
		if err := emailConfirmationCreatedEvent.Insert(ctx, a.db); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		handleServerError(ctx, err, w)
	}

	if _, err := a.pubsubClient.Topic(emailConfirmationCreatedEvent.Topic).Publish(ctx, m).Get(ctx); err != nil {
		logging.Logger(ctx).Error(
			"publish EmailAuthTempRegisteredEvent failed",
			zap.Error(err),
			zap.Any("EmailAuthTempRegisteredEvent", emailConfirmationCreatedEvent),
		)
	}

	w.WriteHeader(200)
}

func (a *AuthHandler) SignUpWithEmail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	params := struct {
		DisplayName           string `json:"displayName"`
		Email                 string `json:"email"`
		EmailConfirmationCode string `json:"emailConfirmationCode"`
		Password              string `json:"password"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		handleClientError(ctx, err, w)
		return
	}

	now := clock.Now()
	err := a.txManager.RunInTx(ctx, func(ctx context.Context) error {
		emailConfirmation, err := infra.GetEmailConfirmationByEmailConfirmationCode(ctx, a.db, params.Email, params.EmailConfirmationCode)
		if err == sql.ErrNoRows {
			return err
		}
		if err != nil {
			return err
		}

		//nolint:staticcheck // SA1019
		emailConfirmation.ConfirmedAt = mysql.NullTime{Time: now, Valid: true}
		if err := emailConfirmation.Update(ctx, a.db); err != nil {
			return err
		}
		user, err := a.createUser(
			ctx,
			params.Email,
			params.DisplayName,
			nil,
			authProviderEmail,
		)
		if err != nil {
			return err
		}

		passwordDigest, err := newPasswordDigest(params.Password)
		if err != nil {
			return err
		}
		emailAuth := &infra.EmailAuth{
			UserID:         user.ID,
			Email:          params.Email,
			PasswordDigest: passwordDigest,
			CreatedAt:      now,
			UpdatedAt:      now,
		}
		if err := emailAuth.Insert(ctx, a.db); err != nil {
			return err
		}

		if err := a.sessionManager.Login(w, r, user.ID); err != nil {
			handleServerError(ctx, err, w)
			return err
		}
		return nil
	})
	if err != nil {
		handleServerError(ctx, err, w)
		return
	}

	w.WriteHeader(200)
}

func (a *AuthHandler) SignUpWithOAuth(w http.ResponseWriter, r *http.Request) {
	r.AddCookie(&http.Cookie{
		Name:     "oauth_for",
		Value:    "signup",
		Path:     "/",
		MaxAge:   0,
		Secure:   true,
		HttpOnly: true,
	})
	gothic.BeginAuthHandler(w, r)
}

func (a *AuthHandler) LoginWithEmail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	params := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		handleClientError(ctx, err, w)
		return
	}

	emailAuth, err := infra.GetEmailAuthByEmail(ctx, a.db, params.Email)
	if err != nil && err != sql.ErrNoRows {
		handleServerError(ctx, err, w)
		return
	}
	if emailAuth == nil {
		w.WriteHeader(401)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(emailAuth.PasswordDigest), []byte(params.Password))
	if err != nil {
		w.WriteHeader(401)
		return
	}

	if err := a.sessionManager.Login(w, r, emailAuth.UserID); err != nil {
		handleServerError(ctx, err, w)
		return
	}

	w.WriteHeader(200)
}

func (a *AuthHandler) LoginWithOAuth(w http.ResponseWriter, r *http.Request) {
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
		googleAuth, err := infra.GetGoogleAuthByGoogleID(ctx, a.db, gothUser.UserID)
		if googleAuth != nil {
			userID = googleAuth.UserID
			break
		}
		if err != nil && err != sql.ErrNoRows {
			handleServerError(ctx, err, w)
			return
		}
		// Create account if not exist
		user, err := a.createOauthUser(ctx, gothUser, authProviderGoogle)
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

	http.Redirect(w, r, a.webAppURL, http.StatusFound)
}

func (a *AuthHandler) createOauthUser(ctx context.Context, gothUser goth.User, provider authProvider) (*infra.User, error) {
	var screenImgURL *string
	if gothUser.AvatarURL != "" {
		screenImgURL = &gothUser.AvatarURL
	}
	var user *infra.User
	err := a.txManager.RunInTx(ctx, func(ctx context.Context) error {
		var err error
		user, err = a.createUser(ctx, gothUser.Email, gothUser.Name, screenImgURL, provider)
		if err != nil {
			return err
		}

		switch provider {
		case authProviderGoogle:
			googleAuth := infra.GoogleAuth{
				UserID:    user.ID,
				GoogleID:  gothUser.UserID,
				CreatedAt: user.CreatedAt,
				UpdatedAt: user.CreatedAt,
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

	return user, nil
}

func (a *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if err := a.sessionManager.Logout(w, r); err != nil {
		handleServerError(r.Context(), err, w)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("%s/login", a.webAppURL), http.StatusFound)
}

func (a *AuthHandler) createUser(
	ctx context.Context,
	email string,
	displayName string,
	screenImgURLStr *string,
	provider authProvider,
) (*infra.User, error) {
	logger := logging.Logger(ctx)

	now := clock.Now()
	user := &infra.User{
		ID:        uuid.Generate(),
		CreatedAt: now,
		UpdatedAt: now,
	}
	userProfile := &infra.UserProfile{
		UserID:       user.ID,
		Email:        email,
		DisplayName:  displayName,
		ScreenImgURL: sqlutil.PtrToNullString(screenImgURLStr),
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	userRegisteredEvent, m, err := newUserRegisteredEvent(userProfile, provider)
	if err != nil {
		return nil, err
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

		return nil
	})
	if err != nil {
		return nil, err
	}

	if _, err := a.pubsubClient.Topic(userRegisteredEvent.Topic).Publish(ctx, m).Get(ctx); err != nil {
		logger.Error("publish UserRegisteredEvent failed", zap.Error(err), zap.Any("UserRegisteredEvent", userRegisteredEvent))
	}

	return user, nil
}

func newEmailConfirmationCreatedEvent(emailConfirmation *infra.EmailConfirmation) (*infra.PubsubEvent, *pubsub.Message, error) {
	data, err := proto.Marshal(&eitan.EmailConfirmationCreatedEvent{
		Email:            emailConfirmation.Email,
		ConfirmationCode: emailConfirmation.ConfirmationCode,
	})
	if err != nil {
		return nil, nil, err
	}

	id := uuid.Generate()
	event := infra.PubsubEvent{
		ID: id,
		DeduplicateKey: sql.NullString{
			String: pubsubevent.NewDeduplicateKey(pubsubevent.EmailConfirmationCreatedTopicName, id),
			Valid:  true,
		},
		Topic:     pubsubevent.EmailConfirmationCreatedTopicName.String(),
		Data:      string(data),
		CreatedAt: clock.Now(),
	}
	m := pubsub.Message{
		Data: data,
		Attributes: pubsubevent.SetDeduplicateKey(
			map[string]string{},
			event.DeduplicateKey.String,
		),
	}
	return &event, &m, nil
}

func newUserRegisteredEvent(userProfile *infra.UserProfile, provider authProvider) (*infra.PubsubEvent, *pubsub.Message, error) {
	data, err := proto.Marshal(&eitan.UserRegisteredEvent{
		UserId:      userProfile.UserID,
		Provider:    string(provider),
		Email:       userProfile.Email,
		DisplayName: userProfile.DisplayName,
	})
	if err != nil {
		return nil, nil, err
	}
	event := infra.PubsubEvent{
		ID: uuid.Generate(),
		DeduplicateKey: sql.NullString{
			String: pubsubevent.NewDeduplicateKey(pubsubevent.UserRegisteredTopicName, userProfile.UserID),
			Valid:  true,
		},
		Topic:     pubsubevent.UserRegisteredTopicName.String(),
		Data:      string(data),
		CreatedAt: userProfile.CreatedAt,
	}
	m := pubsub.Message{
		Data: data,
		Attributes: pubsubevent.SetDeduplicateKey(
			map[string]string{},
			event.DeduplicateKey.String,
		),
	}
	return &event, &m, nil
}

func newPasswordDigest(password string) (string, error) {
	passwordDigest, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(passwordDigest), nil
}

func handleClientError(ctx context.Context, err error, w http.ResponseWriter) {
	logging.Logger(ctx).Warn(err.Error(), zap.Error(err))
	w.WriteHeader(400)
	_, _ = w.Write([]byte(`{"code":"400"}`))
}

func handleServerError(ctx context.Context, err error, w http.ResponseWriter) {
	logging.Logger(ctx).Error(err.Error(), zap.Error(err))
	w.WriteHeader(500)
	_, _ = w.Write([]byte(`{"status":"500"}`))
}
