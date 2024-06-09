package queue

import (
	"email-sendler/internal/email"
	"errors"
	"sync"
)

// Queue представляет интерфейс для работы с очередью сообщений.
type Queue interface {
	Enqueue(message email.Message) error
	Dequeue() (email.Message, error)
}

// InMemoryQueue представляет простую очередь сообщений в памяти (для тестирования).
type InMemoryQueue struct {
	messages []email.Message
	mutex    sync.Mutex
}

// NewInMemoryQueue создает новую очередь.
func NewInMemoryQueue() *InMemoryQueue {
	return &InMemoryQueue{messages: make([]email.Message, 0)}
}

// Enqueue добавляет сообщение в очередь.
func (q *InMemoryQueue) Enqueue(message email.Message) error {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	q.messages = append(q.messages, message)
	return nil
}

// Dequeue извлекает сообщение из очереди.
func (q *InMemoryQueue) Dequeue() (email.Message, error) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	if len(q.messages) == 0 {
		return email.Message{}, errors.New("queue is empty")
	}
	message := q.messages[0]
	q.messages = q.messages[1:]
	return message, nil
}
