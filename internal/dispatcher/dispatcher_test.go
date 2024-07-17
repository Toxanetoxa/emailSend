package dispatcher

import (
	"context"
	"email-sendler/internal/email"
	"email-sendler/internal/logger"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
)

// MockEmailSender mocks email.Sender interface
type MockEmailSender struct {
	mock.Mock
}

func (m *MockEmailSender) Send(ctx context.Context, to string, subject string, body string) error {
	args := m.Called(ctx, to, subject, body)
	return args.Error(0)
}

// MockQueue mocks redis.QueueInterface interface
type MockQueue struct {
	mock.Mock
}

func (m *MockQueue) Len() int {
	args := m.Called()
	return args.Int(0)
}

func (m *MockQueue) Dequeue() (email.Message, error) {
	args := m.Called()
	return args.Get(0).(email.Message), args.Error(1)
}

func (m *MockQueue) Enqueue(message email.Message) error {
	args := m.Called(message)
	return args.Error(0)
}

// MockLogger mocks logger.LoggerInterface interface
type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Success(op string, msg string) {
	m.Called(op, msg)
}

func TestEmailDispatcher_Start(t *testing.T) {
	sender := new(MockEmailSender)
	queue := new(MockQueue)
	l, err := logger.NewLogger("error.log")
	if err != nil {
		return
	}
	stopChan := make(chan struct{})
	dispatcher := NewEmailDispatcher(sender, queue, 5, 1*time.Second)

	message := email.Message{
		To:      "test@example.com",
		Subject: "Test Subject",
		Body:    "Test Body",
	}

	// Setting up the mock expectations
	l.Success("EmailDispatcher.Start", "Начало отправки сообщений")
	queue.On("Len").Return(1).Twice()
	queue.On("Dequeue").Return(message, nil).Twice()
	sender.On("Send", mock.Anything, message.To, message.Subject, message.Body).Return(nil).Twice()

	go dispatcher.Start(stopChan, l)

	// Wait for the dispatcher to process some emails
	time.Sleep(3 * time.Second)

	// Stop the dispatcher
	close(stopChan)
	time.Sleep(1 * time.Second)

	// Assert expectations
	queue.AssertExpectations(t)
	sender.AssertExpectations(t)
}

func TestEmailDispatcher_InvalidRecipient(t *testing.T) {
	sender := new(MockEmailSender)
	queue := new(MockQueue)
	l, err := logger.NewLogger("error.log")
	if err != nil {
		return
	}
	stopChan := make(chan struct{})
	dispatcher := NewEmailDispatcher(sender, queue, 5, 1*time.Second)

	invalidMessage := email.Message{
		To:      "invalid-email",
		Subject: "Test Subject",
		Body:    "Test Body",
	}

	// Setting up the mock expectations
	l.Success("EmailDispatcher.Start", "Начало отправки сообщений")
	queue.On("Len").Return(1).Twice()
	queue.On("Dequeue").Return(invalidMessage, nil).Twice()
	sender.On("Send", mock.Anything, invalidMessage.To, invalidMessage.Subject, invalidMessage.Body).Return(email.NewInvalidRecipientError(invalidMessage.To)).Twice()

	go dispatcher.Start(stopChan, l)

	// Wait for the dispatcher to process some emails
	time.Sleep(3 * time.Second)

	// Stop the dispatcher
	close(stopChan)
	time.Sleep(1 * time.Second)

	// Assert expectations
	queue.AssertExpectations(t)
	sender.AssertExpectations(t)
}
