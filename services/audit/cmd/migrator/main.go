package main

import (
	"fmt"
	"log"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {

	connString := "postgres://audit_user:audit_password@audit-db:5432/online_store_audit?sslmode=disable"
	
	db, err := sqlx.Connect("pgx", connString)
	if err != nil {
		log.Fatalf("Unable to connect to OMS database: %v", err)
	}
	defer db.Close()
	
	goose.SetDialect("postgres")
	
	if err := goose.Up(db.DB, "services/audit/migrations"); err != nil {
		log.Fatalf("Failed to run OMS migrations: %v", err)
	}
	
	fmt.Println("AUDIT migrations applied successfully")
}
