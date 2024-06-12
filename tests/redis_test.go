package tests

import (
	"context"
	"email-sendler/internal/redis"
	"testing"

	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"

	"email-sendler/internal/email"
)

func TestQueue_Enqueue(t *testing.T) {
	const queueKey = "emailQueue"
	msg := email.Message{
		To:      "toxanetoxa@gmail.com",
		Subject: "Test Subject",
		Body:    "Test Body",
	}

	serializedMessage, err := msg.Serialize()
	assert.NoError(t, err)

	db, mock := redismock.NewClientMock()
	mock.ExpectRPush(queueKey, serializedMessage).SetVal(1)

	q := &redis.Queue{
		Client: db,
		Ctx:    context.Background(),
		Key:    queueKey,
	}

	err = q.Enqueue(msg)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

}

func TestQueue_Dequeue(t *testing.T) {
	const queueKey = "emailQueue"
	const serializedMessage = `{"to":"toxanetoxa@gmail.com","subject":"Test Subject","body":"Test Body"}`
	expectedMessage := email.Message{
		To:      "toxanetoxa@gmail.com",
		Subject: "Test Subject",
		Body:    "Test Body",
	}

	db, mock := redismock.NewClientMock()
	mock.ExpectLPop(queueKey).SetVal(serializedMessage)

	q := &redis.Queue{
		Client: db,
		Ctx:    context.Background(),
		Key:    queueKey,
	}

	msg, err := q.Dequeue()
	assert.NoError(t, err)
	assert.Equal(t, expectedMessage, msg)
	assert.NoError(t, mock.ExpectationsWereMet())
}
