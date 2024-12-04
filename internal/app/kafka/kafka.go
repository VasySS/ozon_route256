package kafka

import (
	"workshop-1/internal/app/kafka/event"

	"github.com/IBM/sarama"
)

type Kafka struct {
	SyncProducer sarama.SyncProducer
	EventFactory *event.Factory
}

func New(addrs []string, factory *event.Factory, opts ...Option) (*Kafka, error) {
	config, err := newConfig(opts...)
	if err != nil {
		return nil, err
	}

	syncProducer, err := sarama.NewSyncProducer(addrs, config)
	if err != nil {
		return nil, err
	}

	return &Kafka{
		SyncProducer: syncProducer,
		EventFactory: factory,
	}, nil
}

func (k *Kafka) SendMessage(msg []byte, idempotentKey, topic string) error {
	pMsg := &sarama.ProducerMessage{
		Key:   sarama.StringEncoder(idempotentKey),
		Topic: topic,
		Value: sarama.ByteEncoder(msg),
	}

	_, _, err := k.SyncProducer.SendMessage(pMsg)
	if err != nil {
		return err
	}

	return nil
}

func (k *Kafka) Close() error {
	return k.SyncProducer.Close()
}
