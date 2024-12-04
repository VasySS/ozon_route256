package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewOrder(t *testing.T) {
	type args struct {
		orderID int
		userID  int
		weight  float32
		price   float32
		expiry  string
	}

	tests := []struct {
		name    string
		args    args
		want    Order
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "успешное создание заказа",
			args: args{
				orderID: 1,
				userID:  1,
				weight:  10,
				price:   100,
				expiry:  "01-01-2099",
			},
			want: Order{
				ID:          1,
				UserID:      1,
				Weight:      10,
				Price:       100,
				ExpiryDate:  time.Date(2099, 1, 1, 0, 0, 0, 0, time.UTC),
				ReceiveDate: time.Time{},
			},
			wantErr: assert.NoError,
		},
		{
			name: "попытка создания заказа с неправильным ID",
			args: args{
				orderID: -1,
				userID:  1,
				weight:  10,
				price:   100,
				expiry:  "01-01-2099",
			},
			want:    Order{},
			wantErr: assert.Error,
		},
		{
			name: "попытка создания заказа с неправильным UserID",
			args: args{
				orderID: 1,
				userID:  -1,
				weight:  10,
				price:   100,
				expiry:  "01-01-2099",
			},
			want:    Order{},
			wantErr: assert.Error,
		},
		{
			name: "попытка создания заказа с неправильным весом",
			args: args{
				orderID: 1,
				userID:  1,
				weight:  0,
				price:   100,
				expiry:  "01-01-2099",
			},
			want:    Order{},
			wantErr: assert.Error,
		},
		{
			name: "попытка создания заказа с неправильной ценой",
			args: args{
				orderID: 1,
				userID:  1,
				weight:  10,
				price:   -1,
				expiry:  "01-01-2099",
			},
			want:    Order{},
			wantErr: assert.Error,
		},
		{
			name: "попытка создания заказа с датой истечения в неправильном формате",
			args: args{
				orderID: 1,
				userID:  1,
				weight:  10,
				price:   100,
				expiry:  "2099-01-01",
			},
			want:    Order{},
			wantErr: assert.Error,
		},
		{
			name: "попытка создания заказа с некорректной датой истечения",
			args: args{
				orderID: 1,
				userID:  1,
				weight:  10,
				price:   100,
				expiry:  "01-13-2099",
			},
			want:    Order{},
			wantErr: assert.Error,
		},
		{
			name: "попытка создания заказа с датой истечения в прошлом",
			args: args{
				orderID: 1,
				userID:  1,
				weight:  10,
				price:   100,
				expiry:  "01-01-2000",
			},
			want:    Order{},
			wantErr: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewOrder(tt.args.orderID, tt.args.userID, tt.args.weight, tt.args.price, tt.args.expiry)

			tt.wantErr(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestOrder_CanReturnToCourier(t *testing.T) {
	tests := []struct {
		name        string
		currentTime time.Time
		o           *Order
		wantErr     assert.ErrorAssertionFunc
	}{
		{
			name:        "успешная проверка на возможность возвращения заказа курьеру",
			currentTime: time.Now(),
			o: &Order{
				ID:          1,
				UserID:      1,
				Weight:      10,
				Price:       100,
				ExpiryDate:  time.Now().Add(-(time.Hour)),
				ReceiveDate: time.Time{},
			},
			wantErr: assert.NoError,
		},
		{
			name:        "попытка вернуть уже полученный заказ курьеру",
			currentTime: time.Now(),
			o: &Order{
				ID:          1,
				UserID:      1,
				Weight:      10,
				Price:       100,
				ExpiryDate:  time.Date(2099, 1, 1, 0, 0, 0, 0, time.UTC),
				ReceiveDate: time.Now().Add(-(time.Hour * 1)),
			},
			wantErr: assert.Error,
		},
		{
			name:        "попытка вернуть не полученный заказ курьеру",
			currentTime: time.Now(),
			o: &Order{
				ID:          1,
				UserID:      1,
				Weight:      10,
				Price:       100,
				ExpiryDate:  time.Date(2099, 1, 1, 0, 0, 0, 0, time.UTC),
				ReceiveDate: time.Time{},
			},
			wantErr: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.o.CanReturnToCourier(tt.currentTime)
			tt.wantErr(t, err)
		})
	}
}

func TestOrder_CanAcceptUserReturn(t *testing.T) {
	tests := []struct {
		name        string
		currentTime time.Time
		o           *Order
		wantErr     assert.ErrorAssertionFunc
	}{
		{
			name:        "успешная проверка на возможность принять возврат заказа от пользователя",
			currentTime: time.Now(),
			o: &Order{
				ID:          1,
				UserID:      1,
				Weight:      10,
				Price:       100,
				ExpiryDate:  time.Date(2099, 1, 1, 0, 0, 0, 0, time.UTC),
				ReceiveDate: time.Now().Add(-(time.Hour * 24)),
			},
			wantErr: assert.NoError,
		},
		{
			name:        "попытка принять возврат ещё не полученного заказа",
			currentTime: time.Now(),
			o: &Order{
				ID:          1,
				UserID:      1,
				Weight:      10,
				Price:       100,
				ExpiryDate:  time.Date(2099, 1, 1, 0, 0, 0, 0, time.UTC),
				ReceiveDate: time.Time{},
			},
			wantErr: assert.Error,
		},
		{
			name:        "попытка принять возврат заказа, полученного более 48 часов назад",
			currentTime: time.Now(),
			o: &Order{
				ID:          1,
				UserID:      1,
				Weight:      10,
				Price:       100,
				ExpiryDate:  time.Date(2099, 1, 1, 0, 0, 0, 0, time.UTC),
				ReceiveDate: time.Now().Add(-(time.Hour * 49)),
			},
			wantErr: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.o.CanAcceptUserReturn(tt.currentTime)
			tt.wantErr(t, err)
		})
	}
}
