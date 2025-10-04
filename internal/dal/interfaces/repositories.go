package interfaces

import (
	"context"

	"github.com/Lamafout/online-store-api/internal/dal/models"
)

type IOrderRepository interface {
	CreateOrder(ctx context.Context, order *models.V1OrderDal) error
	BulkInsertOrders(ctx context.Context, orders []models.BulkOrderDalModel) ([]models.V1OrderDal, error)
	GetOrderByID(ctx context.Context, id int64) (*models.V1OrderDal, error)
	QueryOrders(ctx context.Context, req *models.QueryOrdersDalModel) ([]models.V1OrderDal, error)
}

type IOrderItemRepository interface {
	CreateOrderItem(ctx context.Context, item *models.V1OrderItemDal) error
	BulkInsertOrderItems(ctx context.Context, items []models.BulkOrderItemDalModel) ([]models.V1OrderItemDal, error)
	GetOrderItemsByOrderID(ctx context.Context, orderID int64) ([]models.V1OrderItemDal, error)
	QueryOrderItems(ctx context.Context, req *models.QueryOrderItemsDalModel) ([]models.V1OrderItemDal, error)
}