package pubsubevent

import (
	"database/sql"
	"time"
)

// PubSubEvent represents a row from 'pubsub_events'.
type PubSubEvent struct {
	ID             string         `db:"id"`
	DeduplicateKey sql.NullString `db:"deduplicate_key"`
	Topic          string         `db:"topic"`
	Data           string         `db:"data"`
	IsPublished    bool           `db:"is_published"`
	PublishedAt    sql.NullTime   `db:"published_at"`
	CreatedAt      time.Time      `db:"created_at"`
}
