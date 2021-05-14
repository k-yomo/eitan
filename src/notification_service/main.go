package main

import (
	"cloud.google.com/go/datastore"
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/k-yomo/eitan/src/notification_service/config"
	"github.com/k-yomo/eitan/src/notification_service/internal/email"
	"github.com/k-yomo/eitan/src/pkg/event"
	"github.com/k-yomo/eitan/src/pkg/logging"
	"github.com/k-yomo/pm"
	"github.com/k-yomo/pm/middleware/logging/pm_zap"
	"github.com/k-yomo/pm/middleware/pm_autoack"
	"github.com/k-yomo/pm/middleware/pm_recovery"
	"github.com/pkg/errors"
	"github.com/sendgrid/sendgrid-go"
	"go.uber.org/zap"
	"log"
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

	logger, err := logging.NewLogger(!appConfig.IsDeployedEnv())
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
		),
	)
	defer pubsubSubscriber.Close()

	var emailClient email.Client
	if appConfig.IsDeployedEnv() {
		emailClient = email.NewSendgridEmailClient(sendgrid.NewSendClient(appConfig.GCPProjectID))
	} else {
		emailClient = email.NewNoopEmailClient()
	}

	h := NewNotificationHandler(dsClient, emailClient)

	if !appConfig.IsDeployedEnv() {
		if err := createTopicsAndSubs(pubsubClient); err != nil {
			logger.Fatal("create topics and subscription failed", zap.Error(err))
		}
	}
	err = pubsubSubscriber.HandleSubscriptionFunc("notification-service-account-registration-sub", h.HandleAccountRegistration)
	if err != nil {
		logger.Fatal("set subscription handler failed", zap.Error(err))
	}

	pubsubSubscriber.Run(ctxzap.ToContext(context.Background(), logger))
	defer pubsubSubscriber.Close()
	log.Printf("pubsub subscriber started running")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}

func createTopicsAndSubs(pubsubClient *pubsub.Client) error {
	ctx := context.Background()
	t := pubsubClient.Topic(event.AccountRegistrationTopicName)
	topicExist, err := t.Exists(ctx)
	if err != nil {
		return errors.Wrap(err, "check if topic exists failed")
	}
	if !topicExist {
		t, err = pubsubClient.CreateTopic(ctx, event.AccountRegistrationTopicName)
		if err != nil {
			return errors.Wrap(err, "create topic failed")
		}
	}
	subExist, err := pubsubClient.Subscription("notification-service-account-registration-sub").Exists(ctx)
	if err != nil {
		return errors.Wrap(err, "check if subscription exist failed")
	}
	if !subExist {
		c := pubsub.SubscriptionConfig{Topic: t, AckDeadline: 10 * time.Minute}
		if _, err := pubsubClient.CreateSubscription(ctx, "notification-service-account-registration-sub", c); err != nil {
			return errors.Wrap(err, "create subscription failed")
		}
	}
	return nil
}
