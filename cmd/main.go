package main

import (
	"email-sendler/internal/email"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

func main() {
	err := godotenv.Load("../.env.prod")
	if err != nil {
		log.Fatalf("Error loading .env.prod file: %v", err)
	}

	Host := os.Getenv("HOST")
	Username := os.Getenv("USERNAME")
	Password := os.Getenv("PASSWORD")
	Port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatalf("PORT must be an integer")
	}

	sender := &email.SenderConf{
		Host:     Host,
		Port:     Port,
		Username: Username,
		Password: Password,
	}

	// Проверка отправки сообщения (РАБОТАЕТ!!!)
	//err = sender.Send("toxanetoxa@gmail.com", "test 1", "Body 1")
	//if err != nil {
	//	log.Fatalf("Error sending email: %v", err)
	//	return
	//}

	// TODO нужно реализовать метод создания очереди и добаления в нее сообщений
	// TODO нужно реализовать метод отправки сообщений из очереди
	// TODO нужно реализовать метод отправки сообщений из очереди с ограничением по количеству сообщений в единицу времени
}
