package session

import (
	"context"
)

const CookieSessionKey = "eitansid"

type sessionCtxKey struct{}

// SetSessionID sets session id to context
func SetSessionID(ctx context.Context, sid string) context.Context {
	return context.WithValue(ctx, sessionCtxKey{}, sid)
}

// GetSessionID extract session id from context
func GetSessionID(ctx context.Context) (string, bool) {
	switch sid := ctx.Value(sessionCtxKey{}).(type) {
	case string:
		return sid, true
	default:
		return "", false
	}
}
