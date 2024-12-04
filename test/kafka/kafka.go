package kafka

import (
	"fmt"
	"time"
	"workshop-1/config"
	"workshop-1/internal/app/kafka"
	"workshop-1/internal/app/kafka/event"

	"github.com/IBM/sarama"
	"github.com/stretchr/testify/suite"
)

type KafkaSuite struct {
	topic         string
	kafkaProducer *kafka.Kafka
	kafkaConsumer sarama.Consumer
	suite.Suite
}

func (s *KafkaSuite) SetupSuite() {
	addrs := []string{config.KafkaBrokers}
	s.topic = fmt.Sprintf("test-topic-%d", time.Now().UTC().Unix())

	kafkaProducer, err := kafka.New(
		addrs,
		event.NewDefaultFactory(),
		kafka.WithIdempotent(),
		kafka.WithRequiredAcks(sarama.WaitForAll),
		kafka.WithMaxRetries(3),
	)
	s.Require().NoError(err)

	s.kafkaProducer = kafkaProducer

	config := sarama.NewConfig()
	config.Version = sarama.MaxVersion
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	kafkaConsumer, err := sarama.NewConsumer(addrs, config)
	s.Require().NoError(err)

	s.kafkaConsumer = kafkaConsumer
}

func (s *KafkaSuite) TearDownSuite() {
	err := s.kafkaProducer.Close()
	s.Require().NoError(err)

	err = s.kafkaConsumer.Close()
	s.Require().NoError(err)
}
