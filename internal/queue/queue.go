package queue

import (
	"email-sendler/internal/kafka"
	"email-sendler/internal/queueTypes"
	"email-sendler/internal/rabbitmq"
	"email-sendler/internal/redis"
	"fmt"
)

// Queue представляет интерфейс для работы с очередью сообщений.

// QueueType определяет тип очереди.
type QueueType int

const (
	RedisQueueType QueueType = iota
	KafkaQueueType
	RabbitMQQueueType
)

type Factory struct {
}

// CreateQueue создает новую очередь.
func (f *Factory) CreateQueue(queueType QueueType, config map[string]interface{}) (queueTypes.Queue, error) {
	switch queueType {
	case RedisQueueType:
		return redis.NewRedisQueue(
			config["addr"].(string),
			config["password"].(string),
			config["db"].(int),
			config["key"].(string),
		), nil
	case KafkaQueueType:
		return kafka.NewKafkaQueue(
			config["brokers"].([]string),
			config["topic"].(string),
		), nil
	case RabbitMQQueueType:
		return rabbitmq.NewRabbitMQQueue(
			config["url"].(string),
			config["queueName"].(string),
		)
	default:
		return nil, fmt.Errorf("неизвестный тип очереди")
	}
}
