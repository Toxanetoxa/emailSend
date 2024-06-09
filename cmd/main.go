package main

import (
	"email-sendler/internal/dispatcher"
	"email-sendler/internal/email"
	"email-sendler/internal/queue"
	"fmt"
	"os"
	"time"
)

func main() {
	stage := os.Getenv("STAGE")
	if stage == "local" {
		sender := &email.MockEmailSender{}
		q := queue.NewInMemoryQueue()
		emailDispatcher := dispatcher.NewEmailDispatcherMock(*sender, q, 5, time.Minute)

		// Добавим несколько сообщений в очередь
		err := q.Enqueue(email.Message{To: "toxanetoxa@gmail.com", Subject: "Subject 1", Body: "Body 1"})
		if err != nil {
			return
		}
		err = q.Enqueue(email.Message{To: "toxanetoxa@gmail.com", Subject: "Subject 2", Body: "Body 2"})
		if err != nil {
			return
		}
		err = q.Enqueue(email.Message{To: "toxanetoxa@gmail.com", Subject: "Subject 3", Body: "Body 3"})
		if err != nil {
			return
		}

		// Запускаем процесс отправки email
		emailDispatcher.StartMock()
	} else {
		fmt.Println("Real email sender")
		Host := os.Getenv("HOST")
		Port := 587 // Обычно используется этот порт для SMTP с TLS
		Username := os.Getenv("USERNAME")
		Password := os.Getenv("PASSWORD")

		senderReal := &email.RealEmailSender{
			Host:     Host,
			Port:     Port,
			Username: Username,
			Password: Password,
		}

		realQ := queue.NewInMemoryQueue()
		realDispatcher := dispatcher.NewEmailDispatcher(*senderReal, realQ, 5, time.Minute)

		err := realQ.Enqueue(email.Message{To: "toxanetoxa@gmail.com", Subject: "Subject 1", Body: "Body 1"})
		if err != nil {
			return
		}
		err = realQ.Enqueue(email.Message{To: "toxanetoxa@gmail.com", Subject: "Subject 1", Body: "Body 1"})
		if err != nil {
			return
		}

		realDispatcher.Start()
	}
}
