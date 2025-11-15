package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Lamafout/online-store-api/internal/bll/services"
	"github.com/Lamafout/online-store-api/internal/bll/services/consumers"
	"github.com/Lamafout/online-store-api/internal/config"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found: %v", err)
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := sqlx.Connect("pgx", cfg.DbSettings.ConnectionString)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	rabbitURL := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		cfg.RabbitMqSettings.User, cfg.RabbitMqSettings.Password, cfg.RabbitMqSettings.Host, cfg.RabbitMqSettings.Port)

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
		true, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		log.Fatal("failed to declare queue:", err)
	}

	auditService := services.NewAuditService(db)

	consumer, err := consumers.NewOmsOrderCreatedConsumer(rabbitURL, cfg.RabbitMqSettings.OrderCreateQueue, auditService)
	if err != nil {
		log.Fatal("failed to start consumer:", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := consumer.Start(ctx); err != nil {
		log.Fatal("failed to start consumer:", err)
	}

	log.Printf("Audit consumer started. Listening to queue: %s", cfg.RabbitMqSettings.OrderCreateQueue)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	log.Println("Shutting down...")
	time.Sleep(1 * time.Second)
}