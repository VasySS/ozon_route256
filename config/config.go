package config

import (
	"time"
)

const (
	PostgresURL     = "postgres://postgres:postgres@localhost:5432/route256?sslmode=disable"
	ValkeyURL       = "localhost:6379"
	GRPCHost        = "localhost:7000"
	HTTPHost        = "localhost:7001"
	SwaggerHost     = "localhost:7002"
	SwaggerPath     = "./pkg/pvz/v1/pvz_service.swagger.json"
	ShutdownTimeout = time.Second * 3
)

const (
	CacheTTL      = time.Minute * 3
	CacheCapacity = 100
)

const (
	JWTSecret      = "jwtsecret123"
	AccessTokenTTL = 24 * time.Hour
)

var (
	KafkaBrokers         = "localhost:9092"
	KafkaPVZTopic        = "pvz.events-log"
	KafkaConsumerGroupID = "pvz-consumer"
)
