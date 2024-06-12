package main

import (
	"email-sendler/internal/dispatcher"
	"email-sendler/internal/email"
	"email-sendler/internal/queueTypes"
	"email-sendler/internal/redis"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	err := godotenv.Load("../.env.prod")
	if err != nil {
		log.Fatalf("Error loading .env.prod file: %v", err)
	}

	STMTHost := os.Getenv("STMT_HOST")
	STMTUsername := os.Getenv("STMT_USERNAME")
	STMTPassword := os.Getenv("STMT_PASSWORD")
	STMTPort, err := strconv.Atoi(os.Getenv("STMT_PORT"))
	if err != nil {
		log.Fatalf("PORT must be an integer")
	}

	sender := &email.SenderConf{
		Host:     STMTHost,
		Port:     STMTPort,
		Username: STMTUsername,
		Password: STMTPassword,
	}

	// Создание очереди Redis
	redisQue, err := CreateRedisQue()
	if err != nil {
		log.Fatalf("Error creating Redis queue: %v", err)
		return
	}

	////--- тестовое сообщение ---
	//msg := email.Message{
	//	To:      "toxanetoxa@gmail.com",
	//	Subject: "Test1 Subject",
	//	Body:    "Test sdfsafdsfsdfdsfsdfdsf",
	//}
	//
	//// ---- Добавление тестового сообщения в очередь Redis---
	//for i := 0; i < 10; i++ {
	//	err = redisQue.Enqueue(msg)
	//	if err != nil {
	//		log.Fatalf("Error enqueuing message: %v", err)
	//		return
	//	}
	//}

	// Создание диспетчера который будет отправлять сообщения из очереди
	// Реализация с редисом
	emailDispatcher := dispatcher.NewEmailDispatcher(sender, redisQue, 10, time.Second)
	go emailDispatcher.Start()
	select {}

	// --- Проверка отправки сообщения (РАБОТАЕТ!!!) ---
	//err = sender.Send("toxanetoxa@gmail.com", "test 1", "Body 1")
	//if err != nil {
	//	log.Fatalf("Error sending email: %v", err)
	//	return
	//}

	// ---Тестирование подключение, добавление, удаление сообщений из очереди ---
	//testingQue.TestRedis()
	//testingQue.TestKafka()
	//testingQue.TestRabbitMQ()
}

func CreateRedisQue() (queueTypes.Queue, error) {
	const op = "main.SendRedisQue"

	RedisHost := os.Getenv("REDIS_HOST")
	RedisPort := os.Getenv("REDIS_PORT")
	RedisPassword := os.Getenv("REDIS_PASSWORD")
	RedisDB, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		log.Fatalf("%v. Error converting REDIS_DB to int: %v", op, err)
		return nil, err
	}
	RedisKey := os.Getenv("REDIS_KEY_PREFIX")

	redisConfig := map[string]interface{}{
		"addr":     RedisHost + ":" + RedisPort,
		"password": RedisPassword,
		"db":       RedisDB,
		"key":      RedisKey,
	}

	redisQue := redis.NewRedisQueue(
		redisConfig["addr"].(string),
		redisConfig["password"].(string),
		redisConfig["db"].(int),
		redisConfig["key"].(string),
	)

	return redisQue, nil
}
