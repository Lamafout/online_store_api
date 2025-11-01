package repositories

import (
	"context"
	"time"

	"github.com/Lamafout/online-store-api/internal/dal/interfaces"
	"github.com/Lamafout/online-store-api/internal/dal/models"
)

type AuditLogOrderRepository struct {
	db interfaces.DBExecuter
}

func NewAuditLogOrderRepository(db interfaces.DBExecuter) *AuditLogOrderRepository {
	return &AuditLogOrderRepository{db: db}
}

func (r *AuditLogOrderRepository) CreateAuditLog(ctx context.Context, log *models.V1AuditLogOrderDal) error {
	query := `
		INSERT INTO audit_log_order 
			(order_id, order_item_id, customer_id, order_status, created_at, updated_at)
		VALUES
			($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	if log.CreatedAt.IsZero() {
		log.CreatedAt = time.Now()
	}
	if log.UpdatedAt.IsZero() {
		log.UpdatedAt = time.Now()
	}

	return r.db.GetContext(ctx, &log.ID, query,
		log.OrderID,
		log.OrderItemID,
		log.CustomerID,
		log.OrderStatus,
		log.CreatedAt,
		log.UpdatedAt,
	)
}
