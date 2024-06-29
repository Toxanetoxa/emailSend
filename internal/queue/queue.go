package queue

import (
	"email-sendler/internal/queueTypes"
	"email-sendler/internal/redis"
	"fmt"
)

// Type определяет тип очереди.
type Type int

const (
	RedisQueueType Type = iota
	KafkaQueueType
	RabbitMQQueueType
)

type Factory struct {
}

// CreateQueue создает новую очередь.
func (f *Factory) CreateQueue(queueType Type, config map[string]interface{}) (queueTypes.Queue, error) {
	switch queueType {
	case RedisQueueType:
		return redis.NewRedisQueue(
			config["addr"].(string),
			config["password"].(string),
			config["db"].(int),
			config["key"].(string),
		), nil
	case KafkaQueueType:
		//TODO реализовать логику для Kafka
		return nil, fmt.Errorf("реализовать для Kafka")
	case RabbitMQQueueType:
		//TODO реализовать логику для rabbitMQ
		return nil, fmt.Errorf("реализовать для rabbitMQ")
	default:
		return nil, fmt.Errorf("неизвестный тип очереди")
	}
}
