package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/smtp"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../.env.prod")
	if err != nil {
		log.Fatalf("Error loading .env.prod file: %v", err)
	}

	smtpHost := os.Getenv("HOST")
	smtpPort, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatalf("PORT must be an integer")
	}
	smtpUsername := os.Getenv("USERNAME")
	smtpPassword := os.Getenv("PASSWORD")

	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpHost)

	// Подключение к SMTP-серверу
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", smtpHost, smtpPort))
	if err != nil {
		log.Fatalf("Failed to connect to the SMTP server: %v", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		log.Fatalf("Failed to create SMTP client: %v", err)
	}
	defer client.Quit()

	// Установление TLS-соединения
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true, // Установите false, если хотите проверять сертификаты
		ServerName:         smtpHost,
	}

	if err = client.StartTLS(tlsConfig); err != nil {
		log.Fatalf("Failed to start TLS: %v", err)
	}

	// Авторизация
	if err = client.Auth(auth); err != nil {
		log.Fatalf("Failed to authenticate: %v", err)
	}

	log.Println("Successfully connected and authenticated to the SMTP server")

	// Отправка тестового email
	to := []string{"toxanetoxa@gmail.com"}
	msg := []byte("To: toxanetoxa@gmail.com\r\n" +
		"Subject: Test email from Go\r\n" +
		"\r\n" +
		"This is the email body.\r\n")

	if err = client.Mail(smtpUsername); err != nil {
		log.Fatalf("Failed to set sender: %v", err)
	}
	for _, addr := range to {
		if err = client.Rcpt(addr); err != nil {
			log.Fatalf("Failed to set recipient: %v", err)
		}
	}

	wc, err := client.Data()
	if err != nil {
		log.Fatalf("Failed to send data command: %v", err)
	}
	defer wc.Close()

	_, err = wc.Write(msg)
	if err != nil {
		log.Fatalf("Failed to write message: %v", err)
	}

	log.Println("Email sent successfully")
}
