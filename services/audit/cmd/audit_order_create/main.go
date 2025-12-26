package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	contracts "github.com/Lamafout/online-store-api/libs/contracts/messages"
	"github.com/Lamafout/online-store-api/services/audit/internal/bll/services"
	"github.com/Lamafout/online-store-api/services/audit/internal/config"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	 amqp"github.com/rabbitmq/amqp091-go"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found for AUDIT: %v", err)
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load AUDIT config: %v", err)
	}

	db, err := sqlx.Connect("pgx", cfg.DbSettings.ConnectionString)
	if err != nil {
		log.Fatalf("Failed to connect to AUDIT database: %v", err)
	}
	defer db.Close()

	rabbitURL := fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		cfg.RabbitMqSettings.User,
		cfg.RabbitMqSettings.Password,
		cfg.RabbitMqSettings.Host,
		cfg.RabbitMqSettings.Port,
	)

	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		log.Fatalf("failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	channel, err := conn.Channel()
	if err != nil {
		log.Fatalf("failed to open RabbitMQ channel: %v", err)
	}
	defer channel.Close()

	_, err = channel.QueueDeclare(
		cfg.RabbitMqSettings.OrderCreateQueue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("failed to declare queue: %v", err)
	}

	msgs, err := channel.Consume(
		cfg.RabbitMqSettings.OrderCreateQueue,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("failed to start consumer: %v", err)
	}

	auditService := services.NewAuditService(db)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		for d := range msgs {
			var msg contracts.OrderCreatedMessage
			if err := json.Unmarshal(d.Body, &msg); err != nil {
				log.Printf("failed to unmarshal message: %v", err)
				continue
			}
			if err := auditService.LogOrder(ctx, &msg); err != nil {
				log.Printf("failed to log audit: %v", err)
			} else {
				log.Printf("audit logged for order %d", msg.ID)
			}
		}
	}()

	log.Printf("Audit consumer started. Listening to queue: %s", cfg.RabbitMqSettings.OrderCreateQueue)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	log.Println("Shutting down...")
	time.Sleep(1 * time.Second)
}
