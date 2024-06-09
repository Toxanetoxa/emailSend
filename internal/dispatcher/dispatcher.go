package dispatcher

import (
	"email-sendler/internal/email"
	"email-sendler/internal/queue"
	"time"
)

// EmailDispatcherMock отвечает за отправку mock email из очереди.
type EmailDispatcherMock struct {
	sender   email.MockEmailSender
	queue    queue.Queue
	limit    int
	interval time.Duration
}

// NewEmailDispatcherMock создает новый EmailDispatcherMock.
func NewEmailDispatcherMock(sender email.MockEmailSender, queue queue.Queue, limit int, interval time.Duration) *EmailDispatcherMock {
	return &EmailDispatcherMock{
		sender:   sender,
		queue:    queue,
		limit:    limit,
		interval: interval,
	}
}

// StartMock запускает процесс отправки email.
func (d *EmailDispatcherMock) StartMock() {
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

// EmailDispatcher отвечает за отправку real email из очереди.
type EmailDispatcher struct {
	sender   email.RealEmailSender
	queue    queue.Queue
	limit    int
	interval time.Duration
}

// NewEmailDispatcher создает новый EmailDispatcher.
func NewEmailDispatcher(sender email.RealEmailSender, queue queue.Queue, limit int, interval time.Duration) *EmailDispatcher {
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

		if err := d.sender.Send(message.To, message.Subject, message.Body); err == nil {
			count++
		}
	}
}
