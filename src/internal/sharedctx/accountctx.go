package sharedctx

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const GRPCCtxUserIDKey = "userID"

type currentUserIDKey struct{}

func NewUnaryClientCurrentUserInterceptor(getUserID func(ctx context.Context) (string, bool)) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if userID, ok := getUserID(ctx); ok {
			ctx = metadata.AppendToOutgoingContext(ctx, GRPCCtxUserIDKey, userID)
		}
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func NewUnaryServerCurrentUserInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			if v := md.Get(GRPCCtxUserIDKey); len(v) > 0 {
				ctx = context.WithValue(ctx, currentUserIDKey{}, v[0])
			}
		}

		return handler(ctx, req)
	}
}

// GetUserID extract user id from context
func GetUserID(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(currentUserIDKey{}).(string)
	return userID, ok
}