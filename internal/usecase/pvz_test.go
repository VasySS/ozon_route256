package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"workshop-1/internal/domain"
	"workshop-1/internal/domain/strategy"
	"workshop-1/internal/dto"
	"workshop-1/internal/usecase/mock"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	errNoOrderFound       = errors.New("заказ не существует")
	errNoOrderReturnFound = errors.New("возврат не существует")
)

func TestAcceptFromCourier(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	type args struct {
		storage     *mock.StorageMock
		currentTime time.Time
		orderDTO    dto.CreateOrder
		packaging   strategy.Packaging
		setup       func(args)
	}

	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "успешное создание заказа",
			args: args{
				orderDTO:    dto.CreateOrder{ID: 1, UserID: 1, Weight: 10, Price: 100, ExpiryDate: "01-01-2099"},
				packaging:   strategy.Wrap{},
				currentTime: time.Now().UTC(),
				setup: func(a args) {
					createOrderInput, err := domain.NewOrder(
						a.orderDTO.ID,
						a.orderDTO.UserID,
						a.orderDTO.Weight,
						a.orderDTO.Price,
						a.orderDTO.ExpiryDate,
					)
					require.NoError(t, err)

					err = a.packaging.Apply(&createOrderInput)
					require.NoError(t, err)

					a.storage.GetOrderMock.Expect(ctx, a.orderDTO.ID, a.currentTime).Return(domain.Order{}, errNoOrderFound)
					a.storage.CreateOrderMock.Expect(ctx, a.currentTime, createOrderInput).Return(nil)
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "создание заказа с неправильным форматом даты истечения срока",
			args: args{
				orderDTO:    dto.CreateOrder{ID: 1, UserID: 1, Weight: 10, Price: 100, ExpiryDate: "2099-01-01"},
				packaging:   strategy.Wrap{},
				currentTime: time.Now().UTC(),
				setup:       func(a args) {},
			},
			wantErr: assert.Error,
		},
		{
			name: "создание заказа с неверным весом заказа для упаковки",
			args: args{
				orderDTO:    dto.CreateOrder{ID: 1, UserID: 1, Weight: 11, Price: 100, ExpiryDate: "01-01-2099"},
				packaging:   strategy.Bag{},
				currentTime: time.Now().UTC(),
				setup:       func(a args) {},
			},
			wantErr: assert.Error,
		},
		{
			name: "создание уже существующего заказа",
			args: args{
				orderDTO:    dto.CreateOrder{ID: 1, UserID: 1, Weight: 10, Price: 100, ExpiryDate: "01-01-2099"},
				packaging:   strategy.Wrap{},
				currentTime: time.Now().UTC(),
				setup: func(a args) {
					a.storage.GetOrderMock.Expect(ctx, a.orderDTO.ID, a.currentTime).Return(domain.Order{}, nil)
				},
			},
			wantErr: assert.Error,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := minimock.NewController(t)
			tt.args.storage = mock.NewStorageMock(ctrl)
			tt.args.setup(tt.args)

			err := AcceptFromCourier(ctx, tt.args.storage, tt.args.currentTime, tt.args.orderDTO, tt.args.packaging)
			tt.wantErr(t, err)
		})
	}
}

func TestReturnToCourier(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	type args struct {
		storage     *mock.StorageMock
		currentTime time.Time
		orderID     int
		setup       func(args)
	}

	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "успешное возвращение заказа курьеру, который был ранее возвращен пользователем",
			args: args{
				orderID:     1,
				currentTime: time.Now().UTC(),
				setup: func(a args) {
					a.storage.GetReturnMock.Expect(ctx, a.orderID, a.currentTime).Return(domain.OrderReturn{}, nil)
					a.storage.DeleteOrderMock.Expect(ctx, a.orderID).Return(nil)
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "успешное возвращение заказа курьеру, который не был получен пользователем",
			args: args{
				orderID:     1,
				currentTime: time.Now().UTC(),
				setup: func(a args) {
					order := domain.Order{
						ID:          a.orderID,
						ExpiryDate:  time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
						ReceiveDate: time.Time{},
					}

					a.storage.GetReturnMock.Expect(ctx, a.orderID, a.currentTime).Return(domain.OrderReturn{}, errNoOrderFound)
					a.storage.GetOrderMock.Expect(ctx, a.orderID, a.currentTime).Return(order, nil)
					a.storage.DeleteOrderMock.Expect(ctx, a.orderID).Return(nil)
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "ошибка возвращения заказа курьеру, который ещё не был получен пользователем",
			args: args{
				orderID:     1,
				currentTime: time.Now().UTC(),
				setup: func(a args) {
					order := domain.Order{
						ID:          a.orderID,
						ExpiryDate:  time.Date(2099, 1, 1, 0, 0, 0, 0, time.UTC),
						ReceiveDate: time.Time{},
					}

					a.storage.GetReturnMock.Expect(ctx, a.orderID, a.currentTime).
						Return(domain.OrderReturn{}, errNoOrderFound)
					a.storage.GetOrderMock.Expect(ctx, a.orderID, a.currentTime).Return(order, nil)
				},
			},
			wantErr: assert.Error,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := minimock.NewController(t)
			tt.args.storage = mock.NewStorageMock(ctrl)
			tt.args.setup(tt.args)

			err := ReturnToCourier(ctx, tt.args.storage, tt.args.currentTime, tt.args.orderID)
			tt.wantErr(t, err)
		})
	}
}

func TestGiveToUser(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	type args struct {
		storage     *mock.StorageMock
		orderIDs    []int
		currentTime time.Time
		setup       func(args)
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "успешная выдача всех заказов пользователю",
			args: args{
				orderIDs:    []int{1, 2, 3},
				currentTime: time.Now().UTC(),
				setup: func(a args) {
					orders := []domain.Order{
						{ID: 1, UserID: 123, ExpiryDate: time.Date(2099, 1, 1, 0, 0, 0, 0, time.UTC)},
						{ID: 2, UserID: 123, ExpiryDate: time.Date(2099, 1, 1, 0, 0, 0, 0, time.UTC)},
						{ID: 3, UserID: 123, ExpiryDate: time.Date(2099, 1, 1, 0, 0, 0, 0, time.UTC)},
					}

					a.storage.GetOrdersMock.Expect(ctx, a.orderIDs).Return(orders, nil)
					for _, order := range orders {
						a.storage.GetReturnMock.When(ctx, order.ID, a.currentTime).
							Then(domain.OrderReturn{}, errNoOrderFound)
					}
					a.storage.UpdateReceiveDatesMock.Expect(ctx, a.currentTime, a.orderIDs).Return(nil)
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "успешная выдача части заказов пользователю",
			args: args{
				orderIDs:    []int{1, 2, 3},
				currentTime: time.Now().UTC(),
				setup: func(a args) {
					// 1 - успешно выдан, 2 - уже был ранее возвращен в пвз, 3 - вышел срок хранения
					orders := []domain.Order{
						{ID: 1, UserID: 123, ExpiryDate: time.Date(2099, 1, 1, 0, 0, 0, 0, time.UTC)},
						{ID: 2, UserID: 123, ExpiryDate: time.Date(2099, 1, 1, 0, 0, 0, 0, time.UTC)},
						{ID: 3, UserID: 123, ExpiryDate: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)},
					}

					a.storage.GetOrdersMock.Expect(ctx, a.orderIDs).Return(orders, nil)
					a.storage.GetReturnMock.When(ctx, 1, a.currentTime).Then(domain.OrderReturn{}, errNoOrderFound)
					a.storage.GetReturnMock.When(ctx, 2, a.currentTime).Then(domain.OrderReturn{}, nil)
					a.storage.GetReturnMock.When(ctx, 3, a.currentTime).Then(domain.OrderReturn{}, errNoOrderFound)
					a.storage.UpdateReceiveDatesMock.Expect(ctx, a.currentTime, []int{1}).Return(nil)
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "попытка получить заказы для разных пользователей",
			args: args{
				orderIDs:    []int{1, 2, 3},
				currentTime: time.Now(),
				setup: func(a args) {
					orders := []domain.Order{
						{ID: 1, UserID: 123, ExpiryDate: time.Date(2099, 1, 1, 0, 0, 0, 0, time.UTC)},
						{ID: 2, UserID: 321, ExpiryDate: time.Date(2099, 1, 1, 0, 0, 0, 0, time.UTC)},
						{ID: 3, UserID: 123, ExpiryDate: time.Date(2099, 1, 1, 0, 0, 0, 0, time.UTC)},
					}

					a.storage.GetOrdersMock.Expect(ctx, a.orderIDs).Return(orders, nil)
				},
			},
			wantErr: assert.Error,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := minimock.NewController(t)
			tt.args.storage = mock.NewStorageMock(ctrl)
			tt.args.setup(tt.args)

			err := GiveToUser(ctx, tt.args.storage, tt.args.currentTime, tt.args.orderIDs)
			tt.wantErr(t, err)
		})
	}
}

func TestUserOrders(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	type args struct {
		storage   *mock.StorageMock
		userID    int
		lastN     int
		inPVZOnly bool
		setup     func(args)
	}

	tests := []struct {
		name       string
		args       args
		wantOrders []domain.Order
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name: "успешное получение списка всех заказов пользователя",
			args: args{
				userID:    123,
				lastN:     0,
				inPVZOnly: false,
				setup: func(a args) {
					orders := []domain.Order{
						{ID: 1, UserID: 123},
						{ID: 2, UserID: 123},
						{ID: 3, UserID: 123},
					}

					a.storage.GetOrdersByUserIDMock.Expect(ctx, a.userID, a.lastN).Return(orders, nil)
				},
			},
			wantErr: assert.NoError,
			wantOrders: []domain.Order{
				{ID: 1, UserID: 123},
				{ID: 2, UserID: 123},
				{ID: 3, UserID: 123},
			},
		},
		{
			name: "успешное получение списка ещё не полученных заказов пользователя",
			args: args{
				userID:    123,
				lastN:     0,
				inPVZOnly: true,
				setup: func(a args) {
					orders := []domain.Order{
						{ID: 1, UserID: 123, ReceiveDate: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)},
						{ID: 2, UserID: 123, ReceiveDate: time.Time{}},
						{ID: 3, UserID: 123, ReceiveDate: time.Time{}},
					}

					a.storage.GetOrdersByUserIDMock.Expect(ctx, a.userID, a.lastN).Return(orders, nil)
				},
			},
			wantErr: assert.NoError,
			wantOrders: []domain.Order{
				{ID: 2, UserID: 123, ReceiveDate: time.Time{}},
				{ID: 3, UserID: 123, ReceiveDate: time.Time{}},
			},
		},
		{
			name: "ошибка при получении списка заказов пользователя",
			args: args{
				userID:    123,
				lastN:     0,
				inPVZOnly: false,
				setup: func(a args) {
					a.storage.GetOrdersByUserIDMock.Expect(ctx, a.userID, a.lastN).Return(nil, errNoOrderFound)
				},
			},
			wantErr:    assert.Error,
			wantOrders: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := minimock.NewController(t)
			tt.args.storage = mock.NewStorageMock(ctrl)
			tt.args.setup(tt.args)

			orders, err := UserOrders(ctx, tt.args.storage, tt.args.userID, tt.args.lastN, tt.args.inPVZOnly)
			tt.wantErr(t, err)
			assert.Equal(t, tt.wantOrders, orders)
		})
	}
}

func TestAcceptUserReturn(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	type args struct {
		storage     *mock.StorageMock
		userID      int
		orderID     int
		currentTime time.Time
		setup       func(args)
	}

	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "успешное принятие возврата от пользователя",
			args: args{
				userID:      123,
				orderID:     1,
				currentTime: time.Now().UTC(),
				setup: func(a args) {
					order := domain.Order{
						ID:          a.orderID,
						UserID:      a.userID,
						ReceiveDate: time.Now().Add(-(time.Hour * 47)),
						ExpiryDate:  time.Now().Add(-(time.Hour * 23)),
					}

					a.storage.GetReturnMock.Expect(ctx, a.orderID, a.currentTime).
						Return(domain.OrderReturn{}, errNoOrderReturnFound)
					a.storage.GetOrderMock.Expect(ctx, a.orderID, a.currentTime).Return(order, nil)
					a.storage.CreateReturnMock.Expect(ctx, a.currentTime, a.userID, a.orderID).Return(nil)
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "попытка вернуть заказ, полученный более 48 часов назад",
			args: args{
				userID:      123,
				orderID:     1,
				currentTime: time.Now().UTC(),
				setup: func(a args) {
					order := domain.Order{
						ID:          a.orderID,
						UserID:      a.userID,
						ReceiveDate: time.Now().Add(-(time.Hour * 49)),
						ExpiryDate:  time.Now().Add(-(time.Hour * 23)),
					}

					a.storage.GetReturnMock.Expect(ctx, a.orderID, a.currentTime).
						Return(domain.OrderReturn{}, errNoOrderReturnFound)
					a.storage.GetOrderMock.Expect(ctx, a.orderID, a.currentTime).Return(order, nil)
				},
			},
			wantErr: assert.Error,
		},
		{
			name: "попытка вернуть заказ, который уже был возвращён",
			args: args{
				userID:      123,
				orderID:     1,
				currentTime: time.Now().UTC(),
				setup: func(a args) {
					a.storage.GetReturnMock.Expect(ctx, a.orderID, a.currentTime).Return(domain.OrderReturn{}, nil)
				},
			},
			wantErr: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := minimock.NewController(t)
			tt.args.storage = mock.NewStorageMock(ctrl)
			tt.args.setup(tt.args)

			err := AcceptUserReturn(ctx, tt.args.storage, tt.args.currentTime, tt.args.userID, tt.args.orderID)
			tt.wantErr(t, err)
		})
	}
}

func TestUserReturns(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	type args struct {
		storage  *mock.StorageMock
		page     int
		pageSize int
		setup    func(args)
	}

	tests := []struct {
		name        string
		args        args
		wantErr     assert.ErrorAssertionFunc
		wantReturns []domain.OrderReturn
	}{
		{
			name: "успешное получение списка возвратов",
			args: args{
				page:     1,
				pageSize: 100,
				setup: func(a args) {
					returns := []domain.OrderReturn{
						{OrderID: 1, UserID: 555},
						{OrderID: 2, UserID: 444},
						{OrderID: 3, UserID: 777},
					}

					a.storage.GetReturnsMock.Expect(ctx, a.page, a.pageSize).Return(returns, nil)
				},
			},
			wantReturns: []domain.OrderReturn{
				{OrderID: 1, UserID: 555},
				{OrderID: 2, UserID: 444},
				{OrderID: 3, UserID: 777},
			},
			wantErr: assert.NoError,
		},
		{
			name: "ошибка при получении списка возвратов",
			args: args{
				page:     1,
				pageSize: 100,
				setup: func(a args) {
					a.storage.GetReturnsMock.Expect(ctx, a.page, a.pageSize).Return(nil, errNoOrderReturnFound)
				},
			},
			wantErr:     assert.Error,
			wantReturns: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := minimock.NewController(t)
			tt.args.storage = mock.NewStorageMock(ctrl)
			tt.args.setup(tt.args)

			returns, err := UserReturns(ctx, tt.args.storage, tt.args.page, tt.args.pageSize)
			tt.wantErr(t, err)
			assert.Equal(t, tt.wantReturns, returns)
		})
	}
}
