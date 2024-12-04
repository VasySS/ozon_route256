package handler

import (
	"workshop-1/internal/app/kafka"
	"workshop-1/internal/storage/postgres"
	"workshop-1/pkg/pvz/v1"
)

type PVZHandler struct {
	storage postgres.Facade
	kafka   *kafka.Kafka

	pvz.UnimplementedPVZServiceServer
}

func NewHandler(storage postgres.Facade, kafka *kafka.Kafka) *PVZHandler {
	return &PVZHandler{
		storage: storage,
		kafka:   kafka,
	}
}
