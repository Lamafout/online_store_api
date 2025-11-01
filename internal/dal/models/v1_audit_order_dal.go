package models

import "time"

type V1AuditLogOrderDal struct {
	ID          int64     `db:"id"`
	OrderID     int64     `db:"order_id"`
	OrderItemID int64     `db:"order_item_id"`
	CustomerID  int64     `db:"customer_id"`
	OrderStatus string    `db:"order_status"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
