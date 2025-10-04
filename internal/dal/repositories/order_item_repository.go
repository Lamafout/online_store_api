package repositories

import (
	"context"
	"fmt"

	"github.com/Lamafout/online-store-api/internal/dal/models"
	"github.com/Lamafout/online-store-api/internal/dal/interfaces"
)

// OrderItemRepository handles database operations for order items
type OrderItemRepository struct {
	db interfaces.DBExecuter
}

// NewOrderItemRepository creates a new OrderItemRepository
func NewOrderItemRepository(db interfaces.DBExecuter) *OrderItemRepository {
	return &OrderItemRepository{db: db}
}

// CreateOrderItem creates a single order item
func (r *OrderItemRepository) CreateOrderItem(ctx context.Context, item *models.V1OrderItemDal) error {
	query := `
		INSERT INTO order_items (order_id, product_id, quantity, product_title, product_url, price_cents, price_currency, created_at, updated_at)
		VALUES (:order_id, :product_id, :quantity, :product_title, :product_url, :price_cents, :price_currency, :created_at, :updated_at)
		RETURNING id`
	var id int64
	err := r.db.QueryRowxContext(ctx, query, item).Scan(&id)
	if err != nil {
		return fmt.Errorf("failed to create order item: %w", err)
	}
	item.ID = id
	return nil
}

// BulkInsertOrderItems inserts multiple order items using the v1_order_item composite type
func (r *OrderItemRepository) BulkInsertOrderItems(ctx context.Context, items []models.QueryOrderItemsDalModel) error {
	query := `
		INSERT INTO order_items (order_id, product_id, quantity, product_title, product_url, price_cents, price_currency, created_at, updated_at)
		SELECT * FROM UNNEST($1::v1_order_item[])`
	_, err := r.db.ExecContext(ctx, query, items)
	if err != nil {
		return fmt.Errorf("failed to bulk insert order items: %w", err)
	}
	return nil
}

// GetOrderItemsByOrderID retrieves all order items for a given order ID
func (r *OrderItemRepository) GetOrderItemsByOrderID(ctx context.Context, orderID int64) ([]models.V1OrderItemDal, error) {
	query := `SELECT * FROM order_items WHERE order_id = $1`
	var items []models.V1OrderItemDal
	err := r.db.SelectContext(ctx, &items, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order items for order ID %d: %w", orderID, err)
	}
	return items, nil
}