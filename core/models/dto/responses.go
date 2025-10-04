package dto

import (
	"github.com/Lamafout/online-store-api/core/models/common"
)

type V1CreateOrderResponse struct {
    Orders []common.Order `json:"orders"`
}

type V1QueryOrdersResponse struct {
    Orders []common.Order `json:"orders"`
}