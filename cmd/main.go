package main

import (
	"email-sendler/internal/dispatcher"
	"email-sendler/internal/email"
	"email-sendler/internal/logger"
	"email-sendler/internal/redis"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

// initSenderConf метод который инициализация конфиг для сервиса STMT
func initSenderConf() (error, *email.SenderConf, string) {
	const op = "main.initSenderConf"
	var err error
	var STMTHost, STMTUsername, STMTPassword string
	var STMTPort int

	err = godotenv.Load("../.env.prod")
	if err != nil {
		return err, nil, op
	}

	STMTHost = os.Getenv("STMT_HOST")
	STMTUsername = os.Getenv("STMT_USERNAME")
	STMTPassword = os.Getenv("STMT_PASSWORD")
	STMTPort, err = strconv.Atoi(os.Getenv("STMT_PORT"))
	if err != nil {
		return err, nil, op
	}

	conf := &email.SenderConf{
		Host:     STMTHost,
		Port:     STMTPort,
		Username: STMTUsername,
		Password: STMTPassword,
	}

	return nil, conf, op
}

// initLogger метод который инициализация логер для сервиса
func initLogger() (error, *logger.LoggerFile, string) {
	const op = "main.initLogger"

	err := godotenv.Load("../.env.prod")
	if err != nil {
		return err, nil, op
	}
	LoggerFile := os.Getenv("LOGGER_FILE")
	logger, err := logger.NewLogger(LoggerFile)
	if err != nil {
		return err, nil, op
	}

	return nil, logger, op
}

// initServer метод который инициализация сервера
func initServer() (error, *chi.Mux, string) {
	const op = "main.initServer"

	var err error
	var PORT string
	var router *chi.Mux

	err = godotenv.Load("../.env.prod")
	if err != nil {
		return err, nil, op
	}

	PORT = os.Getenv("SERVER_PORT")

	router = chi.NewRouter()
	// запуск сервера
	go func() {
		err := http.ListenAndServe(fmt.Sprint(":", PORT), router)
		if err != nil {
			log.Printf("Error starting server: %v", err)
			return // если сервер не запустился, то программа завершается
		}
	}()

	fmt.Printf("Server started on port %s\n", PORT)
	return nil, router, op
}

// initRedisQue метод который инициализация очередь Redis
func initRedisQue() (string, error, *redis.Queue) {
	const op = "main.SendRedisQue"

	var RedisPort, RedisDB int
	var RedisHost, RedisPassword, RedisKey string
	var err error

	RedisHost = os.Getenv("REDIS_HOST")
	RedisPassword = os.Getenv("REDIS_PASSWORD")
	RedisKey = os.Getenv("REDIS_KEY_PREFIX")
	RedisPort, err = strconv.Atoi(os.Getenv("REDIS_PORT"))
	if err != nil {
		return op, err, nil
	}

	RedisDB, err = strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		log.Fatalf("%v. Error converting REDIS_DB to int: %v", op, err)
		return op, err, nil
	}

	config := &redis.RedisConfig{
		Host:     RedisHost,
		Port:     RedisPort,
		Password: RedisPassword,
		DB:       RedisDB,
		Key:      RedisKey,
	}

	queue := redis.NewRedisQueue(config)

	if err != nil {
		log.Fatalf("%v. Error creating Redis queue: %v", op, err)
		return op, err, nil
	}

	return op, nil, queue
}

func main() {
	var op string
	// Инициализация логера
	err, logger, op := initLogger()
	if err != nil {
		log.Fatalf(op, err)
	}
	defer logger.Close()

	// Инициализация сервера
	var router *chi.Mux
	err, router, op = initServer()
	if err != nil {
		logger.Error(op, err)
		return
	}

	// Инициализация middleware
	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	// Инициализация обработчиков
	router.Route("/api", func(r chi.Router) {})

	// Иницализация конфига
	var sender *email.SenderConf
	err, sender, op = initSenderConf()
	if err != nil {
		logger.Error(op, err)
		return
	}

	// Инициализация очереди в Redis
	var redisQue *redis.Queue
	op, err, redisQue = initRedisQue()
	if err != nil {
		logger.Error(op, err)
		return
	}

	// TODO сделать ручку для принятия сообщений которые нужно отправить
	//////--- тестовое сообщение ---
	//msg := email.Message{
	//	To:      "toxanetoxa@gmail.",
	//	Subject: "Test1 Subject",
	//	Body:    "Test sdfsafdsfsdfdsfsdfdsf",
	//}
	//// ---- Добавление тестового сообщения в очередь Redis---
	//for i := 0; i < 4; i++ {
	//	err = redisQue.Enqueue(msg)
	//	if err != nil {
	//		//log.Fatalf("Error enqueuing message: %v", err)
	//		logger.Error("main: Error enqueuing message", err)
	//		return
	//	}
	//}

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
