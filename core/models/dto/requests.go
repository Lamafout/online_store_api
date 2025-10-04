package dto

type V1CreateOrderRequest struct {
	Orders []V1CreateOrder `json:"orders" validate:"required,dive"`
}

type V1CreateOrder struct {
	CustomerID         int64               `json:"customer_id" validate:"required,gt=0"`
	DeliveryAddress    string              `json:"delivery_address" validate:"required,max=255"`
	TotalPriceCents    int64               `json:"total_price_cents" validate:"required,gt=0"`
	TotalPriceCurrency string              `json:"total_price_currency" validate:"required,oneof=USD EUR"`
	OrderItems         []V1CreateOrderItem `json:"order_items" validate:"required,dive"`
}

type V1CreateOrderItem struct {
	ProductID     int64  `json:"product_id" validate:"required,gt=0"`
	Quantity      int    `json:"quantity" validate:"required,gt=0"`
	ProductTitle  string `json:"product_title" validate:"required,max=255"`
	ProductURL    string `json:"product_url" validate:"required,url"`
	PriceCents    int64  `json:"price_cents" validate:"required,gt=0"`
	PriceCurrency string `json:"price_currency" validate:"required,oneof=USD EUR"`
}

type V1QueryOrdersRequest struct {
	IDs               []int64 `json:"ids"`
	CustomerIDs       []int64 `json:"customer_ids"`
	Page              *int    `json:"page"`
	PageSize          *int    `json:"page_size"`
	IncludeOrderItems bool    `json:"include_order_items"`
}
