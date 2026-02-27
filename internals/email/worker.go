package email

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/suhas-developer07/EdwinNova-Server/internals/infrastructure/mail"
	"github.com/suhas-developer07/EdwinNova-Server/internals/infrastructure/rabbitmq"
)

type Worker struct {
	rabbit *rabbitmq.Connection
	smtp   *mail.SMTPClient
	queue  string
}

func NewWorker(r *rabbitmq.Connection, smtp *mail.SMTPClient, queue string) *Worker {
	return &Worker{
		rabbit: r,
		smtp:   smtp,
		queue:  queue,
	}
}

func (w *Worker) Start(ctx context.Context) error {

	msgs, err := w.rabbit.Consume(w.queue, "email_worker")
	if err != nil {
		return err
	}

	fmt.Println("Email worker started...")

	for msg := range msgs {

		var job EmailJob

		err := json.Unmarshal(msg.Body, &job)
		if err != nil {
			log.Println("failed to unmarshal:", err)
			msg.Nack(false, false) 
			continue
		}

		err = w.smtp.Send(job.To, job.Subject, "Your hackathon registration is successful")

		if err != nil {
			log.Println("failed to send email:", err)
			msg.Nack(false, true) 
			continue
		}

		msg.Ack(false)
	}

	return nil
}