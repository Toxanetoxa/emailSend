package email

import (
	"encoding/json"
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
