package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type DbSettings struct {
	MigrationConnectionString string `mapstructure:"MigrationConnectionString"`
	ConnectionString         string `mapstructure:"ConnectionString"`
}

type Config struct {
	DbSettings DbSettings `mapstructure:"DbSettings"`
}

func LoadConfig(environment string) (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found: %v", err)
	} else {
		log.Println("Successfully loaded .env file")
	}

	user := getEnv("DB_USER", "")
	password := getEnv("DB_PASSWORD", "")
	dbName := getEnv("DB_NAME", "postgres")

	log.Printf("Database config: DB_USER=%s, DB_PASSWORD=%s, DB_NAME=%s", 
		user, maskPassword(password), dbName)

	cfg := &Config{
		DbSettings: DbSettings{
			MigrationConnectionString: buildConnectionString(user, password, dbName, 5432),
			ConnectionString:         buildConnectionString(user, password, dbName, 15432),
		},
	}

	viper.SetConfigName("appsettings." + environment)
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	
	if err := viper.ReadInConfig(); err == nil {
		log.Printf("Loaded appsettings.%s.json", environment)
		
		var fileConfig Config
		if err := viper.Unmarshal(&fileConfig); err != nil {
			return nil, fmt.Errorf("failed to unmarshal config file: %w", err)
		}
		
		if fileConfig.DbSettings.MigrationConnectionString != "" {
			cfg.DbSettings.MigrationConnectionString = fileConfig.DbSettings.MigrationConnectionString
		}
		if fileConfig.DbSettings.ConnectionString != "" {
			cfg.DbSettings.ConnectionString = fileConfig.DbSettings.ConnectionString
		}
	} else {
		log.Printf("No appsettings.%s.json found, using environment variables only", environment)
	}

	if cfg.DbSettings.MigrationConnectionString == "" {
		return nil, fmt.Errorf("MigrationConnectionString is empty")
	}
	if cfg.DbSettings.ConnectionString == "" {
		return nil, fmt.Errorf("ConnectionString is empty")
	}

	log.Printf("Final connection strings: Migration=%s, Main=%s", 
		maskConnectionString(cfg.DbSettings.MigrationConnectionString),
		maskConnectionString(cfg.DbSettings.ConnectionString))

	return cfg, nil
}

func buildConnectionString(user, password, dbName string, port int) string {
	if user == "" || password == "" || dbName == "" {
		return ""
	}
	return fmt.Sprintf("postgres://%s:%s@127.0.0.1:%d/%s?sslmode=disable", 
		user, password, port, dbName)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func maskPassword(password string) string {
	if password == "" {
		return "<empty>"
	}
	return "***"
}

func maskConnectionString(connStr string) string {
	if connStr == "" {
		return "<empty>"
	}

	if len(connStr) > 30 {
		return connStr[:30] + "***"
	}
	return "***"
}