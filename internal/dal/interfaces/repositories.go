package interfaces

import (
	"context"

	"github.com/Lamafout/online-store-api/internal/dal/models"
)

type IOrderRepository interface {
	CreateOrder(ctx context.Context, order *models.V1OrderDal) error
	BulkInsertOrders(ctx context.Context, orders []models.QueryOrdersDalModel) error
	GetOrderByID(ctx context.Context, id int64) (*models.V1OrderDal, error)
}

type IOrderItemRepository interface {
	CreateOrderItem(ctx context.Context, item *models.V1OrderItemDal) error
	BulkInsertOrderItems(ctx context.Context, items []models.QueryOrderItemsDalModel) error
	GetOrderItemsByOrderID(ctx context.Context, orderID int64) ([]models.V1OrderItemDal, error)
}