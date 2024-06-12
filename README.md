# Email Sender
### Email Sender - это проект на Go, который реализует компонент системы для отправки писем через очередь. Этот компонент может использоваться несколькими экземплярами приложения одновременно, обеспечивая stateless отправку писем.


## Особенности
### 1. Определение интерфейсов и структур
    Интерфейсы и структуры определены для отправки email и работы с очередью.
### 2. Реализация компонента отправки писем
    Используется библиотека для отправки email, предоставляющая возможность легко заменить её на любую другую.
### 3. Добавление очереди и подключения к брокеру
    Реализована очередь для хранения писем, без привязки к конкретному брокеру (Redis, RabbitMQ, Kafka и т.д.).
### 3. Добавление ограничения на количество отправленных писем
    Возможность установки лимита на количество отправляемых писем в минуту.
### 4. Написание тестов с использованием моков
    Полное покрытие тестами, использующими моки для отправки писем и работы с очередью.
### 5. Подготовка места для метрик
    Предусмотрено место для внедрения метрик в будущем (например, для Prometheus).

## Установка
1. Создайте .env.prod файл:

2. Перейдите в папку с Docker файлами:
```bash
cd docker
```
3. Создайте Docker сеть:
```bash
docker network create email-sender
```
4. Запустите Docker Compose:
```bash
docker-compose up -d
```

## Использование

#### Интерфейс email
```go
package email

type Message struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}
````

#### Интерфейс для отправки email
```go
package email

type Sender interface {
    Send(to string, subject string, body string) error
}
```

#### Интерфейс для работы с очередью
```go
package queue

type QueueType int

const (
	RedisQueueType QueueType = iota
	KafkaQueueType
	RabbitMQQueueType
)
```
#### Создание очереди через фабрику
```go
package main

 import "email-sendler/internal/queue"

func main()  {
	factory := &queue.Factory{}
	redisQueue, err := factory.CreateQueue(queue.RedisQueueType, redisConfig)
	if err != nil {...}
}
```