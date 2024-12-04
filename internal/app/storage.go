package app

import (
	"context"

	"workshop-1/config"
	"workshop-1/internal/app/logger"
	"workshop-1/internal/domain"
	"workshop-1/internal/storage"
	"workshop-1/internal/storage/postgres"
	"workshop-1/internal/storage/valkey"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewStorageFacade(
	ctx context.Context,
	closer *Closer,
	orderCache storage.Cache[int, domain.Order],
	orderReturnCache storage.Cache[int, domain.OrderReturn],
) postgres.Facade {
	pool, err := pgxpool.New(ctx, config.PostgresURL)
	if err != nil {
		logger.Fatal(err)
	}
	closer.Add(pool.Close)

	txManager := postgres.NewTxManager(pool)
	storage := postgres.NewStorage(txManager)

	return postgres.NewFacade(txManager, storage, orderCache, orderReturnCache)
}

func NewValkeyCache[K comparable, V any](ctx context.Context, closer *Closer) *valkey.Cache[K, V] {
	valkey, err := valkey.New[K, V](config.ValkeyURL, config.CacheTTL)
	if err != nil {
		logger.Fatal(err)
	}
	closer.Add(valkey.Close)

	return valkey
}
