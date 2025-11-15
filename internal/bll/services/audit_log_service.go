package services

import (
	"context"
	"fmt"

	"github.com/Lamafout/online-store-api/core/models/messages"
	"github.com/Lamafout/online-store-api/internal/dal/models"
	dal "github.com/Lamafout/online-store-api/internal/dal/unit_of_work"
	"github.com/jmoiron/sqlx"
)

type AuditService struct {
	db *sqlx.DB
}

func NewAuditService(db *sqlx.DB) *AuditService {
	return &AuditService{
		db: db,
	}
}

func (s *AuditService) LogOrder(ctx context.Context, msg *messages.OrderCreatedMessage) error {
	uow := dal.NewUnitOfWork(s.db)

	if err := uow.Begin(ctx); err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	defer uow.Rollback()
	for _, item := range msg.OrderItems {
		log := &models.V1AuditLogOrderDal{
			OrderID:     msg.ID,
			OrderItemID: item.ID,
			CustomerID:  msg.CustomerID,
			OrderStatus: "Created",
		}

		if err := uow.GetAuditLogOrderRepo().CreateAuditLog(ctx, log); err != nil {
			return fmt.Errorf("failed to create audit log: %w", err)
		}
	}

	if err := uow.Commit(); err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	return nil
}
