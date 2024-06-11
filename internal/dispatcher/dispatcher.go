package dispatcher

import (
	"email-sendler/internal/email"
	"email-sendler/internal/queue"
	"fmt"
	"time"
)

// EmailDispatcher отвечает за отправку real email из очереди.
type EmailDispatcher struct {
	sender   email.Sender
	queue    queue.Queue
	limit    int
	interval time.Duration
}

// NewEmailDispatcher создает новый EmailDispatcher.
func NewEmailDispatcher(sender email.Sender, queue queue.Queue, limit int, interval time.Duration) *EmailDispatcher {
	return &EmailDispatcher{
		sender:   sender,
		queue:    queue,
		limit:    limit,
		interval: interval,
	}
}

// Start запускает процесс отправки email.
func (d *EmailDispatcher) Start() {
	fmt.Println("Начало отправки")
	ticker := time.NewTicker(d.interval)
	defer ticker.Stop()
	count := 0

	for range ticker.C {
		if count >= d.limit {
			count = 0
		}

		message, err := d.queue.Dequeue()
		if err != nil {
			fmt.Println(err)
			continue
		}

		if err := d.sender.Send(message.To, message.Subject, message.Body); err == nil {
			count++
		}
	}
}
