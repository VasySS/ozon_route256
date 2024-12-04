package valkey

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"workshop-1/config"
	"workshop-1/internal/storage"

	"github.com/valkey-io/valkey-go"
)

var _ storage.Cache[struct{}, struct{}] = (*Cache[struct{}, struct{}])(nil)

type Cache[K comparable, V any] struct {
	valkey valkey.Client
	ttl    time.Duration
}

func New[K comparable, V any](valkeyURL string, ttl time.Duration) (*Cache[K, V], error) {
	valkeyClient, err := valkey.NewClient(valkey.ClientOption{
		InitAddress: []string{valkeyURL},
	})
	if err != nil {
		return nil, err
	}

	return &Cache[K, V]{
		valkey: valkeyClient,
		ttl:    ttl,
	}, nil
}

func (c *Cache[K, V]) Close() {
	c.valkey.Close()
}

func (c *Cache[K, V]) Set(ctx context.Context, key K, value V, now time.Time) error {
	valMarshalled, err := json.Marshal(value)
	if err != nil {
		return err
	}

	keyStr := fmt.Sprintf("%v", key)

	setCmd := c.valkey.B().Set().
		Key(keyStr).
		Value(string(valMarshalled)).
		Ex(config.CacheTTL).
		Build()

	return c.valkey.Do(ctx, setCmd).Error()
}

func (c *Cache[K, V]) Get(ctx context.Context, key K, now time.Time) (V, error) {
	keyStr := fmt.Sprintf("%v", key)

	getCmd := c.valkey.B().Get().Key(keyStr).Build()

	var val V
	if err := c.valkey.Do(ctx, getCmd).DecodeJSON(&val); err != nil {
		return val, err
	}

	return val, nil
}
