package main

import (
	"email-sendler/internal/dispatcher"
	"email-sendler/internal/email"
	"email-sendler/internal/emailLogger"
	"email-sendler/internal/queue"
	"email-sendler/internal/queueTypes"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	logger, err := emailLogger.NewLogger("error.log")
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()

	err = godotenv.Load("../.env.prod")
	if err != nil {
		//log.Fatalf("Error loading .env.prod file: %v", err)
		logger.Error("main: Error loading .env.prod file: %v", err)
	}

	STMTHost := os.Getenv("STMT_HOST")
	STMTUsername := os.Getenv("STMT_USERNAME")
	STMTPassword := os.Getenv("STMT_PASSWORD")
	STMTPort, err := strconv.Atoi(os.Getenv("STMT_PORT"))
	if err != nil {
		//log.Fatalf("PORT must be an integer")
		logger.Error("main: PORT must be an integer", err)
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
		//log.Fatalf("Error creating Redis queue: %v", err)
		logger.Error("main: Error creating Redis queue", err)
		return
	}

	//////--- тестовое сообщение ---
	msg := email.Message{
		To:      "toxanetoxa@gmail.",
		Subject: "Test1 Subject",
		Body:    "Test sdfsafdsfsdfdsfsdfdsf",
	}

	//// ---- Добавление тестового сообщения в очередь Redis---
	for i := 0; i < 4; i++ {
		err = redisQue.Enqueue(msg)
		if err != nil {
			//log.Fatalf("Error enqueuing message: %v", err)
			logger.Error("main: Error enqueuing message", err)
			return
		}
	}

	// Создание диспетчера который будет отправлять сообщения из очереди
	emailDispatcher := dispatcher.NewEmailDispatcher(sender, redisQue, 10, time.Second)

	// Канал для остановки диспетчера
	stopChan := make(chan struct{})
	go emailDispatcher.Start(stopChan)
	// Пример остановки диспетчера через 10 секунд
	go func() {
		time.Sleep(10 * time.Second)
		close(stopChan)
	}()

	// Блокировка основного потока, чтобы программа не завершалась
	select {}
}

func CreateRedisQue() (queueTypes.Queue, error) {
	const op = "main.SendRedisQue"

	RedisHost := os.Getenv("REDIS_HOST")
	RedisPort := os.Getenv("REDIS_PORT")
	RedisPassword := os.Getenv("REDIS_PASSWORD")
	RedisKey := os.Getenv("REDIS_KEY_PREFIX")
	RedisDB, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		log.Fatalf("%v. Error converting REDIS_DB to int: %v", op, err)
		return nil, err
	}

	redisConfig := map[string]interface{}{
		"addr":     RedisHost + ":" + RedisPort,
		"password": RedisPassword,
		"db":       RedisDB,
		"key":      RedisKey,
	}

	factory := &queue.Factory{}
	redisQueue, err := factory.CreateQueue(queue.RedisQueueType, redisConfig)
	if err != nil {
		log.Fatalf("%v. Error creating Redis queue: %v", op, err)
		return nil, err
	}

	fmt.Printf("Successfully created Redis queue: %s\n", redisConfig["key"])
	return redisQueue, nil
}
