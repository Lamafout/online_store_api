package services

import (
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/Lamafout/online-store-api/internal/config"
)

type RabbitPublisher struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   string
}

func NewRabbitPublisher(cfg *config.RabbitMqSettings) (*RabbitPublisher, error) {
	url := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		cfg.User, cfg.Password, cfg.Host, cfg.Port)

	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	_, err = ch.QueueDeclare(
		cfg.Queue, // name
		true,      // durable
		false,     // auto-delete
		false,     // exclusive
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	log.Printf("Connected to RabbitMQ queue: %s", cfg.Queue)
	return &RabbitPublisher{conn: conn, channel: ch, queue: cfg.Queue}, nil
}

func (p *RabbitPublisher) Publish(message any) error {
	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	err = p.channel.Publish(
		"",        // exchange
		p.queue,   // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	log.Printf("Message published to queue %s: %s", p.queue, string(body))
	return nil
}

func (p *RabbitPublisher) Close() {
	if err := p.channel.Close(); err != nil {
		log.Printf("Error closing channel: %v", err)
	}
	if err := p.conn.Close(); err != nil {
		log.Printf("Error closing connection: %v", err)
	}
}
