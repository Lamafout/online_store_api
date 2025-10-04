package config

import (
	"fmt"
	"os"
	"log"
	"github.com/joho/godotenv"
)

type DbSettings struct {
	ConnectionString           string
	MigrationConnectionString string
}

type Config struct {
	DbSettings DbSettings
	ServerPort string
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

	if user == "" || password == "" || dbName == "" || port == "" || host == "" || serverPort == "" {
		return nil, fmt.Errorf("missing required environment variables")
	}

	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbName)
	migrationConnString := connString
	return &Config{
		DbSettings: DbSettings{
			ConnectionString:           connString,
			MigrationConnectionString: migrationConnString,
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
