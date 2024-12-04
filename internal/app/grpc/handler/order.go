package handler

import (
	"context"
	"encoding/json"
	"time"

	"workshop-1/config"
	"workshop-1/internal/app/kafka/event"
	"workshop-1/internal/domain"
	"workshop-1/internal/dto"
	"workshop-1/internal/infrastructure/metrics"
	"workshop-1/internal/usecase"
	"workshop-1/pkg/pvz/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (h *PVZHandler) CreateOrder(ctx context.Context, req *pvz.CreateOrderRequest) (*emptypb.Empty, error) {
	handlerName := "CreateOrder"
	var errCode codes.Code

	defer func() {
		if errCode == codes.OK {
			metrics.IncOkResponseTotal(handlerName)
		} else {
			metrics.IncErrResponseTotal(handlerName, errCode.String())
		}
	}()

	if err := req.ValidateAll(); err != nil {
		errCode = codes.InvalidArgument
		return nil, status.Error(errCode, err.Error())
	}

	err := usecase.AcceptFromCourier(
		ctx,
		h.storage,
		time.Now(),
		dto.ToCreateOrder(req),
		dto.ToPackaging(req.PackagingType),
	)
	if err != nil {
		errCode = codes.Internal
		return nil, status.Error(errCode, err.Error())
	}

	event := h.kafka.EventFactory.CreateOrderEvent(
		domain.Order{
			ID:         int(req.GetId()),
			UserID:     int(req.GetUserId()),
			Weight:     req.GetWeight(),
			Price:      req.GetPrice(),
			ExpiryDate: req.GetExpiryDate().AsTime(),
		},
		event.OrderCreatedEvent,
	)

	b, err := json.Marshal(event)
	if err != nil {
		errCode = codes.Internal
		return nil, status.Error(errCode, err.Error())
	}

	if err := h.kafka.SendMessage(b, event.IdempotentKey, config.KafkaPVZTopic); err != nil {
		errCode = codes.Internal
		return nil, status.Error(errCode, err.Error())
	}

	return nil, nil
}

func (h *PVZHandler) GiveOrders(ctx context.Context, req *pvz.GiveOrdersRequest) (*emptypb.Empty, error) {
	handlerName := "GiveOrders"
	var errCode codes.Code

	defer func() {
		if errCode == codes.OK {
			metrics.IncOkResponseTotal(handlerName)
		} else {
			metrics.IncErrResponseTotal(handlerName, errCode.String())
		}
	}()

	if err := req.ValidateAll(); err != nil {
		errCode = codes.InvalidArgument
		return nil, status.Error(errCode, err.Error())
	}

	err := usecase.GiveToUser(
		ctx,
		h.storage,
		time.Now(),
		dto.Uint64ToIntArray(req.GetId()),
	)
	if err != nil {
		errCode = codes.Internal
		return nil, status.Error(errCode, err.Error())
	}

	event := h.kafka.EventFactory.CreateOrdersEvent(
		req.GetId(),
		event.OrderUpdatedEvent,
	)

	b, err := json.Marshal(event)
	if err != nil {
		errCode = codes.Internal
		return nil, status.Error(errCode, err.Error())
	}

	if err := h.kafka.SendMessage(b, event.IdempotentKey, config.KafkaPVZTopic); err != nil {
		errCode = codes.Internal
		return nil, status.Error(errCode, err.Error())
	}

	metrics.AddOrdersGiven(len(req.GetId()))

	return nil, nil
}

func (h *PVZHandler) GetOrders(ctx context.Context, req *pvz.GetOrdersRequest) (*pvz.GetOrdersResponse, error) {
	handlerName := "GetOrders"
	var errCode codes.Code

	defer func() {
		if errCode == codes.OK {
			metrics.IncOkResponseTotal(handlerName)
		} else {
			metrics.IncErrResponseTotal(handlerName, errCode.String())
		}
	}()

	if err := req.ValidateAll(); err != nil {
		errCode = codes.InvalidArgument
		return nil, status.Error(errCode, err.Error())
	}

	orders, err := usecase.UserOrders(
		ctx,
		h.storage,
		int(req.GetUserId()),
		int(req.GetLastN()),
		req.GetPvzOnly(),
	)
	if err != nil {
		errCode = codes.Internal
		return nil, status.Error(errCode, err.Error())
	}

	return &pvz.GetOrdersResponse{Orders: dto.FromDomainOrderArrayToGRPC(orders)}, nil
}
