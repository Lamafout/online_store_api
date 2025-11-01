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
	ServerPort       string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found: %v", err)
	} else {
		log.Println("Successfully loaded .env file")
	}

	user := getEnv("DB_USER", "")
	password := getEnv("DB_PASSWORD", "")
	dbName := getEnv("DB_NAME", "postgres")
	port := getEnv("DB_PORT", "5432")
	host := getEnv("DB_HOST", "localhost")
	serverPort := getEnv("SERVER_PORT", "8080")
	rabbitHost := getEnv("RABBIT_HOST", "localhost")
	rabbitPort := getEnv("RABBIT_PORT", "5672")
	rabbitUser := getEnv("RABBIT_USER", "guest")
	rabbitPassword := getEnv("RABBIT_PASSWORD", "guest")
	rabbitQueue := getEnv("RABBIT_ORDER_CREATED_QUEUE", "oms.order.created")

	if user == "" || password == "" || dbName == "" || port == "" || host == "" || serverPort == "" {
		return nil, fmt.Errorf("missing required environment variables")
	}

	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbName)
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
		ServerPort: serverPort,
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
