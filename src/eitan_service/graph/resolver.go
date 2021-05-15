package graph

import (
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/k-yomo/eitan/src/internal/pb/eitan"
	"github.com/k-yomo/eitan/src/pkg/tx"
)

//go:generate go run github.com/99designs/gqlgen
type Resolver struct {
	db                   *sqlx.DB
	txManager            tx.Manager
	accountServiceClient eitan.AccountServiceClient
	redisClient          *redis.Client
}

func NewResolver(
	db *sqlx.DB,
	txManager tx.Manager,
	accountServiceClient eitan.AccountServiceClient,
	redisClient *redis.Client,
) *Resolver {
	return &Resolver{
		db:                   db,
		txManager:            txManager,
		accountServiceClient: accountServiceClient,
		redisClient:          redisClient,
	}
}
