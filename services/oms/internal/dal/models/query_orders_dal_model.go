package models

type QueryOrdersDalModel struct {
    IDs         []int64 `db:"ids"`
    CustomerIDs []int64 `db:"customer_ids"`
    Limit       int     `db:"limit"`
    Offset      int     `db:"offset"`
}