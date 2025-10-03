package main

import (
	_"context"
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
	_"github.com/jackc/pgx/v5/stdlib"
	"github.com/Lamafout/online-store-api/internal/config"
)

func main() {
	environment := os.Getenv("APP_ENV")
	if environment == "" {
		log.Fatal("APP_ENV is not set")
	}

	cfg, err := config.LoadConfig(environment)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	connString := cfg.DbSettings.MigrationConnectionString
	if connString == "" {
		log.Fatal("MigrationConnectionString is not set")
	}

	db, err := sqlx.Connect("pgx", connString)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer db.Close()

	goose.SetDialect("postgres")
	if err := goose.Up(db.DB, "migrations"); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	fmt.Println("Migrations applied successfully")
}