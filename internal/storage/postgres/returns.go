package postgres

import (
	"context"
	"time"

	"workshop-1/internal/domain"

	"github.com/georgysavva/scany/v2/pgxscan"
)

func (s *Storage) CreateReturn(ctx context.Context, currentTime time.Time, userID, orderID int) error {
	tx := s.txManager.GetQueryEngine(ctx)

	query := `
		INSERT INTO order_return (created_at, updated_at, user_id, order_id)
		VALUES ($1, $1, $2, $3)
	`

	_, err := tx.Exec(ctx, query, currentTime, userID, orderID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetReturn(ctx context.Context, orderID int, now time.Time) (domain.OrderReturn, error) {
	tx := s.txManager.GetQueryEngine(ctx)

	query := `
		SELECT user_id, order_id
		FROM order_return
		WHERE order_id = $1
	`

	var orderReturn domain.OrderReturn
	if err := pgxscan.Get(ctx, tx, &orderReturn, query, orderID); err != nil {
		return domain.OrderReturn{}, err
	}

	return orderReturn, nil
}

func (s *Storage) GetReturns(ctx context.Context, page, pageSize int) ([]domain.OrderReturn, error) {
	tx := s.txManager.GetQueryEngine(ctx)
	offset := (page - 1) * pageSize

	query := `
		SELECT user_id, order_id
		FROM order_return
		LIMIT $1
		OFFSET $2
	`

	var orderReturns []domain.OrderReturn
	if err := pgxscan.Select(ctx, tx, &orderReturns, query, pageSize, offset); err != nil {
		return nil, err
	}

	return orderReturns, nil
}
