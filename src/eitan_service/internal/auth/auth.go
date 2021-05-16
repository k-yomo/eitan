package auth

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
	"github.com/k-yomo/eitan/src/eitan_service/graph/model"
	"github.com/k-yomo/eitan/src/eitan_service/internal/customerror"
	"github.com/k-yomo/eitan/src/internal/pb/eitan"
	"github.com/k-yomo/eitan/src/internal/session"
	"github.com/k-yomo/eitan/src/pkg/logging"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

func NewSessionIDMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			c, err := r.Cookie(session.CookieSessionKey)
			if err == nil {
				r = r.WithContext(session.SetSessionID(ctx, c.Value))
			}

			next.ServeHTTP(w, r)
		})
	}
}

type authenticatedUserIDKey struct {}

// NewHasRole returns a function to authenticate user based on role
// role is defined at defs/graphql/schema.graphql
func NewHasRole(accountServiceClient eitan.AccountServiceClient) func(ctx context.Context, obj interface{}, next graphql.Resolver, role model.Role) (res interface{}, err error) {
	return func(ctx context.Context, obj interface{}, next graphql.Resolver, role model.Role) (res interface{}, err error) {
		switch role {
		case model.RoleUser:
			sid, ok := session.GetSessionID(ctx)
			if !ok {
				return nil, customerror.NewErrUnauthenticated()
			}
			res, err := accountServiceClient.Authenticate(ctx, &eitan.AuthenticateRequest{
				SessionId: sid,
			})
			if status.Code(err) == codes.Unauthenticated {
				return nil, customerror.NewErrUnauthenticated()
			} else if err != nil {
				return nil, customerror.New(err, customerror.ErrInternal)
			}
			logging.AddFields(ctx, zap.String("userID", res.UserProfile.UserId))
			ctx = context.WithValue(ctx, authenticatedUserIDKey{}, res.UserProfile.UserId)
		}
		return next(ctx)
	}
}

// GetUserID extract user id from context
func GetUserID(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(authenticatedUserIDKey{}).(string)
	return userID, ok
}
