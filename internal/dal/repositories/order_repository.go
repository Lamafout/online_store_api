package repositories

import (
	"context"
	"fmt"

	"github.com/Lamafout/online-store-api/internal/dal/models"
	"github.com/Lamafout/online-store-api/internal/dal/interfaces"
)

// OrderRepository handles database operations for orders
type OrderRepository struct {
	db interfaces.DBExecuter
}

// NewOrderRepository creates a new OrderRepository
func NewOrderRepository(db interfaces.DBExecuter) *OrderRepository {
	return &OrderRepository{db: db}
}

// CreateOrder creates a single order
func (r *OrderRepository) CreateOrder(ctx context.Context, order *models.V1OrderDal) error {
	query := `
		INSERT INTO orders (customer_id, delivery_address, total_price_cents, total_price_currency, created_at, updated_at)
		VALUES (:customer_id, :delivery_address, :total_price_cents, :total_price_currency, :created_at, :updated_at)
		RETURNING id`
	var id int64
	err := r.db.QueryRowxContext(ctx, query, order).Scan(&id)
	if err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}
	order.ID = id
	return nil
}

// BulkInsertOrders inserts multiple orders using the v1_order composite type
func (r *OrderRepository) BulkInsertOrders(ctx context.Context, orders []models.QueryOrdersDalModel) error {
	query := `
		INSERT INTO orders (customer_id, delivery_address, total_price_cents, total_price_currency, created_at, updated_at)
		SELECT * FROM UNNEST($1::v1_order[])`
	_, err := r.db.ExecContext(ctx, query, orders)
	if err != nil {
		return fmt.Errorf("failed to bulk insert orders: %w", err)
	}
	return nil
}

// GetOrderByID retrieves an order by its ID
func (r *OrderRepository) GetOrderByID(ctx context.Context, id int64) (*models.V1OrderDal, error) {
	query := `SELECT * FROM orders WHERE id = $1`
	var order models.V1OrderDal
	err := r.db.GetContext(ctx, &order, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get order by ID %d: %w", id, err)
	}
	return &order, nil
}