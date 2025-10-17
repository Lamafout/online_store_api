package main

import (
	"log"
	"net/http"

	"github.com/Lamafout/online-store-api/internal/bll/services"
	"github.com/Lamafout/online-store-api/internal/config"
	v1 "github.com/Lamafout/online-store-api/internal/handlers/v1"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/jackc/pgx/v5/stdlib"

	_ "github.com/Lamafout/online-store-api/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Online Store API
// @version 1.0
// @description API for managing orders in an online store.
// @host localhost:8080
// @BasePath /api/v1
func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found: %v", err)
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := sqlx.Connect("pgx", cfg.DbSettings.ConnectionString)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	orderService := services.NewOrderService()

	r := chi.NewRouter()
	r.Route("/api/v1", func(r chi.Router) {
		r.Mount("/orders", v1.NewOrderHandler(db, orderService).Routes())
	})
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	log.Printf("Starting server on :%s", cfg.ServerPort)
	if err := http.ListenAndServe(":"+cfg.ServerPort, r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
