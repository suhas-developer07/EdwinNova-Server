package email

import (
	"fmt"

	"github.com/suhas-developer07/EdwinNova-Server/internals/infrastructure/rabbitmq"
)

func SetupEmailQueue(r *rabbitmq.Connection, queueName string) error {

	err := r.DeclareQueue(queueName, true)
	if err != nil {
		return fmt.Errorf("failed to setup email queue: %w", err)
	}

	return nil
}