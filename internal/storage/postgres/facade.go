package postgres

import (
	"context"
	"time"

	"workshop-1/internal/domain"
	"workshop-1/internal/storage"
)

type Facade interface {
	CreateOrder(ctx context.Context, currentTime time.Time, order domain.Order) error
	GetOrder(ctx context.Context, orderID int, now time.Time) (domain.Order, error)
	GetReturn(ctx context.Context, orderID int, now time.Time) (domain.OrderReturn, error)
	GetOrders(ctx context.Context, orderIDs []int) ([]domain.Order, error)
	GetOrdersByUserID(ctx context.Context, userID, lastN int) ([]domain.Order, error)
	DeleteOrder(ctx context.Context, orderID int) error
	UpdateReceiveDates(ctx context.Context, currentTime time.Time, orderIDs []int) error
	CreateReturn(ctx context.Context, currentTime time.Time, userID, orderID int) error
	GetReturns(ctx context.Context, page, pageSize int) ([]domain.OrderReturn, error)
}

type storageFacade struct {
	txManager        TransactionManager
	storage          *Storage
	orderCache       storage.Cache[int, domain.Order]
	orderReturnCache storage.Cache[int, domain.OrderReturn]
}

func NewFacade(
	txManager TransactionManager,
	pgStorage *Storage,
	orderCache storage.Cache[int, domain.Order],
	orderReturnCache storage.Cache[int, domain.OrderReturn],
) Facade {
	if orderCache == nil {
		orderCache = storage.NoCache[int, domain.Order]{}
	}

	if orderReturnCache == nil {
		orderReturnCache = storage.NoCache[int, domain.OrderReturn]{}
	}

	return &storageFacade{
		txManager:        txManager,
		storage:          pgStorage,
		orderCache:       orderCache,
		orderReturnCache: orderReturnCache,
	}
}

func (s *storageFacade) CreateOrder(ctx context.Context, currentTime time.Time, order domain.Order) error {
	return s.storage.CreateOrder(ctx, currentTime, order)
}

func (s *storageFacade) GetOrder(ctx context.Context, orderID int, now time.Time) (domain.Order, error) {
	order, err := s.orderCache.Get(ctx, orderID, now)
	if err == nil {
		return order, nil
	}

	_ = s.orderCache.Set(ctx, orderID, order, now)

	return s.storage.GetOrder(ctx, orderID, now)
}

func (s *storageFacade) GetReturn(ctx context.Context, orderID int, now time.Time) (domain.OrderReturn, error) {
	orderReturn, err := s.orderReturnCache.Get(ctx, orderID, now)
	if err == nil {
		return orderReturn, nil
	}

	_ = s.orderReturnCache.Set(ctx, orderID, orderReturn, now)

	return s.storage.GetReturn(ctx, orderID, now)
}

func (s *storageFacade) GetOrders(ctx context.Context, orderIDs []int) ([]domain.Order, error) {
	var orders []domain.Order

	err := s.txManager.RunReadCommitted(ctx, func(ctx context.Context) error {
		o, err := s.storage.GetOrders(ctx, orderIDs)
		if err != nil {
			return err
		}

		orders = o
		return nil
	})

	return orders, err
}

func (s *storageFacade) GetOrdersByUserID(ctx context.Context, userID, lastN int) ([]domain.Order, error) {
	var orders []domain.Order

	err := s.txManager.RunReadCommitted(ctx, func(ctx context.Context) error {
		o, err := s.storage.GetOrdersByUserID(ctx, userID, lastN)
		if err != nil {
			return err
		}

		orders = o
		return nil
	})

	return orders, err
}

func (s *storageFacade) DeleteOrder(ctx context.Context, orderID int) error {
	return s.storage.DeleteOrder(ctx, orderID)
}

func (s *storageFacade) UpdateReceiveDates(ctx context.Context, currentTime time.Time, orderIDs []int) error {
	err := s.txManager.RunSerializable(ctx, func(ctx context.Context) error {
		return s.storage.UpdateReceiveDates(ctx, currentTime, orderIDs)
	})

	return err
}

func (s *storageFacade) CreateReturn(ctx context.Context, currentTime time.Time, userID, orderID int) error {
	return s.storage.CreateReturn(ctx, currentTime, userID, orderID)
}

func (s *storageFacade) GetReturns(ctx context.Context, page, pageSize int) ([]domain.OrderReturn, error) {
	var orderReturns []domain.OrderReturn

	err := s.txManager.RunReadCommitted(ctx, func(ctx context.Context) error {
		o, err := s.storage.GetReturns(ctx, page, pageSize)
		if err != nil {
			return err
		}

		orderReturns = o
		return nil
	})

	return orderReturns, err
}
