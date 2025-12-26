package unit_of_work

import (
	"context"
	"fmt"

	"github.com/Lamafout/online-store-api/services/oms/internal/dal/interfaces"
	"github.com/Lamafout/online-store-api/services/oms/internal/dal/repositories"
	"github.com/jmoiron/sqlx"
)

// UnitOfWork manages database transactions and repositories.
type UnitOfWork struct {
	db            *sqlx.DB
	tx            *sqlx.Tx
	currentDB     interfaces.DBExecuter
	isTransaction bool
}

// NewUnitOfWork creates a new UnitOfWork.
func NewUnitOfWork(db *sqlx.DB) *UnitOfWork {
	return &UnitOfWork{
		db:        db,
		currentDB: db,
	}
}

func (u *UnitOfWork) GetOrderRepo() interfaces.IOrderRepository {
	return repositories.NewOrderRepository(u.currentDB)
}

func (u *UnitOfWork) GetOrderItemRepo() interfaces.IOrderItemRepository {
	return repositories.NewOrderItemRepository(u.currentDB)
}

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

func (u *UnitOfWork) reset() {
	u.tx = nil
	u.currentDB = u.db
	u.isTransaction = false
}
