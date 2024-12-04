package kafka

import (
	"testing"

	"github.com/IBM/sarama"
	"github.com/IBM/sarama/mocks"
	"github.com/stretchr/testify/assert"
)

func TestKafka_SendMessage(t *testing.T) {
	t.Parallel()

	type args struct {
		msg           []byte
		idempotentKey string
		topic         string
	}

	tests := []struct {
		name    string
		args    args
		setup   func(*mocks.SyncProducer)
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "успешное отправление сообщения в кафку",
			args: args{
				msg:           []byte("test-message"),
				idempotentKey: "test-idempotent-key",
				topic:         "test-topic",
			},
			setup: func(k *mocks.SyncProducer) {
				k.ExpectSendMessageAndSucceed()
			},
			wantErr: assert.NoError,
		},
		{
			name: "ошибка при отправке сообщения в кафку",
			args: args{
				msg:           []byte("test-message"),
				idempotentKey: "test-idempotent-key",
				topic:         "test-topic",
			},
			setup: func(k *mocks.SyncProducer) {
				k.ExpectSendMessageAndFail(sarama.ErrBrokerNotAvailable)
			},
			wantErr: assert.Error,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			kafkaMock := mocks.NewSyncProducer(t, nil)
			tt.setup(kafkaMock)

			kafka := Kafka{
				SyncProducer: kafkaMock,
			}
			defer kafka.Close()

			err := kafka.SendMessage(tt.args.msg, tt.args.idempotentKey, tt.args.topic)
			tt.wantErr(t, err)
		})
	}
}
