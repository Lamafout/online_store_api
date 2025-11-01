package services

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitPublisher struct {
	channel *amqp.Channel
}

func NewRabbitPublisher(channel *amqp.Channel) *RabbitPublisher {
	return &RabbitPublisher{channel: channel}
}

func (p *RabbitPublisher) Publish(ctx context.Context, messages []any, queue string) error {
	for _, msg := range messages {
		body, err := json.Marshal(msg)
		if err != nil {
			return fmt.Errorf("marshal failed: %w", err)
		}
		if err := p.channel.PublishWithContext(ctx,
			"",
			queue,
			false,
			false,
			amqp.Publishing{
				ContentType: "application/json",
				Body:        body,
			},
		); err != nil {
			return fmt.Errorf("publish failed: %w", err)
		}
	}
	return nil
}
