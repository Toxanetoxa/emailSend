package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
)

type RedisQueue struct {
	client *redis.Client
	ctx    context.Context
	key    string
}

func NewRedisQueue(addr, key string) *RedisQueue {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	return &RedisQueue{
		client: rdb,
		ctx:    context.Background(),
		key:    key,
	}
}

func (q *RedisQueue) Enqueue(message string) error {
	return q.client.LPush(q.ctx, q.key, message).Err()
}

func (q *RedisQueue) Dequeue() (string, error) {
	result, err := q.client.RPop(q.ctx, q.key).Result()
	if err == redis.Nil {
		return "", nil // Queue is empty
	}
	return result, err
}
