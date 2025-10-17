package services

import (
	"context"
	"fmt"
	"time"

	core "github.com/Lamafout/online-store-api/core/models/common"
	"github.com/Lamafout/online-store-api/core/models/dto"
	"github.com/Lamafout/online-store-api/internal/dal/models"
	"github.com/Lamafout/online-store-api/internal/dal/unit_of_work"
	"github.com/go-playground/validator/v10"
)

type OrderService struct {
	validate *validator.Validate
}

func NewOrderService() *OrderService {
	return &OrderService{
		validate: validator.New(),
	}
}

func (s *OrderService) CreateOrder(
	ctx context.Context,
	uow *dal.UnitOfWork,
	order *core.Order,
) error {
	if err := s.validate.Struct(order); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	dalOrder := &models.V1OrderDal{
		CustomerID:         order.CustomerID,
		DeliveryAddress:    order.DeliveryAddress,
		TotalPriceCents:    order.TotalPriceCents,
		TotalPriceCurrency: order.TotalPriceCurrency,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	if err := uow.GetOrderRepo().CreateOrder(ctx, dalOrder); err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}
	order.ID = dalOrder.ID

	for i := range order.Items {
		item := &order.Items[i]
		if err := s.validate.Struct(item); err != nil {
			return fmt.Errorf("validation failed for item: %w", err)
		}

		dalItem := &models.V1OrderItemDal{
			OrderID:       dalOrder.ID,
			ProductID:     item.ProductID,
			Quantity:      item.Quantity,
			ProductTitle:  item.ProductTitle,
			ProductURL:    item.ProductURL,
			PriceCents:    item.PriceCents,
			PriceCurrency: item.PriceCurrency,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		if err := uow.GetOrderItemRepo().CreateOrderItem(ctx, dalItem); err != nil {
			return fmt.Errorf("failed to create order item: %w", err)
		}

		item.ID = dalItem.ID
		item.OrderID = dalItem.OrderID
	}

	return nil
}

func (s *OrderService) GetOrder(
	ctx context.Context,
	uow *dal.UnitOfWork,
	id int64,
) (*core.Order, error) {
	dalOrder, err := uow.GetOrderRepo().GetOrderByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	dalItems, err := uow.GetOrderItemRepo().GetOrderItemsByOrderID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get order items: %w", err)
	}

	order := &core.Order{
		ID:                 dalOrder.ID,
		CustomerID:         dalOrder.CustomerID,
		DeliveryAddress:    dalOrder.DeliveryAddress,
		TotalPriceCents:    dalOrder.TotalPriceCents,
		TotalPriceCurrency: dalOrder.TotalPriceCurrency,
		CreatedAt:          dalOrder.CreatedAt,
		UpdatedAt:          dalOrder.UpdatedAt,
		Items:              make([]core.OrderItem, len(dalItems)),
	}

	for i, item := range dalItems {
		order.Items[i] = core.OrderItem{
			ID:            item.ID,
			OrderID:       item.OrderID,
			ProductID:     item.ProductID,
			Quantity:      item.Quantity,
			ProductTitle:  item.ProductTitle,
			ProductURL:    item.ProductURL,
			PriceCents:    item.PriceCents,
			PriceCurrency: item.PriceCurrency,
			CreatedAt:     item.CreatedAt,
			UpdatedAt:     item.UpdatedAt,
		}
	}

	return order, nil
}

func (s *OrderService) BatchCreateOrders(
	ctx context.Context,
	uow *dal.UnitOfWork,
	orders []*core.Order,
) ([]*core.Order, error) {
	if len(orders) == 0 {
		return nil, fmt.Errorf("orders list cannot be empty")
	}

	for _, order := range orders {
		if err := s.validate.Struct(order); err != nil {
			return nil, fmt.Errorf("validation failed: %w", err)
		}
		total := int64(0)
		for _, item := range order.Items {
			total += item.PriceCents * int64(item.Quantity)
		}
		if total != order.TotalPriceCents {
			return nil, fmt.Errorf("total price mismatch for order: expected %d, got %d", total, order.TotalPriceCents)
		}
	}

	now := time.Now()
	bulkOrders := make([]models.BulkOrderDalModel, len(orders))
	for i, order := range orders {
		bulkOrders[i] = models.BulkOrderDalModel{
			CustomerID:         order.CustomerID,
			DeliveryAddress:    order.DeliveryAddress,
			TotalPriceCents:    order.TotalPriceCents,
			TotalPriceCurrency: order.TotalPriceCurrency,
			CreatedAt:          now,
			UpdatedAt:          now,
		}
	}

	insertedOrders, err := uow.GetOrderRepo().BulkInsertOrders(ctx, bulkOrders)
	if err != nil {
		return nil, fmt.Errorf("failed to bulk insert orders: %w", err)
	}

	for i := range insertedOrders {
		orders[i].ID = insertedOrders[i].ID
		orders[i].CreatedAt = insertedOrders[i].CreatedAt
		orders[i].UpdatedAt = insertedOrders[i].UpdatedAt

		for j := range orders[i].Items {
			item := &orders[i].Items[j]
			dalItem := &models.V1OrderItemDal{
				OrderID:       insertedOrders[i].ID,
				ProductID:     item.ProductID,
				Quantity:      item.Quantity,
				ProductTitle:  item.ProductTitle,
				ProductURL:    item.ProductURL,
				PriceCents:    item.PriceCents,
				PriceCurrency: item.PriceCurrency,
				CreatedAt:     now,
				UpdatedAt:     now,
			}
			if err := uow.GetOrderItemRepo().CreateOrderItem(ctx, dalItem); err != nil {
				return nil, fmt.Errorf("failed to create order item: %w", err)
			}
			item.ID = dalItem.ID
			item.OrderID = dalItem.OrderID
			item.CreatedAt = dalItem.CreatedAt
			item.UpdatedAt = dalItem.UpdatedAt
		}
	}

	return orders, nil
}

func (s *OrderService) QueryOrders(
	ctx context.Context,
	uow *dal.UnitOfWork,
	req *dto.V1QueryOrdersRequest,
) ([]*core.Order, error) {
	dalReq := &models.QueryOrdersDalModel{
		IDs:         req.IDs,
		CustomerIDs: req.CustomerIDs,
	}

	if req.Page != nil && req.PageSize != nil && *req.Page > 0 && *req.PageSize > 0 {
		dalReq.Offset = (*req.Page - 1) * *req.PageSize
		dalReq.Limit = *req.PageSize
	} else {
		dalReq.Limit = 100
		dalReq.Offset = 0
	}

	dalOrders, err := uow.GetOrderRepo().QueryOrders(ctx, dalReq)
	if err != nil {
		return nil, fmt.Errorf("failed to query orders: %w", err)
	}

	if len(dalOrders) == 0 {
		return []*core.Order{}, nil
	}

	var orderItemsLookup map[int64][]models.V1OrderItemDal
	if req.IncludeOrderItems {
		orderIDs := make([]int64, len(dalOrders))
		for i, order := range dalOrders {
			orderIDs[i] = order.ID
		}

		itemsReq := &models.QueryOrderItemsDalModel{OrderIDs: orderIDs}
		dalItems, err := uow.GetOrderItemRepo().QueryOrderItems(ctx, itemsReq)
		if err != nil {
			return nil, fmt.Errorf("failed to query order items: %w", err)
		}

		orderItemsLookup = make(map[int64][]models.V1OrderItemDal)
		for _, item := range dalItems {
			orderItemsLookup[item.OrderID] = append(orderItemsLookup[item.OrderID], item)
		}
	}

	orders := make([]*core.Order, len(dalOrders))
	for i, dalOrder := range dalOrders {
		order := &core.Order{
			ID:                 dalOrder.ID,
			CustomerID:         dalOrder.CustomerID,
			DeliveryAddress:    dalOrder.DeliveryAddress,
			TotalPriceCents:    dalOrder.TotalPriceCents,
			TotalPriceCurrency: dalOrder.TotalPriceCurrency,
			CreatedAt:          dalOrder.CreatedAt,
			UpdatedAt:          dalOrder.UpdatedAt,
			Items:              []core.OrderItem{},
		}

		if req.IncludeOrderItems && orderItemsLookup != nil {
			if items, exists := orderItemsLookup[dalOrder.ID]; exists {
				order.Items = make([]core.OrderItem, len(items))
				for j, item := range items {
					order.Items[j] = core.OrderItem{
						ID:            item.ID,
						OrderID:       item.OrderID,
						ProductID:     item.ProductID,
						Quantity:      item.Quantity,
						ProductTitle:  item.ProductTitle,
						ProductURL:    item.ProductURL,
						PriceCents:    item.PriceCents,
						PriceCurrency: item.PriceCurrency,
						CreatedAt:     item.CreatedAt,
						UpdatedAt:     item.UpdatedAt,
					}
				}
			}
		}

		orders[i] = order
	}

	return orders, nil
}
