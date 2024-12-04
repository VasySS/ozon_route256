package postgres

import (
	"context"
	"testing"
	"time"

	"workshop-1/config"
	"workshop-1/internal/domain"
	"workshop-1/internal/storage/inmem"
	"workshop-1/internal/storage/valkey"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
)

func BenchmarkGetOrder(b *testing.B) {
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, config.PostgresURL)
	if err != nil {
		b.Fatal(err)
	}
	defer pool.Close()

	valkeyCache, err := valkey.New[int, domain.Order](config.ValkeyURL, time.Minute*5)
	if err != nil {
		b.Fatal(err)
	}
	defer valkeyCache.Close()

	inmemCache := inmem.New[int, domain.Order](config.CacheTTL, inmem.LFU, config.CacheCapacity)

	txManager := NewTxManager(pool)
	pgStorage := NewStorage(txManager)

	facadeWithCache := NewFacade(txManager, pgStorage, valkeyCache, nil)
	facadeWithInmemCache := NewFacade(txManager, pgStorage, inmemCache, nil)
	facadeWithoutCache := NewFacade(txManager, pgStorage, nil, nil)

	order := domain.Order{
		ID:         1337,
		UserID:     556677,
		Weight:     12.228,
		Price:      15.99,
		ExpiryDate: time.Date(2099, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	if err := pgStorage.CreateOrder(ctx, time.Now(), order); err != nil {
		b.Fatal(err)
	}
	defer pgStorage.DeleteOrder(ctx, order.ID)

	// искусственная доп нагрузка
	// for i := 0; i < 10000; i++ {
	// 	go func() {
	// 		for {
	// 			_, _ = pgStorage.GetOrder(ctx, order.ID)
	// 		}
	// 	}()
	// }

	b.Run("WithoutCache", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := facadeWithoutCache.GetOrder(ctx, order.ID, time.Now().UTC())
			assert.NoError(b, err)
		}
	})

	b.Run("WithValkeyCache", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := facadeWithCache.GetOrder(ctx, order.ID, time.Now().UTC())
			assert.NoError(b, err)
		}
	})

	b.Run("WithInmemCache", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := facadeWithInmemCache.GetOrder(ctx, order.ID, time.Now().UTC())
			assert.NoError(b, err)
		}
	})
}
