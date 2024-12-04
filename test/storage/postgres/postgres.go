package postgres

import (
	"context"

	"workshop-1/config"
	"workshop-1/internal/storage/postgres"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/brianvoe/gofakeit/v7/source"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/suite"
)

var (
	migrationsPath = "../../../migrations"
)

type PostgresStorageSuite struct {
	suite.Suite
	storage postgres.Facade
	pool    *pgxpool.Pool
	testTx  pgx.Tx
}

func (s *PostgresStorageSuite) SetupSuite() {
	ctx := context.Background()
	gofakeit.GlobalFaker = gofakeit.NewFaker(source.NewCrypto(), false)

	pool, err := pgxpool.New(ctx, config.PostgresURL)
	s.Require().NoError(err)

	tm := postgres.NewTxManager(pool)
	storage := postgres.NewStorage(tm)
	storageFacade := postgres.NewFacade(tm, storage, nil, nil)

	s.storage = storageFacade
	s.pool = pool

	s.gooseDown()
	s.gooseUp()
}

func (s *PostgresStorageSuite) TearDownSuite() {
	defer s.pool.Close()
	s.gooseDown()
}

func (s *PostgresStorageSuite) SetupTest() {
	tx, err := s.pool.BeginTx(context.Background(), pgx.TxOptions{})
	s.Require().NoError(err, "не удалось создать транзакцию для теста")

	s.testTx = tx
}

func (s *PostgresStorageSuite) TearDownTest() {
	if err := s.testTx.Rollback(context.Background()); err != nil && err != pgx.ErrTxClosed {
		s.T().Fatal(err)
	}
}

func (s *PostgresStorageSuite) gooseUp() {
	if err := goose.SetDialect("postgres"); err != nil {
		s.T().Fatal(err)
	}

	db := stdlib.OpenDBFromPool(s.pool)

	if err := goose.Up(db, migrationsPath); err != nil {
		s.T().Fatal(err)
	}

	if err := db.Close(); err != nil {
		s.T().Fatal(err)
	}
}

func (s *PostgresStorageSuite) gooseDown() {
	if err := goose.SetDialect("postgres"); err != nil {
		s.T().Fatal(err)
	}

	db := stdlib.OpenDBFromPool(s.pool)

	if err := goose.Down(db, migrationsPath); err != nil {
		s.T().Log(err)
	}

	if err := db.Close(); err != nil {
		s.T().Log(err)
	}
}
