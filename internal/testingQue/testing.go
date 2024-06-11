package testingQue

import (
	"context"
	"email-sendler/internal/email"
	"email-sendler/internal/queue"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
)

func TestRedis() {
	const op = "testingQue.TestRedis"

	redisConfig := map[string]interface{}{
		"addr":     "127.0.0.1:6379", // Обязательно проверьте, что здесь указан правильный порт
		"password": "",
		"db":       0,
		"key":      "emailQueue",
	}

	// Создание фабрики очередей
	factory := &queue.Factory{}

	// Создание очереди типа Redis
	redisQueue, err := factory.CreateQueue(queue.RedisQueueType, redisConfig)
	if err != nil {
		log.Fatalf("%v. Error creating Redis queue: %v", op, err)
	}

	// Проверка подключения к Redis
	client := redis.NewClient(&redis.Options{
		Addr:     redisConfig["addr"].(string),
		Password: redisConfig["password"].(string),
		DB:       redisConfig["db"].(int),
	})
	ctx := context.Background()

	pong, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("%v. Failed to connect to Redis: %v", op, err)
	}
	fmt.Printf("Connected to Redis: %v\n", pong)

	msg := email.Message{
		To:      "toxanetoxa@gmail.com",
		Subject: "Test Subject",
		Body:    "Test Body",
	}

	err = redisQueue.Enqueue(msg)
	if err != nil {
		log.Fatalf("%v. Error enqueuing message: %v", op, err)
	}
	fmt.Println("Message enqueued")

	// Проверка наличия ключа в очереди после добавления сообщения
	length, err := client.LLen(ctx, redisConfig["key"].(string)).Result()
	if err != nil {
		log.Fatalf("%v. Error checking queue length: %v", op, err)
	}
	fmt.Printf("Queue length after enqueue: %d\n", length)

	// Пример извлечения сообщения из очереди
	dequeuedMsg, err := redisQueue.Dequeue()
	if err != nil {
		log.Fatalf("%v. Error dequeuing message: %v", op, err)
	}

	fmt.Printf("Dequeued message: %+v\n", dequeuedMsg)
}

func TestKafka() {
	kafkaConfig := map[string]interface{}{
		"brokers": []string{"127.0.0.1:9092"},
		"topic":   "emailTopic",
	}

	factory := &queue.Factory{}

	kafkaQueue, err := factory.CreateQueue(queue.KafkaQueueType, kafkaConfig)
	if err != nil {
		log.Fatalf("Error creating Kafka queue: %v", err)
	}

	msg := email.Message{
		To:      "toxanetoxa@gmail.com",
		Subject: "Test Subject",
		Body:    "Test Body",
	}

	// Пример добавления сообщения в очередь Kafka
	err = kafkaQueue.Enqueue(msg)
	if err != nil {
		log.Fatalf("Error enqueuing message to Kafka: %v", err)
	}
	fmt.Println("Message enqueued to Kafka")

	// Пример извлечения сообщения из очереди Kafka
	dequeuedMsg, err := kafkaQueue.Dequeue()
	if err != nil {
		log.Fatalf("Error dequeuing message from Kafka: %v", err)
	}

	fmt.Printf("Dequeued message from Kafka: %+v\n", dequeuedMsg)
}
