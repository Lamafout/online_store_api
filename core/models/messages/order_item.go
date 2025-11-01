package messages

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