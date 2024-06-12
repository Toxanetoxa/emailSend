package redis

import (
	"context"
	"email-sendler/internal/email"
	"email-sendler/internal/queueTypes"
	"github.com/go-redis/redis/v8"
)

type Queue struct {
	Client *redis.Client
	Ctx    context.Context
	Key    string
}

func NewRedisQueue(addr string, password string, db int, key string) queueTypes.Queue {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &Queue{
		Client: rdb,
		Ctx:    context.Background(),
		Key:    key,
	}
}

func (q *Queue) Enqueue(message email.Message) error {
	data, err := message.Serialize()
	if err != nil {
		return err
	}

	return q.Client.RPush(q.Ctx, q.Key, data).Err()
}

func (q *Queue) Dequeue() (email.Message, error) {
	data, err := q.Client.LPop(q.Ctx, q.Key).Result()
	if err != nil {
		return email.Message{}, err
	}

	msg, err := email.Deserialize([]byte(data))
	if err != nil {
		return email.Message{}, err
	}

	return *msg, nil
}
