package dto

import (
	"workshop-1/internal/domain"
	"workshop-1/pkg/pvz/v1"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type CreateOrder struct {
	ID         int
	UserID     int
	ExpiryDate string
	Weight     float32
	Price      float32
}

func ToCreateOrder(req *pvz.CreateOrderRequest) CreateOrder {
	if req == nil {
		return CreateOrder{}
	}

	co := CreateOrder{
		ID:         int(req.GetId()),
		UserID:     int(req.GetUserId()),
		ExpiryDate: req.ExpiryDate.AsTime().Format("02-01-2006"),
		Weight:     req.Weight,
		Price:      req.Price,
	}

	return co
}

func FromDomainOrderToGRPC(order domain.Order) *pvz.Order {
	return &pvz.Order{
		Id:          int64(order.ID),
		UserId:      int64(order.UserID),
		Weight:      order.Weight,
		Price:       order.Price,
		ReceiveDate: timestamppb.New(order.ReceiveDate),
		ExpiryDate:  timestamppb.New(order.ExpiryDate),
	}
}

func FromDomainOrderArrayToGRPC(orders []domain.Order) []*pvz.Order {
	res := make([]*pvz.Order, 0, len(orders))
	for _, order := range orders {
		res = append(res, FromDomainOrderToGRPC(order))
	}
	return res
}

func Uint64ToIntArray(ids []uint64) []int {
	res := make([]int, 0, len(ids))
	for _, id := range ids {
		res = append(res, int(id))
	}
	return res
}
