package session

import "context"

const CookieSessionKey = "sid"

// GetSessionID extract session id from context
func SetSessionID(ctx context.Context, sid string) context.Context {
	return context.WithValue(ctx, CookieSessionKey, sid)
}

// GetSessionID extract session id from context
func GetSessionID(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(CookieSessionKey).(string)
	return userID, ok
}
