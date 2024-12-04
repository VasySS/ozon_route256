package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type txManagerKey struct{}

type QueryEngine interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type TransactionManager interface {
	RunReadCommitted(ctx context.Context, fn func(context.Context) error) error
	RunSerializable(ctx context.Context, fn func(context.Context) error) error
	GetQueryEngine(ctx context.Context) QueryEngine
}

type TxManager struct {
	pool *pgxpool.Pool
}

func NewTxManager(pool *pgxpool.Pool) *TxManager {
	return &TxManager{
		pool: pool,
	}
}

func (tm *TxManager) RunReadCommitted(ctx context.Context, fn func(context.Context) error) error {
	opts := pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadOnly,
	}

	return tm.beginFunc(ctx, opts, fn)
}

func (tm *TxManager) RunSerializable(ctx context.Context, fn func(context.Context) error) error {
	opts := pgx.TxOptions{
		IsoLevel:   pgx.Serializable,
		AccessMode: pgx.ReadWrite,
	}

	return tm.beginFunc(ctx, opts, fn)
}

func (tm *TxManager) beginFunc(
	ctx context.Context,
	opts pgx.TxOptions,
	fn func(txCtx context.Context) error,
) error {
	tx, err := tm.pool.BeginTx(ctx, opts)
	if err != nil {
		return err
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	ctx = context.WithValue(ctx, txManagerKey{}, tx)
	if err := fn(ctx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (tm *TxManager) GetQueryEngine(ctx context.Context) QueryEngine {
	// симуляция долгого запроса
	// time.Sleep(time.Second * 5)

	tx, ok := ctx.Value(txManagerKey{}).(QueryEngine)
	if ok && tx != nil {
		return tx
	}

	return tm.pool
}
