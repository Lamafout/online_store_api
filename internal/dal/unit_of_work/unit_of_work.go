package dal

import (
	"context"
	"fmt"

	"github.com/Lamafout/online-store-api/internal/dal/interfaces"
	"github.com/Lamafout/online-store-api/internal/dal/repositories"
	"github.com/jmoiron/sqlx"
)

// UnitOfWork manages database transactions and repositories
type UnitOfWork struct {
	db            *sqlx.DB
	tx            *sqlx.Tx
	currentDB     interfaces.DBExecuter
	isTransaction bool
}

// NewUnitOfWork creates a new UnitOfWork
func NewUnitOfWork(db *sqlx.DB) *UnitOfWork {
	return &UnitOfWork{
		db:            db,
		currentDB:     db,
		isTransaction: false,
	}
}

// GetOrderRepo lazily initializes and returns the OrderRepository
func (u *UnitOfWork) GetOrderRepo() interfaces.IOrderRepository {
	return repositories.NewOrderRepository(u.currentDB)
}

// GetOrderItemRepo lazily initializes and returns the OrderItemRepository
func (u *UnitOfWork) GetOrderItemRepo() interfaces.IOrderItemRepository {
	return repositories.NewOrderItemRepository(u.currentDB)
}

func (u *UnitOfWork) GetAuditLogOrderRepo() interfaces.IAuditLogOrderRepository {
	return repositories.NewAuditLogOrderRepository(u.currentDB)
}

// Begin starts a new transaction
func (u *UnitOfWork) Begin(ctx context.Context) error {
	if u.isTransaction {
		return fmt.Errorf("transaction already started")
	}
	tx, err := u.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	u.tx = tx
	u.currentDB = tx
	u.isTransaction = true
	return nil
}

// Commit commits the transaction
func (u *UnitOfWork) Commit() error {
	if !u.isTransaction {
		return fmt.Errorf("no transaction to commit")
	}
	if err := u.tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	u.reset()
	return nil
}

// Rollback rolls back the transaction
func (u *UnitOfWork) Rollback() error {
	if !u.isTransaction {
		return fmt.Errorf("no transaction to rollback")
	}
	if err := u.tx.Rollback(); err != nil {
		return fmt.Errorf("failed to rollback transaction: %w", err)
	}
	u.reset()
	return nil
}

// reset resets the UnitOfWork to non-transactional state
func (u *UnitOfWork) reset() {
	u.tx = nil
	u.currentDB = u.db
	u.isTransaction = false
	// No need to reset repos, as they use u.currentDB
}