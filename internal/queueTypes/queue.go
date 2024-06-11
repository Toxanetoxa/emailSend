package queueTypes

import "email-sendler/internal/email"

type Queue interface {
	Enqueue(message email.Message) error
	Dequeue() (email.Message, error)
}
