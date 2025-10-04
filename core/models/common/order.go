package common

import "time"

type Order struct {
	ID                 int64       `json:"id"`
	CustomerID         int64       `json:"customer_id" validate:"required,gt=0"`
	DeliveryAddress    string      `json:"delivery_address" validate:"required,max=255"`
	TotalPriceCents    int64       `json:"total_price_cents" validate:"required,gt=0"`
	TotalPriceCurrency string      `json:"total_price_currency" validate:"required,oneof=USD EUR"`
	CreatedAt          time.Time   `json:"created_at"`
	UpdatedAt          time.Time   `json:"updated_at"`
	Items              []OrderItem `json:"items"`
}