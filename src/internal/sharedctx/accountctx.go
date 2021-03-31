package sharedctx

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const GRPCCtxAccountIDKey = "accountID"

type currentAccountIDCtxKey int

var currentAccountAccountIDKey currentAccountIDCtxKey

func NewUnaryClientCurrentAccountInterceptor(getAccountID func(ctx context.Context) (string, bool)) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if accountID, ok := getAccountID(ctx); ok {
			ctx = metadata.AppendToOutgoingContext(ctx, GRPCCtxAccountIDKey, accountID)
		}
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func NewUnaryServerCurrentAccountInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			if v := md.Get(GRPCCtxAccountIDKey); len(v) > 0 {
				ctx = context.WithValue(ctx, currentAccountAccountIDKey, v[0])
			}
		}

		return handler(ctx, req)
	}
}

// GetAccountID extract account id from context
func GetAccountID(ctx context.Context) (string, bool) {
	accountID, ok := ctx.Value(currentAccountAccountIDKey).(string)
	return accountID, ok
}