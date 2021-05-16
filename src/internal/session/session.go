package session

import (
	"context"
)

const CookieSessionKey = "sid"

// SetSessionID sets session id to context
func SetSessionID(ctx context.Context, sid string) context.Context {
	return context.WithValue(ctx, CookieSessionKey, sid)
}

// GetSessionID extract session id from context
func GetSessionID(ctx context.Context) (string, bool) {
	switch sid := ctx.Value(CookieSessionKey).(type) {
	case string:
		return sid, true
	default:
		return "", false
	}
}
