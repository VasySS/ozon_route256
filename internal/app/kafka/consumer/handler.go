package consumer

import (
	"workshop-1/internal/app/logger"

	"github.com/IBM/sarama"
)

var _ sarama.ConsumerGroupHandler = (*handler)(nil)

type handler struct {
}

func NewHandler() *handler {
	return &handler{}
}

func (h *handler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *handler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h *handler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				return nil
			}

			session.MarkMessage(message, "")

			logger.Infow("получено сообщение",
				"тема", message.Topic,
				"партиция", message.Partition,
				"смещение", message.Offset,
				"ключ", string(message.Key),
				"значение", string(message.Value))
		case <-session.Context().Done():
			return nil
		}
	}
}
