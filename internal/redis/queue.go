package redis

import (
	"context"
	"email-sendler/internal/email"
	"email-sendler/internal/queueTypes"
	"github.com/go-redis/redis/v8"
)

type Queue struct {
	client *redis.Client
	ctx    context.Context
	key    string
}

func NewRedisQueue(addr string, password string, db int, key string) queueTypes.Queue {
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
	data, err := message.Serialize()
	if err != nil {
		return err
	}

	return q.client.RPush(q.ctx, q.key, data).Err()
}

func (q *Queue) Dequeue() (email.Message, error) {
	data, err := q.client.LPop(q.ctx, q.key).Result()
	if err != nil {
		return email.Message{}, err
	}

	msg, err := email.Deserialize([]byte(data))
	if err != nil {
		return email.Message{}, err
	}

	return *msg, nil

	//result, err := q.client.RPop(q.ctx, q.key).Result()
	//if errors.Is(err, redis.Nil) {
	//	return email.Message{}, nil // Очередь пуста
	//}
	//return email.Message{Body: result}, err
}
