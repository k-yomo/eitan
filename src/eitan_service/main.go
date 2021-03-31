package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	texporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/k-yomo/eitan/src/eitan_service/config"
	"github.com/k-yomo/eitan/src/eitan_service/graph"
	gql "github.com/k-yomo/eitan/src/eitan_service/graph/generated"
	"github.com/k-yomo/eitan/src/eitan_service/internal/auth"
	"github.com/k-yomo/eitan/src/eitan_service/internal/csrf"
	"github.com/k-yomo/eitan/src/internal/pb/eitan"
	"github.com/k-yomo/eitan/src/internal/sharedctx"
	"github.com/k-yomo/eitan/src/pkg/gqlopentelemetry"
	"github.com/k-yomo/eitan/src/pkg/logging"
	"github.com/k-yomo/eitan/src/pkg/tx"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net/http"
	"time"
)

func main() {
	appConfig, err := config.NewAppConfig()
	if err != nil {
		panic(fmt.Sprintf("initialize app config failed: %v", err))
	}
	apiConfig, err := config.NewAPIConfig()
	if err != nil {
		panic(fmt.Sprintf("initialize api config failed: %v", err))
	}
	logger, err := logging.NewLogger(!appConfig.IsDeployedEnv())
	if err != nil {
		panic(fmt.Sprintf("initialize logger failed: %v", err))
	}
	dbConfig, err := config.NewDBConfig()
	if err != nil {
		logger.Fatal("initialize db config failed", zap.Error(err))
	}

	if appConfig.IsDeployedEnv() {
		initTracer(logger)
	}

	db, err := sqlx.Connect(dbConfig.Driver, dbConfig.Dsn())
	if err != nil {
		logger.Fatal("initialize db failed", zap.Error(err))
	}
	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)
	defer db.Close()

	accountServiceClient, closeAccountServiceClient := newAccountServiceClient(context.Background(), apiConfig.AccountServiceGRPCURL, appConfig.IsDeployedEnv())
	defer closeAccountServiceClient()

	gqlConfig := gql.Config{Resolvers: graph.NewResolver(db, tx.NewManager(db), accountServiceClient)}
	gqlConfig.Directives.HasRole = auth.NewHasRole(accountServiceClient)
	srv := handler.NewDefaultServer(gql.NewExecutableSchema(gqlConfig))
	srv.SetErrorPresenter(graph.NewErrorPresenter())
	srv.Use(gqlopentelemetry.Tracer{})
	srv.Use(logging.GraphQLResponseInterceptor{})

	r := newRouter(appConfig, logger)
	r.Handle("/query", srv)
	if appConfig.Env == config.Local {
		r.Handle("/", playground.Handler("GraphQL playground", "/query"))
	}

	log.Printf("server listening on port: %d", appConfig.Port)
	logger.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", appConfig.Port), withCors(appConfig.AllowedOrigins)(r)).Error())
}

func newAccountServiceClient(ctx context.Context, accountServiceGRPCURL string, isDeployedEnv bool) (client eitan.AccountServiceClient, closeConn func()) {
	conn, err := grpc.DialContext(ctx, accountServiceGRPCURL, grpcOptions(isDeployedEnv)...)
	if err != nil {
		logging.Logger(ctx).Fatal("initialize payments api connection failed", zap.Error(err))
	}
	return eitan.NewAccountServiceClient(conn), func() {
		conn.Close()
	}
}

func grpcOptions(isDeployedEnv bool) []grpc.DialOption {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithUnaryInterceptor(sharedctx.NewUnaryClientCurrentAccountInterceptor(auth.GetAccountID)))
	if isDeployedEnv {
		creds := credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})
		gapsTraceInterceptor := otelgrpc.UnaryClientInterceptor()
		opts = append(
			opts,
			grpc.WithUnaryInterceptor(gapsTraceInterceptor),
			grpc.WithTransportCredentials(creds),
		)
	} else {
		opts = append(opts, grpc.WithInsecure())
	}
	return opts
}

func newRouter(appConfig *config.AppConfig, logger *zap.Logger) *mux.Router {
	r := mux.NewRouter()
	r.Use(logging.NewMiddleware(appConfig.GCPProjectID, logger))
	r.Use(csrf.NewCSRFCheckMiddleware(appConfig.IsDeployedEnv()))
	r.Use(auth.NewSessionIDMiddleware())
	return r
}

func withCors(allowedOrigins []string) func(http.Handler) http.Handler {
	return handlers.CORS(
		handlers.AllowedOrigins(allowedOrigins),
		handlers.AllowedHeaders([]string{"X-Requested-By", "Origin", "Authorization", "Accept", "Content-Type"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowCredentials(),
	)
}

func initTracer(logger *zap.Logger) {
	exporter, err := texporter.NewExporter()
	if err != nil {
		logger.Fatal("initialize exporter failed", zap.Error(err))
	}
	tp := sdktrace.NewTracerProvider(sdktrace.WithBatcher(exporter))
	if err != nil {
		logger.Fatal("initialize provider failed", zap.Error(err))
	}
	otel.SetTracerProvider(tp)
}
