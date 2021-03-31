package logging

import (
	"context"
	"encoding/json"
	"github.com/99designs/gqlgen/graphql"
	"github.com/blendle/zapdriver"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"net/http"
	"strings"
)

type ctxKey int

const queryCtxKey ctxKey = 0

func NewLogger(isDev bool) (*zap.Logger, error) {
	var logger *zap.Logger
	var err error
	if isDev {
		config := zapdriver.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		logger, err = config.Build()
	} else {
		logger, err = zapdriver.NewProduction()
	}
	if err != nil {
		return nil, err
	}
	return logger, nil
}

// Logger returns logger from context if exist or from global variable
func Logger(ctx context.Context) *zap.Logger {
	logger := ctxzap.Extract(ctx)
	if query, ok := ctx.Value(queryCtxKey).(string); ok {
		logger = logger.With(zap.String("query", query))
	}
	return logger
}

// Middleware is set logger with request id to context
func NewMiddleware(gcpProjectID string, logger *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t := newTraceFromTraceContext(gcpProjectID, r.Header.Get("X-Cloud-Trace-Context"))
			zapFields := append(
				zapdriver.TraceContext(t.TraceID, t.SpanID, true, t.ProjectID),
				zap.String("ip", r.Header.Get("X-Forwarded-For")),
			)
			ctx := ctxzap.ToContext(r.Context(), logger.With(zapFields...))
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// UnaryServerInterceptor returns a new unary server interceptors that logs request.
func NewUnaryServerInterceptor(gcpProjectID string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		incomingMetadata, _ := metadata.FromIncomingContext(ctx)
		metadataCopy := incomingMetadata.Copy()
		_, spanCtx := otelgrpc.Extract(ctx, &metadataCopy)

		if spanCtx.IsValid() {
			t := traceInfo{ProjectID: gcpProjectID, TraceID: spanCtx.TraceID().String(), SpanID: spanCtx.SpanID().String()}
			ctx = ctxzap.ToContext(ctx, ctxzap.Extract(ctx).With(
				zapdriver.TraceContext(t.TraceID, t.SpanID, true, t.ProjectID)...,
			))
		}

		reqJsonBytes, _ := json.Marshal(req)
		Logger(ctx).Info(
			info.FullMethod,
			zap.ByteString("params", reqJsonBytes),
		)

		return handler(ctx, req)
	}
}

type traceInfo struct {
	ProjectID string
	TraceID   string
	SpanID    string
}

func newTraceFromTraceContext(projectID, traceContext string) traceInfo {
	t := traceInfo{ProjectID: projectID}
	if traceContext != "" {
		params := strings.Split(traceContext, "/")
		if len(params) >= 2 {
			t.TraceID = params[0]
			t.SpanID = params[1]
		}
	}
	return t
}

type GraphQLResponseInterceptor struct{}

var _ interface {
	graphql.HandlerExtension
	graphql.ResponseInterceptor
} = GraphQLResponseInterceptor{}

func (g GraphQLResponseInterceptor) ExtensionName() string {
	return "Logging"
}

func (g GraphQLResponseInterceptor) Validate(schema graphql.ExecutableSchema) error {
	return nil
}

func (g GraphQLResponseInterceptor) InterceptResponse(ctx context.Context, next graphql.ResponseHandler) *graphql.Response {
	oc := graphql.GetOperationContext(ctx)
	Logger(ctx).Info(oc.OperationName, zap.Any("variables", oc.Variables))
	return next(context.WithValue(ctx, queryCtxKey, oc.RawQuery))
}
