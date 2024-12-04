package event

import (
	"time"
	"workshop-1/internal/domain"

	"github.com/google/uuid"
)

type idempotentKeyGenerator interface {
	Generate() string
}

type operationMomentGenerator interface {
	Generate() time.Time
}

type Factory struct {
	IdempotentKey   idempotentKeyGenerator
	OperationMoment operationMomentGenerator
}

func NewDefaultFactory() *Factory {
	return &Factory{
		NewUUIDv4Generator(),
		&Clock{},
	}
}

func (f *Factory) CreateOrderEvent(order domain.Order, eventType OrderEventType) OrderEvent {
	return OrderEvent{
		Order:     order,
		EventType: eventType,
		event: event{
			IdempotentKey: f.IdempotentKey.Generate(),
			Timestamp:     f.OperationMoment.Generate(),
		},
	}
}

func (f *Factory) CreateOrdersEvent(ids []uint64, eventType OrderEventType) OrdersEvent {
	return OrdersEvent{
		IDs:       ids,
		EventType: eventType,
		event: event{
			IdempotentKey: f.IdempotentKey.Generate(),
			Timestamp:     f.OperationMoment.Generate(),
		},
	}
}

func (f *Factory) CreateOrderReturnEvent(orderID int, eventType OrderReturnEventType) OrderReturnEvent {
	return OrderReturnEvent{
		OrderID:   orderID,
		EventType: eventType,
		event: event{
			IdempotentKey: f.IdempotentKey.Generate(),
			Timestamp:     f.OperationMoment.Generate(),
		},
	}
}

func New(
	idemKeyGen idempotentKeyGenerator,
	momentGen operationMomentGenerator,
) *Factory {
	return &Factory{
		idemKeyGen,
		momentGen,
	}
}

type Clock struct{}

func (c *Clock) Generate() time.Time {
	return time.Now().UTC()
}

type UUIDv4Generator struct {
	gen uuid.UUID
}

func NewUUIDv4Generator() *UUIDv4Generator {
	gen := uuid.New()

	return &UUIDv4Generator{
		gen: gen,
	}
}

func (g *UUIDv4Generator) Generate() string {
	return g.gen.String()
}
