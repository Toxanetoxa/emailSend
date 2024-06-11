package redis

import (
	"context"
	"email-sendler/internal/email"
	"email-sendler/internal/queueTypes"
	"errors"
	"github.com/go-redis/redis/v8"
)

type Queue struct {
	client *redis.Client
	ctx    context.Context
	key    string
}

func NewRedisQueue(addr, password string, db int, key string) queueTypes.Queue {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &Queue{
		client: rdb,
		ctx:    context.Background(),
		key:    key,
	}
}

func (q *Queue) Enqueue(message email.Message) error {
	return q.client.LPush(q.ctx, q.key, message.Body).Err()
}

func (q *Queue) Dequeue() (email.Message, error) {
	result, err := q.client.RPop(q.ctx, q.key).Result()
	if errors.Is(err, redis.Nil) {
		return email.Message{}, nil // Очередь пуста
	}
	return email.Message{Body: result}, err
}
