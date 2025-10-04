package repositories

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/Lamafout/online-store-api/internal/config"
	"github.com/Lamafout/online-store-api/internal/dal/interfaces"
	"github.com/Lamafout/online-store-api/internal/dal/models"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func TestOrderItemRepository(t *testing.T) {
	// Load .env file
	if err := godotenv.Load("../../../.env"); err != nil {
		t.Fatalf("Failed to load .env file: %v", err)
	}

	environment := os.Getenv("APP_ENV")
	if environment == "" {
		t.Fatal("APP_ENV is not set in .env file")
	}

	cfg, err := config.LoadConfig(environment)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	db, err := sqlx.Connect("pgx", cfg.DbSettings.ConnectionString)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Use db as interfaces.DBExecuter
	var dbExecuter interfaces.DBExecuter = db

	// Create an order first
	orderRepo := NewOrderRepository(dbExecuter)
	ctx := context.Background()
	order := &models.V1OrderDal{
		CustomerID:        1,
		DeliveryAddress:   "123 Test St",
		TotalPriceCents:   10000,
		TotalPriceCurrency: "USD",
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
	if err := orderRepo.CreateOrder(ctx, order); err != nil {
		t.Fatalf("Failed to create order: %v", err)
	}

	repo := NewOrderItemRepository(dbExecuter)

	// Test CreateOrderItem
	item := &models.V1OrderItemDal{
		OrderID:       order.ID,
		ProductID:     1,
		Quantity:      2,
		ProductTitle:  "Test Product",
		ProductURL:    "http://example.com/product",
		PriceCents:    5000,
		PriceCurrency: "USD",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	if err := repo.CreateOrderItem(ctx, item); err != nil {
		t.Fatalf("Failed to create order item: %v", err)
	}
	t.Logf("Created order item with ID: %d", item.ID)

	// Test GetOrderItemsByOrderID
	items, err := repo.GetOrderItemsByOrderID(ctx, order.ID)
	if err != nil {
		t.Fatalf("Failed to get order items: %v", err)
	}
	t.Logf("Fetched order items: %+v", items)

	// Test BulkInsertOrderItems
	orderItems := []models.QueryOrderItemsDalModel{
		{
			OrderID:       order.ID,
			ProductID:     2,
			Quantity:      3,
			ProductTitle:  "Another Product",
			ProductURL:    "http://example.com/another",
			PriceCents:    7500,
			PriceCurrency: "USD",
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
	}
	if err := repo.BulkInsertOrderItems(ctx, orderItems); err != nil {
		t.Fatalf("Failed to bulk insert order items: %v", err)
	}
	t.Log("Bulk inserted order items successfully")
}