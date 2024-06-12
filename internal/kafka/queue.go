package kafka

import (
	"context"
	"email-sendler/internal/email"
	"email-sendler/internal/queueTypes"
	"github.com/segmentio/kafka-go"
)

type Queue struct {
	writer *kafka.Writer
	reader *kafka.Reader
	ctx    context.Context
}

func NewKafkaQueue(brokers []string, topic string) queueTypes.Queue {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
	})

	return &Queue{
		writer: writer,
		reader: reader,
		ctx:    context.Background(),
	}
}

func (q *Queue) Enqueue(message email.Message) error {
	return q.writer.WriteMessages(q.ctx, kafka.Message{
		Value: []byte(message.Body),
	})
}

func (q *Queue) Dequeue() (email.Message, error) {
	msg, err := q.reader.ReadMessage(q.ctx)
	if err != nil {
		return email.Message{}, err
	}
	return email.Message{Body: string(msg.Value)}, nil
}
