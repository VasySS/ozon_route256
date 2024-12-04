package main

import (
	"context"
	"log"
	"time"

	"workshop-1/config"
	"workshop-1/internal/domain"
	"workshop-1/internal/storage/postgres"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/brianvoe/gofakeit/v7/source"
	"github.com/jackc/pgx/v5/pgxpool"
)

const dataAmount = 10000

func main() {
	gofakeit.GlobalFaker = gofakeit.NewFaker(source.NewCrypto(), false)
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, config.PostgresURL)
	if err != nil {
		log.Fatal(err)
	}

	storage := newStorage(pool)
	generateTestData(ctx, storage)
}

func newStorage(pool *pgxpool.Pool) postgres.Facade {
	txManager := postgres.NewTxManager(pool)
	storage := postgres.NewStorage(txManager)

	return storage
}

func generateTestData(ctx context.Context, storage postgres.Facade) {
	orders := make([]domain.Order, 0, dataAmount)
	currentTime := time.Now().UTC()

	for i := 0; i < dataAmount; i++ {
		order := domain.Order{
			ID:          int(gofakeit.Int32()),
			UserID:      int(gofakeit.Int32()),
			ExpiryDate:  gofakeit.FutureDate().Truncate(time.Millisecond).UTC(),
			ReceiveDate: gofakeit.PastDate().Truncate(time.Millisecond).UTC(),
			Weight:      gofakeit.Float32Range(1, 100),
			Price:       gofakeit.Float32Range(1, 1000),
		}

		err := storage.CreateOrder(ctx, currentTime, order)
		if err != nil {
			log.Fatal(err)
		}

		orders = append(orders, order)
	}

	for _, order := range orders {
		err := storage.CreateReturn(ctx, currentTime, order.UserID, order.ID)
		if err != nil {
			log.Fatal(err)
		}
	}
}
