package dispatcher

import (
	"context"
	"email-sendler/internal/email"
	"email-sendler/internal/queueTypes"
	"errors"
	"fmt"
	"time"
)

// EmailDispatcher отвечает за отправку real email из очереди.
type EmailDispatcher struct {
	sender   email.Sender
	queue    queueTypes.Queue
	limit    int
	interval time.Duration
}

// NewEmailDispatcher создает новый EmailDispatcher.
func NewEmailDispatcher(sender email.Sender, queue queueTypes.Queue, limit int, interval time.Duration) *EmailDispatcher {
	return &EmailDispatcher{
		sender:   sender,
		queue:    queue,
		limit:    limit,
		interval: interval,
	}
}

// Start запускает процесс отправки email.
func (d *EmailDispatcher) Start(stopChan chan struct{}) {
	const op = "EmailDispatcher.Start"

	fmt.Println("Начало отправки")
	ticker := time.NewTicker(d.interval)
	defer ticker.Stop()
	count := 0

	messageChan := make(chan email.Message)
	errChan := make(chan error)

	// Горутина для извлечения сообщений
	go func() {
		for {
			select {
			case <-stopChan:
				fmt.Println("Stopping queue reading")
				return
			default:
				if d.queue.Len() == 0 {
					time.Sleep(1 * time.Second) // Можно добавить задержку перед повторной попыткой
					continue
				}

				message, err := d.queue.Dequeue()
				if err != nil {
					errChan <- err
					continue
				}
				messageChan <- message
			}
		}
	}()

	for {
		select {
		case <-ticker.C:
			if count >= d.limit {
				fmt.Println("Достигнут лимит отправки")
				return
			}

		case message := <-messageChan:
			fmt.Printf("Отправка сообщения: %+v\n", message)

			// Создание контекста с таймаутом для отправки
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

			if err := d.sender.Send(ctx, message.To, message.Subject, message.Body); err == nil {
				count++
			} else {
				if errors.Is(err, email.ErrInvalidRecipient) {
					fmt.Printf("%v. Invalid recipient: %v\n", op, err)
				} else if errors.Is(err, email.ErrSendFailed) {
					fmt.Printf("%v. Error sending email: %v\n", op, err)
				} else {
					fmt.Printf("%v. Unknown error: %v\n", op, err)
				}
			}

			cancel()

		case err := <-errChan:
			fmt.Printf("%v. Error Dequeue message: %v \n", op, err)

		case <-stopChan:
			fmt.Println("Stopping dispatcher")
			return
		}
	}
}
