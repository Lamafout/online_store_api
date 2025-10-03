package models

import (
	"time"
)

type QueryOrderItemsDalModel struct {
	ID             int64     `db:"id"`
	OrderID        int64     `db:"order_id"`
	ProductID      int64     `db:"product_id"`
	Quantity       int       `db:"quantity"`
	ProductTitle   string    `db:"product_title"`
	ProductURL     string    `db:"product_url"`
	PriceCents     int64     `db:"price_cents"`
	PriceCurrency  string    `db:"price_currency"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}