package dispatcher

import (
	"email-sendler/internal/email"
	"email-sendler/internal/queue"
	"time"
)

// EmailDispatcher отвечает за отправку email из очереди.
type EmailDispatcher struct {
	sender   email.MockEmailSender
	queue    queue.Queue
	limit    int
	interval time.Duration
}

// NewEmailDispatcher создает новый EmailDispatcher.
func NewEmailDispatcher(sender email.MockEmailSender, queue queue.Queue, limit int, interval time.Duration) *EmailDispatcher {
	return &EmailDispatcher{
		sender:   sender,
		queue:    queue,
		limit:    limit,
		interval: interval,
	}
}

// Start запускает процесс отправки email.
func (d *EmailDispatcher) Start() {
	ticker := time.NewTicker(d.interval)
	defer ticker.Stop()
	count := 0

	for range ticker.C {
		if count >= d.limit {
			count = 0
		}

		message, err := d.queue.Dequeue()
		if err != nil {
			continue
		}

		if err := d.sender.MockSend(message.To, message.Subject, message.Body); err == nil {
			count++
		}
	}
}
