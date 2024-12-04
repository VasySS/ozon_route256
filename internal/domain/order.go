package domain

import (
	"fmt"
	"time"
)

type Order struct {
	ID          int
	UserID      int
	ExpiryDate  time.Time
	ReceiveDate time.Time
	Weight      float32
	Price       float32
}

//nolint:gocyclo
func NewOrder(orderID, userID int, weight, price float32, expiry string) (Order, error) {
	o := Order{}

	if err := o.SetOrderID(orderID); err != nil {
		return Order{}, err
	}

	if err := o.SetUserID(userID); err != nil {
		return Order{}, err
	}

	if err := o.SetWeight(weight); err != nil {
		return Order{}, err
	}

	if err := o.SetPrice(price); err != nil {
		return Order{}, err
	}

	if err := o.SetExpiryDate(expiry); err != nil {
		return Order{}, err
	}

	return o, nil
}

func (o *Order) SetOrderID(orderID int) error {
	if orderID < 0 {
		return OrderError{"неверный ID заказа: " + fmt.Sprint(orderID)}
	}

	o.ID = orderID
	return nil
}

func (o *Order) SetUserID(userID int) error {
	if userID < 0 {
		return OrderError{"неверный ID пользователя: " + fmt.Sprint(userID)}
	}

	o.UserID = userID
	return nil
}

func (o *Order) SetExpiryDate(expiry string) error {
	expiryDate, err := time.Parse("02-01-2006", expiry)
	if err != nil {
		return fmt.Errorf("неверный формат даты: %w", err)
	}
	expiryDate = expiryDate.UTC()

	if expiryDate.Before(time.Now().UTC()) {
		return OrderError{"дата окончания хранения заказа в прошлом: " + expiry}
	}

	o.ExpiryDate = expiryDate
	return nil
}

func (o *Order) SetWeight(weight float32) error {
	if weight <= 0 {
		return OrderError{"неверный вес: " + fmt.Sprint(weight)}
	}

	o.Weight = weight
	return nil
}

func (o *Order) SetPrice(price float32) error {
	if price < 0 {
		return OrderError{"неверная цена: " + fmt.Sprint(price)}
	}

	o.Price = price
	return nil
}

func (o *Order) CanReturnToCourier(currentTime time.Time) error {
	if !o.ReceiveDate.IsZero() {
		return OrderError{"заказ уже был получен"}
	}

	if o.ExpiryDate.After(currentTime) {
		return OrderError{"время хранения заказа ещё не истекло"}
	}

	return nil
}

func (o *Order) CanAcceptUserReturn(currentTime time.Time) error {
	if o.ReceiveDate.IsZero() {
		return OrderError{"заказ ещё не был получен"}
	}

	if o.ReceiveDate.Add(time.Hour * 48).Before(currentTime) {
		return OrderError{"заказ был получен более 48 часов назад"}
	}

	return nil
}
