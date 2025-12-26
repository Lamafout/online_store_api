package models

type QueryOrderItemsDalModel struct {
    IDs      []int64 `db:"ids"`
    OrderIDs []int64 `db:"order_ids"`
    Limit    int     `db:"limit"`
    Offset   int     `db:"offset"`
}