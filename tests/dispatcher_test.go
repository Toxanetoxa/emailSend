package tests_test

import (
	"email-sendler/internal/dispatcher"
	"email-sendler/internal/email"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type MockSender struct {
	sendCalled bool
}

func (m *MockSender) Send(to, subject, body string) error {
	m.sendCalled = true
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
	return email.Message{}, nil
}

func TestEmailDispatcher_Start(t *testing.T) {
	mockSender := &MockSender{}
	mockQueue := &MockQueue{}
	emailDispatcher := dispatcher.NewEmailDispatcher(mockSender, mockQueue, 5, time.Second)

	go func() {
		time.Sleep(2 * time.Second)
		// Просто завершаем цикл, чтобы тест завершился.
		return
	}()

	emailDispatcher.Start()

	assert.True(t, mockSender.sendCalled, "Expected sender to be called")
	assert.Equal(t, 5, mockQueue.dequeueCount, "Expected dequeue to be called 5 times")
}
