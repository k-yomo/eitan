package sessionmanager

import (
	"github.com/gorilla/sessions"
	"github.com/k-yomo/eitan/src/internal/session"
	"github.com/k-yomo/eitan/src/pkg/logging"
	"go.uber.org/zap"
	"net/http"
)

const sessionAccountIDKey = "account_id"

type SessionManager interface {
	Authenticate(sid string) (accountID string, err error)
	Login(w http.ResponseWriter, r *http.Request, accountID string) error
	Logout(w http.ResponseWriter, r *http.Request) error
}

type sessionManagerImpl struct {
	cookieStore         *sessions.CookieStore
}

type AuthenticatedUserInfo struct {
	ID           string
	Email        string
	DisplayName  string
	ScreenImgURL *string
}

func NewSessionManager(cookieStore *sessions.CookieStore) SessionManager {
	return &sessionManagerImpl{
		cookieStore:         cookieStore,
	}
}

func (s *sessionManagerImpl) Authenticate(sid string) (accountID string, err error) {
	// create pseudo http request to validate session
	r := http.Request{Header: map[string][]string{}}
	r.AddCookie(&http.Cookie{
		Name:  session.CookieSessionKey,
		Value: sid,
	})
	sess, err := s.cookieStore.Get(&r, session.CookieSessionKey)
	if err != nil {
		return "", err
	}
	return sess.Values[sessionAccountIDKey].(string), nil
}

func (s *sessionManagerImpl) Login(w http.ResponseWriter, r *http.Request, accountID string) error {
	ctx := r.Context()
	sess := sessions.NewSession(s.cookieStore, session.CookieSessionKey)
	sess, err := s.cookieStore.New(r, session.CookieSessionKey)
	if err != nil {
		return err
	}
	sess.Values[sessionAccountIDKey] = accountID

	if err := sess.Save(r, w); err != nil {
		return err
	}
	logging.Logger(ctx).Info("Login Success", zap.String("accountID", accountID))
	return nil
}

func (s *sessionManagerImpl) Logout(w http.ResponseWriter, r *http.Request) error {
	sess, err := s.cookieStore.Get(r, session.CookieSessionKey)
	if err != nil {
		return err
	}
	sess.Options.MaxAge = -1
	if err := sess.Save(r, w); err != nil {
		return err
	}

	return nil
}
