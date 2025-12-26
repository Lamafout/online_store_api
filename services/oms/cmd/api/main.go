package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Lamafout/online-store-api/services/oms/internal/bll/services"
	"github.com/Lamafout/online-store-api/services/oms/internal/config"
	v1 "github.com/Lamafout/online-store-api/services/oms/internal/handlers/v1"

	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"

	_ "github.com/Lamafout/online-store-api/services/oms/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Online Store OMS API
// @version 1.0
// @description API for managing orders in an online store.
// @host localhost:8080
// @BasePath /api/v1
func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found for OMS: %v", err)
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load OMS config: %v", err)
	}

	db, err := sqlx.Connect("pgx", cfg.DbSettings.ConnectionString)
	if err != nil {
		log.Fatalf("Failed to connect to OMS database: %v", err)
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
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,
	)
	if err != nil {
		log.Fatalf("failed to declare queue: %v", err)
	}

	publisher := services.NewRabbitPublisher(channel)
	orderService := services.NewOrderService(*publisher, cfg.RabbitMqSettings, db)

	r := chi.NewRouter()

	r.Get("/swagger/*", httpSwagger.WrapHandler)
	r.Route("/api/v1", func(r chi.Router) {
		r.Mount("/orders", v1.NewOrderHandler(db, orderService).Routes())
	})

	log.Printf("Starting OMS server on %s", cfg.ServerPort)
	if err := http.ListenAndServe(cfg.ServerPort, r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
