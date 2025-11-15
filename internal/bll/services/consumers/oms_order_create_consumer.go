package consumers

import (
	"context"
	"encoding/json"
	"log"

	core "github.com/Lamafout/online-store-api/core/models/messages"
	"github.com/Lamafout/online-store-api/internal/bll/services"
	amqp "github.com/rabbitmq/amqp091-go"
)

type OmsOrderCreatedConsumer struct {
	conn         *amqp.Connection
	channel      *amqp.Channel
	queueName    string
	auditService *services.AuditService
}

func NewOmsOrderCreatedConsumer(amqpURL, queueName string, audit *services.AuditService) (*OmsOrderCreatedConsumer, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &OmsOrderCreatedConsumer{
		conn:         conn,
		channel:      ch,
		queueName:    queueName,
		auditService: audit,
	}, nil
}

func (c *OmsOrderCreatedConsumer) Start(ctx context.Context) error {
	msgs, err := c.channel.Consume(
		c.queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for d := range msgs {
			var msg core.OrderCreatedMessage
			if err := json.Unmarshal(d.Body, &msg); err != nil {
				log.Printf("failed to unmarshal message: %v", err)
				continue
			}

			if err := c.auditService.LogOrder(ctx, &msg); err != nil {
				log.Printf("failed to log audit: %v", err)
			} else {
				log.Printf("audit logged for order %d", msg.ID)
			}
		}
	}()

	return nil
}

func (c *OmsOrderCreatedConsumer) Stop() {
	c.channel.Close()
	c.conn.Close()
}
