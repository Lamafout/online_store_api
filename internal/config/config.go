package config

import (
	"fmt"
	"os"
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
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	host := os.Getenv("DB_HOST")
	serverPort := os.Getenv("SERVER_PORT")

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