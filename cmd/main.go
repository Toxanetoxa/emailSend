package main

import (
	"email-sendler/internal/dispatcher"
	"email-sendler/internal/email"
	"email-sendler/internal/queue"
	"time"
)

func main() {
	sender := &email.MockEmailSender{}
	q := queue.NewInMemoryQueue()
	emailDispatcher := dispatcher.NewEmailDispatcher(sender, q, 5, time.Minute)

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
	emailDispatcher.Start()
}
