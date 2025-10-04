package repositories

import (
	"context"
	"fmt"
	"strings"

	"github.com/Lamafout/online-store-api/internal/dal/interfaces"
	"github.com/Lamafout/online-store-api/internal/dal/models"
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
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`
	var id int64
	err := r.db.QueryRowxContext(ctx, query, order.CustomerID, order.DeliveryAddress, order.TotalPriceCents, order.TotalPriceCurrency, order.CreatedAt, order.UpdatedAt).Scan(&id)
	if err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}
	order.ID = id
	return nil
}

// BulkInsertOrders inserts multiple orders
func (r *OrderRepository) BulkInsertOrders(ctx context.Context, orders []models.BulkOrderDalModel) ([]models.V1OrderDal, error) {
    if len(orders) == 0 {
        return []models.V1OrderDal{}, nil
    }

    // Build VALUES clause
    var values []interface{}
    var placeholders []string
    for i, order := range orders {
        placeholders = append(placeholders, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d)", i*6+1, i*6+2, i*6+3, i*6+4, i*6+5, i*6+6))
        values = append(values, order.CustomerID, order.DeliveryAddress, order.TotalPriceCents, order.TotalPriceCurrency, order.CreatedAt, order.UpdatedAt)
    }

    query := fmt.Sprintf(`
        INSERT INTO orders (customer_id, delivery_address, total_price_cents, total_price_currency, created_at, updated_at)
        VALUES %s 
        RETURNING id, customer_id, delivery_address, total_price_cents, total_price_currency, created_at, updated_at`, 
        strings.Join(placeholders, ", "))
    
    var insertedOrders []models.V1OrderDal
    err := r.db.SelectContext(ctx, &insertedOrders, query, values...)
    if err != nil {
        return nil, fmt.Errorf("failed to bulk insert orders: %w", err)
    }
    
    return insertedOrders, nil
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

func (r *OrderRepository) QueryOrders(ctx context.Context, req *models.QueryOrdersDalModel) ([]models.V1OrderDal, error) {
    query := `SELECT * FROM orders WHERE 1=1`
    var args []interface{}
    var conditions []string

    if len(req.IDs) > 0 {
        conditions = append(conditions, fmt.Sprintf("id = ANY($%d)", len(args)+1))
        args = append(args, req.IDs)
    }
    
    if len(req.CustomerIDs) > 0 {
        conditions = append(conditions, fmt.Sprintf("customer_id = ANY($%d)", len(args)+1))
        args = append(args, req.CustomerIDs)
    }

    if len(conditions) > 0 {
        query += " AND " + strings.Join(conditions, " AND ")
    }

    if req.Limit > 0 {
        query += fmt.Sprintf(" LIMIT $%d", len(args)+1)
        args = append(args, req.Limit)
    }

    if req.Offset > 0 {
        query += fmt.Sprintf(" OFFSET $%d", len(args)+1)
        args = append(args, req.Offset)
    }

    var orders []models.V1OrderDal
    err := r.db.SelectContext(ctx, &orders, query, args...)
    if err != nil {
        return nil, fmt.Errorf("failed to query orders: %w", err)
    }
    return orders, nil
}