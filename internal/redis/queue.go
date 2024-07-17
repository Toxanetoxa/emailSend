package redis

import (
	"context"
	"email-sendler/internal/email"
	"fmt"

	"github.com/go-redis/redis/v8"
)

type Queue struct {
	Client *redis.Client
	Ctx    context.Context
	Key    string
}

type QueueInterface interface {
	Len() int
	Dequeue() (email.Message, error)
	Enqueue(message email.Message) error
}

// Config интерфейс для Redis
type Config interface {
	GetAddr() string
	GetPassword() string
	GetDB() int
	GetKey() string
}

// RedisConfig redisConfig структура которая имплементирует Config
type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
	Key      string
}

func (r *RedisConfig) GetAddr() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}

func (r *RedisConfig) GetPassword() string {
	return r.Password
}

func (r *RedisConfig) GetDB() int {
	return r.DB
}

func (r *RedisConfig) GetKey() string {
	return r.Key
}

// NewRedisQueue создание новой очериди в Redis
func NewRedisQueue(config Config) *Queue {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.GetAddr(),
		Password: config.GetPassword(),
		DB:       config.GetDB(),
	})

	return &Queue{
		Client: rdb,
		Ctx:    context.Background(),
		Key:    config.GetKey(),
	}
}

// Enqueue добавление message в очередь
func (q *Queue) Enqueue(message email.Message) error {
	data, err := message.Serialize()
	if err != nil {
		return err
	}

	return q.Client.RPush(q.Ctx, q.Key, data).Err()
}

// Dequeue удаление message из очереди
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

// Len возращает длинну очереди
func (q *Queue) Len() int {
	l, err := q.Client.LLen(context.Background(), "queue").Result()
	if err != nil {
		return 0
	}
	return int(l)
}

// Ensure Queue implements Queue
var _ QueueInterface = (*Queue)(nil)
