package main

import (
	"log"
	"net/http"

	"github.com/Lamafout/online-store-api/internal/bll/services"
	"github.com/Lamafout/online-store-api/internal/config"
	"github.com/Lamafout/online-store-api/internal/dal/unit_of_work"
	"github.com/Lamafout/online-store-api/internal/handlers/v1"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found: %v", err)
	}

	cfg, err := config.LoadConfig("Development")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := sqlx.Connect("pgx", cfg.DbSettings.ConnectionString)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	uow := dal.NewUnitOfWork(db)
	orderService := services.NewOrderService(uow)

	r := chi.NewRouter()
	r.Route("/api/v1", func(r chi.Router) {
		r.Mount("/orders", v1.NewOrderHandler(orderService).Routes())
	})

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}