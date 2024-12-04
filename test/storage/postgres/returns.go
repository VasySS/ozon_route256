package postgres

import (
	"context"
	"time"

	"workshop-1/internal/domain"

	"github.com/brianvoe/gofakeit/v7"
)

func (s *PostgresStorageSuite) TestOrderReturn() {
	ctx := context.Background()
	currentTime := time.Now().UTC()

	orderInput := domain.Order{
		ID:          int(gofakeit.Int32()),
		UserID:      int(gofakeit.Int32()),
		ExpiryDate:  gofakeit.FutureDate().Truncate(time.Millisecond),
		ReceiveDate: time.Time{},
		Weight:      gofakeit.Float32Range(1, 100),
		Price:       gofakeit.Float32Range(1, 1000),
	}

	err := s.storage.CreateOrder(ctx, currentTime, orderInput)
	s.Require().NoError(err)

	returnInput := domain.OrderReturn{
		UserID:  orderInput.UserID,
		OrderID: orderInput.ID,
	}

	s.Run("создание возврата", func() {
		err := s.storage.CreateReturn(ctx, currentTime, returnInput.UserID, returnInput.OrderID)
		s.Require().NoError(err)
	})

	s.Run("получение возврата", func() {
		orderReturn, err := s.storage.GetReturn(ctx, returnInput.OrderID, currentTime)
		s.Require().NoError(err)
		s.Require().Equal(returnInput, orderReturn)
	})
}

func (s *PostgresStorageSuite) TestGetReturns() {
	ctx := context.Background()
	currentTime := time.Now().UTC()

	returnsInput := make([]domain.OrderReturn, 0, 3)
	ordersInput := make([]domain.Order, 0, 3)

	for i := 0; i < 5; i++ {
		order := domain.Order{
			ID:          int(gofakeit.Int32()),
			UserID:      int(gofakeit.Int32()),
			ExpiryDate:  gofakeit.FutureDate().Truncate(time.Millisecond),
			ReceiveDate: time.Time{},
			Weight:      gofakeit.Float32Range(1, 100),
			Price:       gofakeit.Float32Range(1, 1000),
		}

		orderReturn := domain.OrderReturn{
			UserID:  order.UserID,
			OrderID: order.ID,
		}

		err := s.storage.CreateOrder(ctx, currentTime, order)
		s.Require().NoError(err)

		err = s.storage.CreateReturn(ctx, currentTime, orderReturn.UserID, orderReturn.OrderID)
		s.Require().NoError(err)

		returnsInput = append(returnsInput, orderReturn)
		ordersInput = append(ordersInput, order)
	}

	s.Run("получение созданных возвратов", func() {
		orderReturns, err := s.storage.GetReturns(ctx, 1, 100)
		s.Require().NoError(err)
		s.Require().ElementsMatch(returnsInput, orderReturns)
	})
}
