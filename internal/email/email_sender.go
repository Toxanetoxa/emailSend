package email

import (
	"fmt"
	"gopkg.in/gomail.v2"
)

// Sender представляет интерфейс для отправки email.
type Sender interface {
	Send(to string, subject string, body string) error
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
	To      string
	Subject string
	Body    string
}

// Send метод отправления email.
// Sender is an interface that defines the contract for sending emails.
// It requires a Send method that takes a recipient email address, a subject, and a body, all as strings.
func (r *SenderConf) Send(to string, subject string, body string) error {
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
