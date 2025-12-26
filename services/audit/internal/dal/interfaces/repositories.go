package interfaces

import (
	"context"

	"github.com/Lamafout/online-store-api/services/audit/internal/dal/models"
)

type IAuditLogOrderRepository interface {
	CreateAuditLog(ctx context.Context, log *models.V1AuditLogOrderDal) error
}
