package main

import (
	"cloud.google.com/go/pubsub"
	"context"
	"crypto/tls"
	"fmt"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-redis/redis/v8"
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
	"github.com/k-yomo/eitan/src/internal/tracing"
	"github.com/k-yomo/eitan/src/pkg/gqlopentelemetry"
	"github.com/k-yomo/eitan/src/pkg/logging"
	"github.com/k-yomo/eitan/src/pkg/tx"
	"github.com/k-yomo/pm"
	"github.com/k-yomo/pm/middleware/logging/pm_zap"
	"github.com/k-yomo/pm/middleware/pm_autoack"
	"github.com/k-yomo/pm/middleware/pm_recovery"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
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
		if err := tracing.InitTracer(); err != nil {
			logger.Fatal("set trace provider failed", zap.Error(err))
		}
	}

	db, err := sqlx.Connect(dbConfig.Driver, dbConfig.Dsn())
	if err != nil {
		logger.Fatal("initialize db failed", zap.Error(err))
	}
	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)
	defer db.Close()

	redisClient := redis.NewClient(&redis.Options{
		Addr: appConfig.RedisURL,
	})
	if _, err := redisClient.Ping(context.Background()).Result(); err != nil {
		logger.Fatal("initialize redis client failed", zap.Error(err))
	}

	pubsubClient, err := pubsub.NewClient(context.Background(), appConfig.GCPProjectID)
	if err != nil {
		logger.Fatal("initialize pubsub client failed", zap.Error(err))
	}
	defer pubsubClient.Close()

	pubsubSubscriber := pm.NewSubscriber(
		pubsubClient,
		pm.WithSubscriptionInterceptor(
			pm_recovery.SubscriptionInterceptor(),
			pm_zap.SubscriptionInterceptor(logger),
			pm_autoack.SubscriptionInterceptor(),
		),
	)
	defer pubsubSubscriber.Close()

	pubsubHandler := NewPubSubHandler(db)
	err = pubsubSubscriber.HandleSubscriptionFunc("eitan.account.user-registered", pubsubHandler.HandleUserRegisteredEvent)
	if err != nil {
		logger.Fatal("set pubsub subscription handler func failed", zap.Error(err))
	}

	accountServiceClient, closeAccountServiceClient := newAccountServiceClient(context.Background(), apiConfig.AccountServiceGRPCURL, appConfig.IsDeployedEnv())
	defer closeAccountServiceClient()

	gqlConfig := gql.Config{Resolvers: graph.NewResolver(db, tx.NewManager(db), accountServiceClient, redisClient)}
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

	pubsubSubscriber.Run(context.Background())
	defer pubsubSubscriber.Close()
	log.Printf("pubsub subscriber started running")

	log.Printf("server listening on port: %d", appConfig.Port)
	logger.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", appConfig.Port), r).Error())
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
	opts = append(opts, grpc.WithUnaryInterceptor(sharedctx.NewUnaryClientCurrentUserInterceptor(auth.GetUserID)))
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
	corsMiddleware := handlers.CORS(
		handlers.AllowedOrigins(appConfig.AllowedOrigins),
		handlers.AllowedHeaders([]string{"X-Requested-By", "Origin", "Authorization", "Accept", "Content-Type"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowCredentials(),
	)

	r.Use(corsMiddleware)
	r.Use(logging.NewMiddleware(appConfig.GCPProjectID, logger))
	r.Use(csrf.NewCSRFCheckMiddleware(appConfig.IsDeployedEnv()))
	r.Use(auth.NewSessionIDMiddleware())
	return r
}
