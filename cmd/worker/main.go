package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/suhas-developer07/EdwinNova-Server/internals/email"
	"github.com/suhas-developer07/EdwinNova-Server/internals/infrastructure/mail"
	"github.com/suhas-developer07/EdwinNova-Server/internals/infrastructure/rabbitmq"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}
	rabbit, err := rabbitmq.New(os.Getenv("RABBITMQ_URI"))
	if err != nil {
		log.Fatal(err)
	}
	defer rabbit.Close()

	queue := "email_queue"

	err = rabbit.DeclareQueue(queue, true)
	if err != nil {
		log.Fatal(err)
	}

	smtpClient, err := mail.NewSMTPClient()
	if err != nil {
		log.Fatal(err)
	}

	worker := email.NewWorker(rabbit, smtpClient, queue)

	err = worker.Start(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}
