package tests

import (
	"email-sendler/internal/dispatcher"
	"email-sendler/internal/email"
	"email-sendler/internal/queue"
	"testing"
	"time"
)

func TestEmailDispatcher(t *testing.T) {
	sender := &email.MockEmailSender{}
	q := queue.NewInMemoryQueue()
	dispatcher := dispatcher.NewEmailDispatcher(sender, q, 5, time.Minute)

	q.Enqueue(email.Message{To: "test@example.com", Subject: "Test", Body: "This is a test email."})
	go dispatcher.Start()

	// Подождем немного, чтобы удостовериться, что сообщение было отправлено
	time.Sleep(2 * time.Minute)

	// Проверяем, что очередь пуста
	if _, err := q.Dequeue(); err == nil {
		t.Errorf("Queue should be empty after dispatching email")
	}
}
