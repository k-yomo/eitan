package sessionmanager

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/sessions"
	"github.com/k-yomo/eitan/src/account_service/config"
	"github.com/k-yomo/eitan/src/internal/session"
	"github.com/k-yomo/eitan/src/pkg/logging"
	"github.com/pkg/errors"
	"github.com/rbcervilla/redisstore/v8"
	"go.uber.org/zap"
	"net/http"
)

const sessionUserIDKey = "user_id"

type SessionManager interface {
	Authenticate(sid string) (userID string, err error)
	Login(w http.ResponseWriter, r *http.Request, userID string) error
	Logout(w http.ResponseWriter, r *http.Request) error
}

type sessionManagerImpl struct {
	redisStore *redisstore.RedisStore
}

type AuthenticatedUserInfo struct {
	ID           string
	Email        string
	DisplayName  string
	ScreenImgURL *string
}

func NewSessionManager(appConfig *config.AppConfig, redisClient *redis.Client) (SessionManager, error) {
	redisStore, err := redisstore.NewRedisStore(context.Background(), redisClient)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create redis store")
	}
	redisStore.Options(sessions.Options{
		Path:     "/",
		Domain:   appConfig.SessionCookieDomain,
		MaxAge:   60 * 60 * 24 * 365, // 1 year
		Secure:   appConfig.IsDeployedEnv(),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	redisStore.KeyPrefix("session:")
	return &sessionManagerImpl{
		redisStore: redisStore,
	}, nil
}

func (s *sessionManagerImpl) Authenticate(sid string) (userID string, err error) {
	// create pseudo http request to validate session
	r := http.Request{Header: map[string][]string{}}
	r.AddCookie(&http.Cookie{
		Name:  session.CookieSessionKey,
		Value: sid,
	})
	sess, err := s.redisStore.Get(&r, session.CookieSessionKey)
	if err != nil {
		return "", err
	}
	return sess.Values[sessionUserIDKey].(string), nil
}

func (s *sessionManagerImpl) Login(w http.ResponseWriter, r *http.Request, userID string) error {
	ctx := r.Context()
	sess, err := s.redisStore.New(r, session.CookieSessionKey)
	if err != nil {
		return err
	}
	sess.Values[sessionUserIDKey] = userID

	if err := sess.Save(r, w); err != nil {
		return err
	}
	logging.Logger(ctx).Info("Login Success", zap.String("userID", userID))
	return nil
}

func (s *sessionManagerImpl) Logout(w http.ResponseWriter, r *http.Request) error {
	sess, err := s.redisStore.Get(r, session.CookieSessionKey)
	if err != nil {
		return err
	}
	sess.Options.MaxAge = -1
	if err := sess.Save(r, w); err != nil {
		return err
	}

	return nil
}
