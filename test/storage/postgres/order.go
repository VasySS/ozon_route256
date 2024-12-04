package postgres

import (
	"context"
	"time"

	"workshop-1/internal/domain"

	"github.com/brianvoe/gofakeit/v7"
)

func (s *PostgresStorageSuite) TestOrder() {
	ctx := context.Background()
	currentTime := time.Now().UTC()

	orderInput := domain.Order{
		ID:          int(gofakeit.Int32()),
		UserID:      int(gofakeit.Int32()),
		ExpiryDate:  gofakeit.FutureDate().Truncate(time.Millisecond).UTC(),
		ReceiveDate: time.Time{},
		Weight:      gofakeit.Float32Range(1, 100),
		Price:       gofakeit.Float32Range(1, 1000),
	}

	s.Run("создание заказа", func() {
		err := s.storage.CreateOrder(ctx, currentTime, orderInput)
		s.Require().NoError(err)
	})

	s.Run("получение заказа", func() {
		order, err := s.storage.GetOrder(ctx, orderInput.ID, currentTime)
		s.Require().NoError(err)
		s.Require().Equal(order, orderInput)
	})

	s.Run("удаление заказа", func() {
		err := s.storage.DeleteOrder(ctx, orderInput.ID)
		s.Require().NoError(err)
	})
}

func (s *PostgresStorageSuite) TestOrders() {
	ctx := context.Background()
	currentTime := time.Now().UTC()

	orders := make([]domain.Order, 0, 10)
	orderIDs := make([]int, 0, 10)
	for i := 0; i < 10; i++ {
		order := domain.Order{
			ID:          int(gofakeit.Int32()),
			UserID:      1,
			ExpiryDate:  gofakeit.FutureDate().Truncate(time.Millisecond).UTC(),
			ReceiveDate: time.Time{},
			Weight:      gofakeit.Float32Range(1, 100),
			Price:       gofakeit.Float32Range(1, 1000),
		}

		err := s.storage.CreateOrder(ctx, currentTime, order)
		s.Require().NoError(err)

		orders = append(orders, order)
		orderIDs = append(orderIDs, order.ID)
	}

	s.Run("получение заказов по ID", func() {
		ordersGet, err := s.storage.GetOrders(ctx, orderIDs)
		s.Require().NoError(err)
		s.Require().ElementsMatch(orders, ordersGet)
	})

	s.Run("получение заказов по ID пользователя", func() {
		ordersGet, err := s.storage.GetOrdersByUserID(ctx, 1, 10)

		s.Require().NoError(err)
		s.Require().ElementsMatch(orders, ordersGet)
	})

	s.Run("обновление дат доставки", func() {
		err := s.storage.UpdateReceiveDates(ctx, currentTime, orderIDs)
		s.Require().NoError(err)
	})
}
