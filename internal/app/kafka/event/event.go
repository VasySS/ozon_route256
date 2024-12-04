package event

import (
	"time"
	"workshop-1/internal/domain"
)

type OrderEventType string
type OrderReturnEventType string

const (
	OrderCreatedEvent       OrderEventType       = "order-created"
	OrderUpdatedEvent       OrderEventType       = "order-updated"
	OrderReturnCreatedEvent OrderReturnEventType = "order-return-created"
	OrderReturnUpdatedEvent OrderReturnEventType = "order-return-updated"
)

type event struct {
	IdempotentKey string    `json:"idempotentKey"`
	Timestamp     time.Time `json:"timestamp"`
}

type OrderEvent struct {
	Order     domain.Order
	EventType OrderEventType `json:"eventType"`
	event
}

type OrdersEvent struct {
	IDs       []uint64
	EventType OrderEventType `json:"eventType"`
	event
}

type OrderReturnEvent struct {
	OrderID   int
	EventType OrderReturnEventType `json:"eventType"`
	event
}
