package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type DbSettings struct {
	ConnectionString          string
	MigrationConnectionString string
}

type RabbitMqSettings struct {
	Host             string
	Port             string
	User             string
	Password         string
	OrderCreateQueue string
}

type Config struct {
	DbSettings       DbSettings
	RabbitMqSettings RabbitMqSettings
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found for AUDIT: %v", err)
	} else {
		log.Println("Successfully loaded .env file for AUDIT")
	}

	user := getEnv("AUDIT_DB_USER", "")
	password := getEnv("AUDIT_DB_PASSWORD", "")
	dbName := getEnv("AUDIT_DB_NAME", "online_store_audit")
	port := getEnv("AUDIT_DB_PORT", "5432")
	host := getEnv("AUDIT_DB_HOST", "localhost")

	rabbitHost := getEnv("RABBIT_HOST", "localhost")
	rabbitPort := getEnv("RABBIT_PORT", "5672")
	rabbitUser := getEnv("RABBIT_USER", "guest")
	rabbitPassword := getEnv("RABBIT_PASSWORD", "guest")
	rabbitQueue := getEnv("RABBIT_ORDER_CREATED_QUEUE", "oms.order.created")

	if user == "" || password == "" || dbName == "" || port == "" || host == "" {
		return nil, fmt.Errorf("missing required environment variables for AUDIT")
	}

	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user, password, host, port, dbName,
	)
	migrationConnString := connString

	return &Config{
		DbSettings: DbSettings{
			ConnectionString:          connString,
			MigrationConnectionString: migrationConnString,
		},
		RabbitMqSettings: RabbitMqSettings{
			Host:             rabbitHost,
			Port:             rabbitPort,
			User:             rabbitUser,
			Password:         rabbitPassword,
			OrderCreateQueue: rabbitQueue,
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
