package common

import "time"

type OrderItem struct {
	ID            int64     `json:"id"`
	OrderID       int64     `json:"order_id"`
	ProductID     int64     `json:"product_id" validate:"required,gte=0"`
	Quantity      int       `json:"quantity" validate:"required,gte=0"`
	ProductTitle  string    `json:"product_title" validate:"required,max=255"`
	ProductURL    string    `json:"product_url" validate:"required,url"`
	PriceCents    int64     `json:"price_cents" validate:"required,gte=0"`
	PriceCurrency string    `json:"price_currency" validate:"required,oneof=USD EUR"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
