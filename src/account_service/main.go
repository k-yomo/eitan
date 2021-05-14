package main

import (
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/jmoiron/sqlx"
	"github.com/k-yomo/eitan/src/account_service/config"
	"github.com/k-yomo/eitan/src/account_service/internal/sessionmanager"
	"github.com/k-yomo/eitan/src/internal/pb/eitan"
	"github.com/k-yomo/eitan/src/internal/sharedctx"
	"github.com/k-yomo/eitan/src/internal/tracing"
	"github.com/k-yomo/eitan/src/pkg/csrf"
	"github.com/k-yomo/eitan/src/pkg/logging"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"net/http"
)

func main() {
	appConfig, err := config.NewAppConfig()
	if err != nil {
		panic(fmt.Sprintf("initialize app config failed: %v", err))
	}

	logger, err := logging.NewLogger(!appConfig.IsDeployedEnv())
	if err != nil {
		panic(fmt.Sprintf("initialize logger failed: %v", err))
	}

	dbConfig, err := config.NewDBConfig()
	if err != nil {
		logger.Fatal("initialize db config failed", zap.Error(err))
	}
	db, err := sqlx.Connect(dbConfig.Driver, dbConfig.Dsn())
	if err != nil {
		logger.Fatal("initialize db failed", zap.Error(err))
	}

	redisClient := redis.NewClient(&redis.Options{Addr: appConfig.RedisURL})

	sessionManager := sessionmanager.NewSessionManager(appConfig, redisClient)

	pubsubClient, err := pubsub.NewClient(context.Background(), appConfig.GCPProjectID)
	if err != nil {
		logger.Fatal("initialize pubsub client failed", zap.Error(err))
	}
	defer pubsubClient.Close()

	if appConfig.IsDeployedEnv() {
		if err := tracing.InitTracer(); err != nil {
			logger.Fatal("set trace provider failed", zap.Error(err))
		}
	}

	r := newRouter(appConfig, logger)

	// healthcheck
	r.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"200"}`))
	}).Methods("GET")

	goth.UseProviders(
		google.New(appConfig.GoogleAuthClientKey, appConfig.GoogleAuthSecret, fmt.Sprintf("%s/auth/google/callback", appConfig.AppRootURL), "email", "profile"),
	)

	authHandler := NewAuthHandler(sessionManager, db, pubsubClient, appConfig.WebAppURL)
	r.HandleFunc("/auth/logout", authHandler.Logout).Methods("GET")
	r.HandleFunc("/auth/{provider}", authHandler.HandleOAuth).Methods("GET")
	r.HandleFunc("/auth/{provider}/callback", authHandler.HandleOAuthCallback).Methods("GET")

	s := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_recovery.UnaryServerInterceptor(),
			otelgrpc.UnaryServerInterceptor(),
			grpc_zap.UnaryServerInterceptor(logger),
			logging.NewUnaryServerInterceptor(appConfig.GCPProjectID),
			sharedctx.NewUnaryServerCurrentAccountInterceptor(),
		),
	)

	eitan.RegisterAccountServiceServer(s, NewAccountServiceServer(db, sessionManager))
	reflection.Register(s)

	eg := errgroup.Group{}
	eg.Go(func() error {
		fmt.Println("Rest Server listening on port:", appConfig.RestPort)
		return http.ListenAndServe(fmt.Sprintf(":%d", appConfig.RestPort), r)
	})
	eg.Go(func() error {
		fmt.Println("GRPC Server listening on port:", appConfig.GRPCPort)
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", appConfig.GRPCPort))
		if err != nil {
			return err
		}
		return s.Serve(lis)
	})
	logger.Fatal(eg.Wait().Error())
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
	r.Use(csrf.NewCSRFValidationMiddleware(appConfig.IsDeployedEnv()))
	return r
}
