package handler

import (
	"context"
	"encoding/json"
	"time"

	"workshop-1/config"
	"workshop-1/internal/app/kafka/event"
	"workshop-1/internal/dto"
	"workshop-1/internal/usecase"
	"workshop-1/pkg/pvz/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (h *PVZHandler) CreateOrderReturn(
	ctx context.Context,
	req *pvz.CreateOrderReturnRequest,
) (*emptypb.Empty, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err := usecase.AcceptUserReturn(
		ctx,
		h.storage,
		time.Now(),
		int(req.UserId),
		int(req.OrderId),
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	event := h.kafka.EventFactory.CreateOrderReturnEvent(
		int(req.GetOrderId()),
		event.OrderReturnCreatedEvent,
	)

	b, err := json.Marshal(event)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if err := h.kafka.SendMessage(b, event.IdempotentKey, config.KafkaPVZTopic); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return nil, nil
}

func (h *PVZHandler) GiveOrderToCourier(
	ctx context.Context,
	req *pvz.GiveOrderToCourierRequest,
) (*emptypb.Empty, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err := usecase.ReturnToCourier(
		ctx,
		h.storage,
		time.Now(),
		int(req.OrderId),
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	event := h.kafka.EventFactory.CreateOrderReturnEvent(
		int(req.GetOrderId()),
		event.OrderReturnUpdatedEvent,
	)

	b, err := json.Marshal(event)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if err := h.kafka.SendMessage(b, event.IdempotentKey, config.KafkaPVZTopic); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return nil, nil
}

func (h *PVZHandler) GetOrderReturns(
	ctx context.Context,
	req *pvz.GetOrderReturnsRequest,
) (*pvz.GetOrderReturnsResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	returns, err := usecase.UserReturns(
		ctx,
		h.storage,
		int(req.GetPage()),
		int(req.GetPageSize()),
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pvz.GetOrderReturnsResponse{OrderReturns: dto.FromDomainReturnArrayToGRPC(returns)}, nil
}
