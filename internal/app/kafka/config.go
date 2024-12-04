package kafka

import (
	"github.com/IBM/sarama"
)

func newConfig(opts ...Option) (*sarama.Config, error) {
	conf := sarama.NewConfig()

	conf.Version = sarama.MaxVersion
	conf.Producer.Return.Successes = true
	conf.Producer.Return.Errors = true
	conf.Net.MaxOpenRequests = 1

	for _, opt := range opts {
		if err := opt.Apply(conf); err != nil {
			return nil, err
		}
	}

	return conf, nil
}
