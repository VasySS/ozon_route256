package consumer

import (
	"time"

	"github.com/IBM/sarama"
)

func newConfig(opts ...Option) (*sarama.Config, error) {
	conf := sarama.NewConfig()

	conf.Version = sarama.MaxVersion
	conf.Consumer.Offsets.AutoCommit.Enable = true
	conf.Consumer.Offsets.AutoCommit.Interval = 1 * time.Second
	conf.Consumer.Group.Heartbeat.Interval = 3 * time.Second
	conf.Consumer.Offsets.Initial = sarama.OffsetOldest
	conf.Consumer.Return.Errors = true

	for _, opt := range opts {
		if err := opt.Apply(conf); err != nil {
			return nil, err
		}
	}

	return conf, nil
}
