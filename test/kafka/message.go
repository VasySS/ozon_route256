package kafka

func (s *KafkaSuite) TestSendMessage() {
	msg := []byte("test-message")
	idempotentKey := s.kafkaProducer.EventFactory.IdempotentKey.Generate()

	s.Run("отправка сообщения в кафку", func() {
		s.kafkaProducer.SendMessage(msg, idempotentKey, s.topic)
	})

	s.Run("получение сообщения из кафки", func() {
		receivedMsg, err := s.consumeMessage()
		s.Require().NoError(err)

		s.Equal(string(msg), string(receivedMsg.Value))
		s.Equal(idempotentKey, string(receivedMsg.Key))
	})
}
