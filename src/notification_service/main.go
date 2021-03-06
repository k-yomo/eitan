package main

import (
	"cloud.google.com/go/datastore"
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"github.com/k-yomo/eitan/src/notification_service/internal/config"
	"github.com/k-yomo/eitan/src/notification_service/internal/email"
	"github.com/k-yomo/eitan/src/pkg/logging"
	"github.com/k-yomo/pm"
	"github.com/k-yomo/pm/middleware/logging/pm_zap"
	"github.com/k-yomo/pm/middleware/pm_autoack"
	"github.com/k-yomo/pm/middleware/pm_effectively_once"
	"github.com/k-yomo/pm/middleware/pm_recovery"
	"github.com/sendgrid/sendgrid-go"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"syscall"
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

	dsClient, err := datastore.NewClient(context.Background(), appConfig.GCPProjectID)
	if err != nil {
		logger.Fatal("initialize datastore client failed", zap.Error(err))
	}
	defer dsClient.Close()

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
			pm_effectively_once.SubscriptionInterceptor(pm_effectively_once.NewDatastoreMutexer("PubSubEvent", dsClient)),
		),
	)
	defer pubsubSubscriber.Close()

	var emailClient email.Client
	if appConfig.Env.IsDeployed() {
		emailClient = email.NewSendgridEmailClient(sendgrid.NewSendClient(appConfig.GCPProjectID))
	} else {
		emailClient = email.NewNoopEmailClient()
	}

	h := NewPubSubHandler(dsClient, emailClient, appConfig.WebAppURL)
	err = pubsubSubscriber.HandleSubscriptionFuncMap(map[string]pm.MessageHandler{
		"notification.account.user-registered":            h.HandleUserRegisteredEvent,
		"notification.account.email-confirmation-created": h.HandleEmailConfirmationCreatedEvent,
	})
	if err != nil {
		logger.Fatal("set subscription handler failed", zap.Error(err))
	}

	pubsubSubscriber.Run(context.Background())
	defer pubsubSubscriber.Close()
	log.Printf("pubsub subscriber started running")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}
