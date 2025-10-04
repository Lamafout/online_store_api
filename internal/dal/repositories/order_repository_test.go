package repositories

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/Lamafout/online-store-api/internal/config"
	"github.com/Lamafout/online-store-api/internal/dal/models"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func TestOrderRepository(t *testing.T) {
	// Load .env file
	if err := godotenv.Load("../../../.env"); err != nil {
		t.Fatalf("Failed to load .env file: %v", err)
	}

	environment := os.Getenv("APP_ENV")
	if environment == "" {
		t.Fatal("APP_ENV is not set in .env file")
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	db, err := sqlx.Connect("pgx", cfg.DbSettings.ConnectionString)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	repo := NewOrderRepository(db)
	ctx := context.Background()

	// Test CreateOrder
	order := &models.V1OrderDal{
		CustomerID:        1,
		DeliveryAddress:   "123 Test St",
		TotalPriceCents:   10000,
		TotalPriceCurrency: "USD",
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
	if err := repo.CreateOrder(ctx, order); err != nil {
		t.Fatalf("Failed to create order: %v", err)
	}
	t.Logf("Created order with ID: %d", order.ID)

	// Test GetOrderByID
	fetchedOrder, err := repo.GetOrderByID(ctx, order.ID)
	if err != nil {
		t.Fatalf("Failed to get order by ID: %v", err)
	}
	t.Logf("Fetched order: %+v", fetchedOrder)

	// Test BulkInsertOrders
	orders := []models.QueryOrdersDalModel{
		{
			CustomerID:        2,
			DeliveryAddress:   "456 Test St",
			TotalPriceCents:   20000,
			TotalPriceCurrency: "USD",
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		},
	}
	if err := repo.BulkInsertOrders(ctx, orders); err != nil {
		t.Fatalf("Failed to bulk insert orders: %v", err)
	}
	t.Log("Bulk inserted orders successfully")
}