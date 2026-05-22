package rmq

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
)

type Producer struct {
	client *redis.Client
	opts   *ProducerOptions
}

func NewProducer(client *redis.Client, opts ...ProducerOption) *Producer {
	p := Producer{
		client: client,
		opts:   &ProducerOptions{},
	}

	for _, opt := range opts {
		opt(p.opts)
	}

	repairProducer(p.opts)

	return &p
}

// 生产一条消息
func (p *Producer) SendMsg(ctx context.Context, topic, key string, val []byte) (string, error) {
	if topic == "" {
		return "", errors.New("redis XADD topic can't be empty")
	}
	return p.client.XAdd(ctx, &redis.XAddArgs{
		Stream: topic,
		MaxLen: p.opts.msgQueueLen,
		Values: map[string]interface{}{key: val},
	}).Result()
}
