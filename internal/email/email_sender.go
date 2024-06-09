package email

import (
	"fmt"
	"gopkg.in/gomail.v2"
)

// Sender представляет интерфейс для отправки email.
type Sender interface {
	Send(to string, subject string, body string) error
}

// RealEmailSender представляет реальный интерфейс для отправки email.
type RealEmailSender struct {
	Host     string
	Port     int
	Username string
	Password string
}

// Message представляет структуру email сообщения.
type Message struct {
	To      string
	Subject string
	Body    string
}

// MockEmailSender представляет мок для отправки email.
type MockEmailSender struct{}

// Send отправляет email.
func (m *MockEmailSender) MockSend(to string, subject string, body string) error {
	fmt.Printf("Email sent to %s with subject %s\n", to, subject)
	return nil
}

func (r *RealEmailSender) Send(to string, subject string, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", r.Username)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer(r.Host, r.Port, r.Username, r.Password)

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("could not send email: %v", err)
	}
	return nil
}
