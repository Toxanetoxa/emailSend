package email

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/gomail.v2"
	"net/mail"
)

// Определение базовых ошибок
var (
	ErrInvalidRecipient = errors.New("invalid recipient address")
	ErrSendFailed       = errors.New("failed to send email")
)

// NewInvalidRecipientError Функция для создания ошибок с дополнительным контекстом
func NewInvalidRecipientError(recipient string) error {
	return fmt.Errorf("%w: %s", ErrInvalidRecipient, recipient)
}

// NewSendFailedError Функция для создания ошибок с дополнительным контекстом
func NewSendFailedError(err error) error {
	return fmt.Errorf("%w: %v", ErrSendFailed, err)
}

// Sender представляет интерфейс для отправки email.
type Sender interface {
	Send(ctx context.Context, to string, subject string, body string) error
}

// SenderConf представляет реальный интерфейс для отправки email.
type SenderConf struct {
	Host     string
	Port     int
	Username string
	Password string
}

// Message представляет структуру email сообщения.
type Message struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

// Serialize сериализует сообщение в JSON.
func (m *Message) Serialize() ([]byte, error) {
	return json.Marshal(m)
}

// Deserialize десериализует JSON в сообщение.
func Deserialize(data []byte) (*Message, error) {
	var msg Message
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

// Send метод отправления email.
func (r *SenderConf) Send(ctx context.Context, to string, subject string, body string) error {
	if !isValidEmail(to) {
		return NewInvalidRecipientError(to)
	}

	m := gomail.NewMessage()
	m.SetHeader("From", r.Username)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer(r.Host, r.Port, r.Username, r.Password)

	// Используйте контекст для отправки email.
	done := make(chan error, 1)
	go func() {
		done <- d.DialAndSend(m)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-done:
		if err != nil {
			return NewSendFailedError(err)
		}
		return nil
	}
}

// isValidEmail Валидация email
func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
