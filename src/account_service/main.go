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
	"github.com/k-yomo/eitan/src/account_service/internal/config"
	"github.com/k-yomo/eitan/src/account_service/internal/sessionmanager"
	"github.com/k-yomo/eitan/src/internal/pb/eitan"
	"github.com/k-yomo/eitan/src/internal/sharedctx"
	"github.com/k-yomo/eitan/src/internal/tracing"
	"github.com/k-yomo/eitan/src/pkg/healthserver"
	"github.com/k-yomo/eitan/src/pkg/logging"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"
	"github.com/oklog/run"
	"github.com/pkg/errors"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	appConfig, err := config.NewAppConfig()
	if err != nil {
		panic(fmt.Sprintf("initialize app config failed: %v", err))
	}

	logger, err := logging.NewLogger(!appConfig.Env.IsDeployed())
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

	sessionManager, err := sessionmanager.NewSessionManager(appConfig, redisClient)
	if err != nil {
		logger.Fatal("initialize session manager failed", zap.Error(err))
	}

	pubsubClient, err := pubsub.NewClient(context.Background(), appConfig.GCPProjectID)
	if err != nil {
		logger.Fatal("initialize pubsub client failed", zap.Error(err))
	}
	defer pubsubClient.Close()

	if appConfig.Env.IsDeployed() {
		if err := tracing.InitTracer(); err != nil {
			logger.Fatal("set trace provider failed", zap.Error(err))
		}
	}

	r := newRouter(appConfig, logger)

	goth.UseProviders(
		google.New(appConfig.GoogleAuthClientKey, appConfig.GoogleAuthSecret, fmt.Sprintf("%s/auth/google/callback", appConfig.AppRootURL), "email", "profile"),
	)

	authHandler := NewAuthHandler(sessionManager, db, pubsubClient, appConfig.WebAppURL)
	r.HandleFunc("/auth/logout", authHandler.Logout).Methods("GET")
	r.HandleFunc("/auth/{provider}", authHandler.HandleOAuth).Methods("GET")
	r.HandleFunc("/auth/{provider}/callback", authHandler.HandleOAuthCallback).Methods("GET")

	grpcServer := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_recovery.UnaryServerInterceptor(),
			otelgrpc.UnaryServerInterceptor(),
			grpc_zap.UnaryServerInterceptor(logger),
			logging.NewUnaryServerInterceptor(appConfig.GCPProjectID),
			sharedctx.NewUnaryServerCurrentUserInterceptor(),
		),
	)

	eitan.RegisterAccountServiceServer(grpcServer, NewAccountServiceServer(db, sessionManager))
	grpc_health_v1.RegisterHealthServer(grpcServer, healthserver.NewHealthServer())
	reflection.Register(grpcServer)

	httpServer := &http.Server{Addr: fmt.Sprintf(":%d", appConfig.HTTPPort), Handler: r}

	var g run.Group
	quitChan := make(chan os.Signal, 1)
	signal.Notify(quitChan, syscall.SIGTERM, os.Interrupt)
	closeChan := make(chan struct{})

	g.Add(
		func() error {
			select {
			case sig := <-quitChan:
				logger.Info("Signal received, shutting down gracefully...", zap.Any("signal", sig))
			case <-closeChan:
			}
			return nil
		},
		func(err error) {
			close(closeChan)
		},
	)

	g.Add(
		func() error {
			if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
				return errors.Wrap(err, "serve rest server")
			}
			return nil
		},
		func(err error) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			if err := httpServer.Shutdown(ctx); err != nil {
				logger.Error("shutdown rest server failed", zap.Error(err))
			}
		},
	)

	g.Add(
		func() error {
			lis, err := net.Listen("tcp", fmt.Sprintf(":%d", appConfig.GRPCPort))
			if err != nil {
				return errors.Wrap(err, "listen to grpc server")
			}
			if err := grpcServer.Serve(lis); err != nil {
				return errors.Wrap(err, "serve grpc server")
			}
			return nil
		},
		func(err error) {
			grpcServer.GracefulStop()
		},
	)

	logger.Fatal(g.Run().Error())
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
	// r.Use(csrf.NewCSRFValidationMiddleware(appConfig.Env.IsDeployed()))
	return r
}
