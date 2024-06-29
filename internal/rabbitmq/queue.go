package rabbitmq

//
//import (
//	"email-sendler/internal/email"
//	"email-sendler/internal/queueTypes"
//	"github.com/streadway/amqp"
//)
//
//type Queue struct {
//	conn    *amqp.Connection
//	channel *amqp.Channel
//	queue   amqp.Queue
//}
//
//func NewRabbitMQQueue(url, queueName string) (queueTypes.Queue, error) {
//	conn, err := amqp.Dial(url)
//	if err != nil {
//		return nil, err
//	}
//
//	channel, err := conn.Channel()
//	if err != nil {
//		return nil, err
//	}
//
//	queue, err := channel.QueueDeclare(
//		queueName,
//		false,
//		false,
//		false,
//		false,
//		nil,
//	)
//	if err != nil {
//		return nil, err
//	}
//
//	return &Queue{
//		conn:    conn,
//		channel: channel,
//		queue:   queue,
//	}, nil
//}
//
//func (q *Queue) Enqueue(message email.Message) error {
//	return q.channel.Publish(
//		"",
//		q.queue.Name,
//		false,
//		false,
//		amqp.Publishing{
//			ContentType: "text/plain",
//			Body:        []byte(message.Body),
//		},
//	)
//}
//
//func (q *Queue) Dequeue() (email.Message, error) {
//	msgs, err := q.channel.Consume(
//		q.queue.Name,
//		"",
//		true,
//		false,
//		false,
//		false,
//		nil,
//	)
//	if err != nil {
//		return email.Message{}, err
//	}
//
//	for msg := range msgs {
//		return email.Message{Body: string(msg.Body)}, nil
//	}
//
//	return email.Message{}, nil
//}
//
//func (q *Queue) Len() (int, error) {
//	// Реализация метода Len для RabbitMQ очереди
//	queueStats, err := q.channel.QueueInspect(q.queue.Name)
//	if err != nil {
//		return 0, err
//	}
//	return queueStats.Messages, nil
//}
