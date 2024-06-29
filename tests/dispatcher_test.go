package tests_test

import (
	"context"
	"email-sendler/internal/dispatcher"
	"email-sendler/internal/email"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type MockSender struct {
	sendCalledCount int
}

func (m *MockSender) Send(ctx context.Context, to, subject, body string) error {
	m.sendCalledCount++
	return nil
}

type MockQueue struct {
	dequeueCount int
}

func (m *MockQueue) Enqueue(message email.Message) error {
	return nil
}

func (m *MockQueue) Dequeue() (email.Message, error) {
	m.dequeueCount++
	if m.dequeueCount <= 5 {
		return email.Message{To: "test@example.com", Subject: "Test", Body: "Test Body"}, nil
	}
	return email.Message{}, nil
}

func (m *MockQueue) Len() int {
	// Реализация метода Len для тестовых нужд
	return 0
}

func TestEmailDispatcher_Start(t *testing.T) {
	mockSender := &MockSender{}
	mockQueue := &MockQueue{}
	emailDispatcher := dispatcher.NewEmailDispatcher(mockSender, mockQueue, 5, time.Millisecond*100)
	stopChan := make(chan struct{})

	go func() {
		time.Sleep(time.Second) // Даем достаточно времени для выполнения 5 итераций
		close(stopChan)         // Останавливаем диспетчер
	}()

	emailDispatcher.Start(stopChan)

	assert.Equal(t, 5, mockSender.sendCalledCount, "Expected sender to be called 5 times")
	assert.Equal(t, 5, mockQueue.dequeueCount, "Expected dequeue to be called 5 times")
}
