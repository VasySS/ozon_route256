package consumer

import (
	"context"
	"sync"
	"workshop-1/internal/app/logger"

	"github.com/IBM/sarama"
)

type consumerGroup struct {
	sarama.ConsumerGroup
	handler sarama.ConsumerGroupHandler
	topics  []string
}

func NewConsumerGroup(
	brokers []string,
	groupID string,
	topics []string,
	handler sarama.ConsumerGroupHandler,
	opts ...Option,
) (*consumerGroup, error) {
	config, err := newConfig(opts...)
	if err != nil {
		return nil, err
	}

	cg, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, err
	}

	return &consumerGroup{
		ConsumerGroup: cg,
		handler:       handler,
		topics:        topics,
	}, nil
}

func (c *consumerGroup) Run(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)

	go func() {
		defer wg.Done()
		logger.Info("kafka consumer запущен")

		for {
			if err := c.Consume(ctx, c.topics, c.handler); err != nil {
				logger.Error(err)
			}

			if ctx.Err() != nil {
				logger.Info("kafka consumer канал закрыт")
				return
			}
		}
	}()
}
