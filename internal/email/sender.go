package email

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/mail"

	"gopkg.in/gomail.v2"
)

// Errors
var (
	ErrInvalidRecipient = errors.New("invalid recipient address")
	ErrSendFailed       = errors.New("failed to send email")
)

// NewInvalidRecipientError creates a new InvalidRecipientError
func NewInvalidRecipientError(recipient string) error {
	return fmt.Errorf("%w: %s", ErrInvalidRecipient, recipient)
}

// NewSendFailedError creates a new SendFailedError
func NewSendFailedError(err error) error {
	return fmt.Errorf("%w: %v", ErrSendFailed, err)
}

// Sender represents an email sender interface
type Sender interface {
	Send(ctx context.Context, to string, subject string, body string) error
}

// SenderConf represents the configuration for sending emails
type SenderConf struct {
	Host     string
	Port     int
	Username string
	Password string
}

// Message represents an email message
type Message struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

// Serialize serializes the message to JSON
func (m *Message) Serialize() ([]byte, error) {
	return json.Marshal(m)
}

// Deserialize deserializes JSON to a message
func Deserialize(data []byte) (*Message, error) {
	var msg Message
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

// Send sends an email
func (r *SenderConf) Send(ctx context.Context, to string, subject string, body string) error {
	const op = "email.Send"

	if !isValidEmail(to) {
		return NewInvalidRecipientError(to)
	}

	m := gomail.NewMessage()
	m.SetHeader("From", r.Username)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer(r.Host, r.Port, r.Username, r.Password)

	err := d.DialAndSend(m)
	if err != nil {
		return NewSendFailedError(err)
	}
	return nil
}

// isValidEmail validates an email address
func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
