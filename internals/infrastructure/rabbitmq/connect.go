package rabbitmq

import (
	"context"
	"fmt"

	"github.com/streadway/amqp"
)

type Connection struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func New(url string) (*Connection, error) {

	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to rabbitmq: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	fmt.Print("Connected to RabbitMQ")
	return &Connection{
		conn:    conn,
		channel: ch,
	}, nil
}

func (r *Connection) DeclareQueue(name string, durable bool) error {

	_, err := r.channel.QueueDeclare(
		name,
		durable,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	return nil
}

func (r *Connection) Publish(ctx context.Context, queue string, body []byte) error {

	err := r.channel.Publish(
		"",    
		queue,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		},
	)

	if err != nil {
		return fmt.Errorf("failed to publish: %w", err)
	}

	return nil
}

func (r *Connection) Consume(queue string, consumer string) (<-chan amqp.Delivery, error) {

	msgs, err := r.channel.Consume(
		queue,
		consumer,
		false, 
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to consume: %w", err)
	}

	return msgs, nil
}

func (r *Connection) Close() error {
	if err := r.channel.Close(); err != nil {
		return err
	}
	return r.conn.Close()
}