package kafka

import (
	"errors"
	"time"

	"github.com/IBM/sarama"
	"golang.org/x/net/context"
)

func (s *KafkaSuite) consumeMessage() (*sarama.ConsumerMessage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	partitions, err := s.kafkaConsumer.Partitions(s.topic)
	s.Require().NoError(err)

	for _, partition := range partitions {
		partitionConsumer, err := s.kafkaConsumer.ConsumePartition(s.topic, partition, sarama.OffsetOldest)
		s.Require().NoError(err)

		for {
			select {
			case <-ctx.Done():
				return nil, errors.New("превышено время ожидания сообщения из кафки")
			case msg, ok := <-partitionConsumer.Messages():
				if !ok {
					return nil, errors.New("канал сообщений закрыт")
				}

				return msg, nil
			}
		}
	}

	return nil, errors.New("не удалось получить сообщения из кафки")
}
