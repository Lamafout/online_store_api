package repositories

import (
	"context"
	"fmt"
	"strings"

	"github.com/Lamafout/online-store-api/internal/dal/interfaces"
	"github.com/Lamafout/online-store-api/internal/dal/models"
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
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id`
	var id int64
	err := r.db.QueryRowxContext(ctx, query, item.OrderID, item.ProductID, item.Quantity, item.ProductTitle, item.ProductURL, item.PriceCents, item.PriceCurrency, item.CreatedAt, item.UpdatedAt).Scan(&id)
	if err != nil {
		return fmt.Errorf("failed to create order item: %w", err)
	}
	item.ID = id
	return nil
}

// BulkInsertOrderItems inserts multiple order items
func (r *OrderItemRepository) BulkInsertOrderItems(ctx context.Context, items []models.QueryOrderItemsDalModel) error {
	if len(items) == 0 {
		return nil
	}

	// Build VALUES clause
	var values []interface{}
	var placeholders []string
	for i, item := range items {
		placeholders = append(placeholders, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)", i*9+1, i*9+2, i*9+3, i*9+4, i*9+5, i*9+6, i*9+7, i*9+8, i*9+9))
		values = append(values, item.OrderID, item.ProductID, item.Quantity, item.ProductTitle, item.ProductURL, item.PriceCents, item.PriceCurrency, item.CreatedAt, item.UpdatedAt)
	}

	query := fmt.Sprintf(`
		INSERT INTO order_items (order_id, product_id, quantity, product_title, product_url, price_cents, price_currency, created_at, updated_at)
		VALUES %s`, strings.Join(placeholders, ", "))
	_, err := r.db.ExecContext(ctx, query, values...)
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