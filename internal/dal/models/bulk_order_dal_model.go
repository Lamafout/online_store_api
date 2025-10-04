package models

import "time"

type BulkOrderDalModel struct {
    CustomerID         int64     `db:"customer_id"`
    DeliveryAddress    string    `db:"delivery_address"`
    TotalPriceCents    int64     `db:"total_price_cents"`
    TotalPriceCurrency string    `db:"total_price_currency"`
    CreatedAt          time.Time `db:"created_at"`
    UpdatedAt          time.Time `db:"updated_at"`
}