package messages

import "time"

type OrderCreatedMessage struct {
	ID                 int64              `json:"id"`
	CustomerID         int64              `json:"customer_id"`
	DeliveryAddress    string             `json:"delivery_address"`
	TotalPriceCents    int64              `json:"total_price_cents"`
	TotalPriceCurrency string             `json:"total_price_currency"`
	CreatedAt          time.Time          `json:"created_at"`
	OrderItems         []OrderItemMessage `json:"order_items"`
}

type OrderItemMessage struct {
	ID            int64  `json:"id"`
	OrderID       int64  `json:"order_id"`
	ProductID     int64  `json:"product_id"`
	ProductTitle  string `json:"product_title"`
	ProductURL    string `json:"product_url"`
	Quantity      int    `json:"quantity"`
	PriceCents    int64  `json:"price_cents"`
	PriceCurrency string `json:"price_currency"`
}
