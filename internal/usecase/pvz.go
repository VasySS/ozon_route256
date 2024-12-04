package usecase

import (
	"context"
	"fmt"
	"slices"
	"time"

	"workshop-1/internal/domain"
	"workshop-1/internal/domain/strategy"
	"workshop-1/internal/dto"
)

type Storage interface {
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

type PVZError struct {
	Msg string
}

func (e PVZError) Error() string {
	return e.Msg
}

// AcceptFromCourier принять заказ от курьера
func AcceptFromCourier(
	ctx context.Context,
	storage Storage,
	currentTime time.Time,
	orderDTO dto.CreateOrder,
	packaging strategy.Packaging,
) error {
	order, err := domain.NewOrder(
		orderDTO.ID,
		orderDTO.UserID,
		orderDTO.Weight,
		orderDTO.Price,
		orderDTO.ExpiryDate,
	)
	if err != nil {
		return err
	}

	if err := packaging.Apply(&order); err != nil {
		return err
	}

	if _, err := storage.GetOrder(ctx, order.ID, currentTime); err == nil {
		return PVZError{"заказ уже существует"}
	}

	return storage.CreateOrder(ctx, currentTime, order)
}

// ReturnToCourier вернуть заказ курьеру
func ReturnToCourier(ctx context.Context, storage Storage, currentTime time.Time, orderID int) error {
	// if order was returned by user
	if _, err := storage.GetReturn(ctx, orderID, currentTime); err == nil {
		if err := storage.DeleteOrder(ctx, orderID); err != nil {
			return err
		}

		return nil
	}

	order, err := storage.GetOrder(ctx, orderID, currentTime)
	if err != nil {
		return err
	}

	if err := order.CanReturnToCourier(currentTime); err != nil {
		return err
	}

	return storage.DeleteOrder(ctx, orderID)
}

// GiveToUser выдать заказы клиенту
func GiveToUser(ctx context.Context, storage Storage, currentTime time.Time, orderIDs []int) error {
	orders, err := storage.GetOrders(ctx, orderIDs)
	if err != nil {
		return err
	}

	userID := orders[0].UserID
	for _, order := range orders {
		if order.UserID != userID {
			return PVZError{"у заказов разные ID пользователя"}
		}
	}

	orderIDs = getAvailableOrders(ctx, storage, orders, currentTime)

	return storage.UpdateReceiveDates(ctx, currentTime, orderIDs)
}

func getAvailableOrders(ctx context.Context, storage Storage, orders []domain.Order, currentTime time.Time) []int {
	receivedOrders := make([]int, 0, len(orders))

	for _, order := range orders {
		// check that order was not returned by user earlier
		if _, err := storage.GetReturn(ctx, order.ID, currentTime); err == nil {
			fmt.Printf("заказ %d уже был возвращён клиентом ранее\n", order.ID)
			continue
		}

		// check that order has not expired
		if order.ExpiryDate.Before(time.Now()) {
			fmt.Printf("срок хранения заказа %d истёк, выдача невозможна\n", order.ID)
			continue
		}

		receivedOrders = append(receivedOrders, order.ID)
	}

	return receivedOrders
}

// UserOrders получить список заказов
func UserOrders(ctx context.Context, storage Storage, userID, lastN int, inPVZOnly bool) ([]domain.Order, error) {
	orders, err := storage.GetOrdersByUserID(ctx, userID, lastN)
	if err != nil {
		return nil, err
	}

	if inPVZOnly {
		orders = slices.DeleteFunc(orders, func(order domain.Order) bool {
			return !order.ReceiveDate.IsZero()
		})
	}

	return orders, nil
}

// AcceptUserReturn принять возврат от клиента
func AcceptUserReturn(ctx context.Context, storage Storage, currentTime time.Time, userID, orderID int) error {
	if _, err := storage.GetReturn(ctx, orderID, currentTime); err == nil {
		return PVZError{"заказ уже был возвращён курьеру"}
	}

	order, err := storage.GetOrder(ctx, orderID, currentTime)
	if err != nil {
		return err
	}

	if err := order.CanAcceptUserReturn(currentTime); err != nil {
		return err
	}

	if err := storage.CreateReturn(ctx, currentTime, userID, orderID); err != nil {
		return err
	}

	return nil
}

// UserReturns получить список возвратов
func UserReturns(ctx context.Context, storage Storage, page, pageSize int) ([]domain.OrderReturn, error) {
	returns, err := storage.GetReturns(ctx, page, pageSize)
	if err != nil {
		return nil, err
	}

	return returns, nil
}
