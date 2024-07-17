package main

import (
	"email-sendler/internal/dispatcher"
	"email-sendler/internal/email"
	"email-sendler/internal/http-server/handler"
	"email-sendler/internal/logger"
	"email-sendler/internal/redis"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

type HTTPServer interface {
	GetAddr() string
	GetTimeout() time.Duration
	GetIdleTimeout() time.Duration
	GetHandler() *chi.Mux
}

type ServerConfig struct {
	Timeout     time.Duration
	IdleTimeout time.Duration
	Handler     *chi.Mux
}

type Address struct {
	Host string
	Port string
}

func (r *Address) GetAddr() string {
	return fmt.Sprintf("%s:%s", r.Host, r.Port)
}

func (r *ServerConfig) GetHandler() *chi.Mux {
	return r.Handler
}

func (r *ServerConfig) GetTimeout() time.Duration {
	return r.Timeout
}

func (r *ServerConfig) GetIdleTimeout() time.Duration {
	return r.IdleTimeout
}

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
func initLogger() (error, *logger.File, string) {
	const op = "main.initLogger"

	var err error
	var l *logger.File

	err = godotenv.Load("../.env.prod")
	if err != nil {
		return err, nil, op
	}
	LoggerFile := os.Getenv("LOGGER_FILE")
	l, err = logger.NewLogger(LoggerFile)
	if err != nil {
		return err, nil, op
	}

	return nil, l, op
}

// initRouter метод который инициализирует роутер
func initRouter() (error, *chi.Mux, string) {
	const op = "main.initRouter"

	router := chi.NewRouter()

	return nil, router, op
}

// initServer метод который инициализирует сервера
func initServer(router *chi.Mux) (error, *http.Server, string) {
	const op = "main.initServer"
	var HOST, PORT string
	var TIMEOUT, IdleTimeout time.Duration
	var err error
	var srvConf *ServerConfig

	err = godotenv.Load("../.env.prod")
	if err != nil {
		return err, nil, op
	}

	PORT = os.Getenv("SERVER_PORT")
	HOST = os.Getenv("SERVER_HOST")
	TIMEOUT, err = time.ParseDuration(os.Getenv("SERVER_TIMEOUT"))
	if err != nil {
		return err, nil, op
	}
	IdleTimeout, err = time.ParseDuration(os.Getenv("SERVER_IDLE_TIMEOUT"))
	if err != nil {
		return err, nil, op
	}

	address := &Address{
		Host: HOST,
		Port: PORT,
	}

	srvConf = &ServerConfig{
		Timeout:     TIMEOUT,
		IdleTimeout: IdleTimeout,
		Handler:     router,
	}

	srv := &http.Server{
		Addr:         address.GetAddr(),
		Handler:      srvConf.GetHandler(),
		ReadTimeout:  srvConf.GetTimeout(),
		WriteTimeout: srvConf.GetTimeout(),
		IdleTimeout:  srvConf.GetIdleTimeout(),
	}

	return nil, srv, op
}

// initRedisQue метод который инициализация очередь Redis
func initRedisQue() (string, error, *redis.Queue) {
	const op = "main.initRedisQue"

	var RedisPort, RedisDB int
	var RedisHost, RedisPassword, RedisKey string
	var err error

	err = godotenv.Load("../.env.prod")
	if err != nil {
		return op, err, nil
	}

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
	err, emailLogger, op := initLogger()
	if err != nil {
		log.Fatalf(op, err)
	}
	defer emailLogger.Close()

	// Инициализация роутера
	var router *chi.Mux
	err, router, op = initRouter()
	if err != nil {
		emailLogger.Error(op, err)
		return
	}

	// Инициализация middleware
	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	// Инициализация сервера
	var srv *http.Server

	err, srv, op = initServer(router)
	if err != nil {
		emailLogger.Error(op, err)
		return
	}

	// Запуск сервера
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("%s: Error starting server: %v", op, err)
			return
		}
	}()

	emailLogger.Info(op, "Server started")

	// Иницализация конфига sender
	var senderConf *email.SenderConf
	err, senderConf, op = initSenderConf()
	if err != nil {
		emailLogger.Error(op, err)
		return
	}

	// Инициализация очереди в Redis
	var redisQue *redis.Queue
	op, err, redisQue = initRedisQue()
	if err != nil {
		emailLogger.Error(op, err)
		return
	}

	// TODO сделать ручку для принятия сообщений которые нужно отправить
	// Инициализация обработчиков
	router.Route("/api",
		func(r chi.Router) {
			r.Post("/send", handler.New(emailLogger, redisQue))
		},
	)

	// Создание диспетчера который будет отправлять сообщения из очереди
	emailDispatcher := dispatcher.NewEmailDispatcher(senderConf, redisQue, 10, time.Second)

	// Канал для остановки диспетчера
	stopChan := make(chan struct{})

	// Пример остановки диспетчера через 10 секунд
	go func() {
		emailDispatcher.Start(stopChan, emailLogger)
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	close(stopChan)

	os.Exit(0)
}
