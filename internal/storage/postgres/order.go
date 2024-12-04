package postgres

import (
	"context"
	"math"
	"time"

	"workshop-1/internal/domain"

	"github.com/georgysavva/scany/v2/pgxscan"
)

func (s *Storage) CreateOrder(ctx context.Context, currentTime time.Time, order domain.Order) error {
	tx := s.txManager.GetQueryEngine(ctx)

	query := `
		INSERT INTO user_order (id, created_at, updated_at, user_id, weight, price, expiry_date)
		VALUES ($1, $2, $2, $3, $4, $5, $6)
		RETURNING id
	`
	_, err := tx.Exec(ctx, query,
		order.ID, currentTime, order.UserID, order.Weight, order.Price, order.ExpiryDate)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetOrder(ctx context.Context, id int, now time.Time) (domain.Order, error) {
	tx := s.txManager.GetQueryEngine(ctx)

	query := `
		SELECT id, user_id, weight, price, expiry_date, receive_date
		FROM user_order
		WHERE id = $1
	`

	var order domain.Order
	if err := pgxscan.Get(ctx, tx, &order, query, id); err != nil {
		return domain.Order{}, err
	}

	return order, nil
}

func (s *Storage) GetOrders(ctx context.Context, ids []int) ([]domain.Order, error) {
	tx := s.txManager.GetQueryEngine(ctx)

	query := `
		SELECT id, user_id, weight, price, expiry_date, receive_date
		FROM user_order
		WHERE id = ANY($1)
	`

	var orders []domain.Order
	if err := pgxscan.Select(ctx, tx, &orders, query, ids); err != nil {
		return nil, err
	}

	return orders, nil
}

func (s *Storage) GetOrdersByUserID(ctx context.Context, userID, lastN int) ([]domain.Order, error) {
	if lastN == 0 {
		lastN = math.MaxInt
	}

	tx := s.txManager.GetQueryEngine(ctx)

	query := `
		SELECT id, user_id, weight, price, expiry_date, receive_date
		FROM user_order
		WHERE user_id = $1
		ORDER BY id DESC
		LIMIT $2
	`

	var orders []domain.Order
	if err := pgxscan.Select(ctx, tx, &orders, query, userID, lastN); err != nil {
		return nil, err
	}

	return orders, nil
}

func (s *Storage) DeleteOrder(ctx context.Context, id int) error {
	tx := s.txManager.GetQueryEngine(ctx)

	query := `
		DELETE FROM user_order
		WHERE id = $1
	`

	if _, err := tx.Exec(ctx, query, id); err != nil {
		return err
	}

	return nil
}

func (s *Storage) UpdateReceiveDates(ctx context.Context, currentTime time.Time, orderIDs []int) error {
	tx := s.txManager.GetQueryEngine(ctx)

	query := `
		UPDATE user_order
		SET receive_date = $1, updated_at = $1
		WHERE id = ANY($2)
	`

	if _, err := tx.Exec(ctx, query, currentTime, orderIDs); err != nil {
		return err
	}

	return nil
}
