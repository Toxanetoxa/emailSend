package email

import (
	"context"
	"gopkg.in/gomail.v2"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIsValidEmail(t *testing.T) {
	validEmails := []string{
		"test@example.com",
		"user.name+tag+sorting@example.com",
		"admin@mailserver1",
	}
	invalidEmails := []string{
		"plainaddress",
		"@missingusername.com",
		"username@.com",
	}

	for _, email := range validEmails {
		assert.True(t, isValidEmail(email), "Expected email to be valid: %s", email)
	}
	for _, email := range invalidEmails {
		assert.False(t, isValidEmail(email), "Expected email to be invalid: %s", email)
	}
}

func TestSerializeDeserializeMessage(t *testing.T) {
	const op = "sender_test.TestSerializeDeserializeMessage"
	msg := &Message{
		To:      "test@example.com",
		Subject: "Test Subject",
		Body:    "Test Body",
	}

	data, err := msg.Serialize()
	assert.NoError(t, err)

	deserializedMsg, err := Deserialize(data)
	assert.NoError(t, err)
	assert.Equal(t, msg, deserializedMsg)
}

func TestSend_InvalidEmail(t *testing.T) {
	const op = "sender_test.TestSend_InvalidEmail"
	ctx := context.Background()
	sender := &SenderConf{
		Host:     "smtp.example.com",
		Port:     587,
		Username: "user@example.com",
		Password: "password",
	}

	err := sender.Send(ctx, "invalid-email", "Test Subject", "Test Body")
	assert.ErrorIs(t, err, ErrInvalidRecipient)
	assert.Equal(t, "email.Send", op)
}

// MockDialer is a mock for gomail.Dialer
type MockDialer struct{}

func (d *MockDialer) DialAndSend(m ...*gomail.Message) error {
	return nil
}

func TestSend_ValidEmail(t *testing.T) {
	const op = "sender_test.TestSend_ValidEmail"
	ctx := context.Background()
	sender := &SenderConf{
		Host:     "smtp.example.com",
		Port:     587,
		Username: "user@example.com",
		Password: "password",
	}

	// Use a short timeout context to test context cancellation
	ctx, cancel := context.WithTimeout(ctx, 1*time.Millisecond)
	defer cancel()

	go func() {
		time.Sleep(2 * time.Millisecond)
		cancel()
	}()

	err := sender.Send(ctx, "test@example.com", "Test Subject", "Test Body")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "context deadline exceeded")
	assert.Equal(t, "email.Send", op)
}
