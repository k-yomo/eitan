package main

import (
	"cloud.google.com/go/pubsub"
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/jmoiron/sqlx"
	"github.com/k-yomo/eitan/src/internal/pubsubevent"
	"github.com/k-yomo/eitan/src/pkg/clock"
	"github.com/k-yomo/eitan/src/pkg/logging"
	"go.uber.org/zap"
	"sync"
)

func main() {
	logger, err := logging.NewLogger(false)
	if err != nil {
		panic(fmt.Sprintf("initialize logger failed: %v", err))
	}

	ctx := ctxzap.ToContext(context.Background(), logger)

	config, err := newConfig()
	if err != nil {
		logger.Fatal("failed to load config", zap.Error(err))
	}

	db, err := sqlx.Connect(config.DBDriver, config.dbDsn())
	if err != nil {
		logger.Fatal("initialize db failed", zap.Error(err))
	}

	pubsubClient, err := pubsub.NewClient(ctx, config.GCPProjectID)
	if err != nil {
		logger.Fatal("initialize pubsub client failed", zap.Error(err))
	}

	if err := publishMessages(ctx, db, pubsubClient); err != nil {
		logger.Fatal("publish messages failed", zap.Error(err))
	}
}

func publishMessages(ctx context.Context, db *sqlx.DB, pubsubClient *pubsub.Client) error {
	logger := logging.Logger(ctx)
	logger.Info("[START] publishMessages")

	events, err := getUnpublishedEvents(ctx, db)
	if err != nil {
		return err
	}
	logger.Info(fmt.Sprintf("Fetch %d pubsub events", len(events)))

	wg := sync.WaitGroup{}
	for _, e := range events {
		wg.Add(1)
		e := e
		attributes := make(map[string]string)
		if e.DeduplicateKey.Valid {
			attributes = pubsubevent.SetDeduplicateKey(attributes, e.DeduplicateKey.String)
		}
		res := pubsubClient.Topic(e.Topic).Publish(ctx, &pubsub.Message{
			Data:       []byte(e.Data),
			Attributes: attributes,
		})
		go func() {
			defer wg.Done()
			if _, err := res.Get(ctx); err != nil {
				logging.Logger(ctx).Error("publish event failed", zap.Error(err), zap.Any("event", e))
				return
			}
			if err := updateEventToPublished(ctx, db, e.ID); err != nil {
				logging.Logger(ctx).Error("update pubsub_event to published failed", zap.Error(err), zap.Any("event", e))
			}
		}()
	}
	wg.Wait()

	logger.Info("[END] publishMessages")
	return nil
}

func getUnpublishedEvents(ctx context.Context, db *sqlx.DB) ([]*pubsubevent.PubSubEvent, error) {
	const sqlstr = `
SELECT 
	id,
	deduplicate_key,
	topic,
	data,
	is_published,
	created_at,
	published_at
FROM pubsub_events
WHERE is_published = FALSE
`

	var events []*pubsubevent.PubSubEvent
	err := db.SelectContext(ctx, &events, sqlstr)
	if err != nil {
		return nil, err
	}
	return events, nil
}

func updateEventToPublished(ctx context.Context, db *sqlx.DB, eventID string) error {
	_, err := db.NamedExecContext(
		ctx,
		`UPDATE pubsub_events SET is_published = :is_published, published_at = :published_at WHERE id = :id`,
		map[string]interface{}{
			"id":           eventID,
			"is_published": true,
			"published_at": sql.NullTime{Time: clock.Now(), Valid: true},
		},
	)
	return err
}
