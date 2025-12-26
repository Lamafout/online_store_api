package main

import (
	"fmt"
	"log"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	connString := "postgres://oms_user:oms_password@oms-db:5432/online_store_oms?sslmode=disable"
	
	db, err := sqlx.Connect("pgx", connString)
	if err != nil {
		log.Fatalf("Unable to connect to OMS database: %v", err)
	}
	defer db.Close()

	goose.SetDialect("postgres")
	if err := goose.Up(db.DB, "services/oms/migrations"); err != nil {
		log.Fatalf("Failed to run OMS migrations: %v", err)
	}

	fmt.Println("OMS migrations applied successfully")
}
