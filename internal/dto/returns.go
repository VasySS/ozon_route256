package dto

import (
	"workshop-1/internal/domain"
	"workshop-1/pkg/pvz/v1"
)

func FromDomainReturnToGRPC(ret domain.OrderReturn) *pvz.OrderReturn {
	return &pvz.OrderReturn{
		UserId:  int64(ret.UserID),
		OrderId: int64(ret.OrderID),
	}
}

func FromDomainReturnArrayToGRPC(returns []domain.OrderReturn) []*pvz.OrderReturn {
	res := make([]*pvz.OrderReturn, 0, len(returns))
	for _, ret := range returns {
		res = append(res, FromDomainReturnToGRPC(ret))
	}
	return res
}
