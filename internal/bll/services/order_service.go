package services

import (
	"context"
	"fmt"
	"time"

	core "github.com/Lamafout/online-store-api/core/models/common"
	"github.com/Lamafout/online-store-api/core/models/dto"
	"github.com/Lamafout/online-store-api/core/models/messages"
	"github.com/Lamafout/online-store-api/internal/config"
	"github.com/Lamafout/online-store-api/internal/dal/models"
	dal "github.com/Lamafout/online-store-api/internal/dal/unit_of_work"
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
)

type OrderService struct {
	validate  *validator.Validate
	publisher RabbitPublisher
	rs        config.RabbitMqSettings
	db        *sqlx.DB
}

func NewOrderService(publisher RabbitPublisher, rs config.RabbitMqSettings, db *sqlx.DB) *OrderService {
	return &OrderService{
		validate:  validator.New(),
		publisher: publisher,
		rs:        rs,
		db:        db,
	}
}

func (s *OrderService) BatchCreateOrders(
	ctx context.Context,
	orders []*core.Order,
) ([]*core.Order, error) {

	uow := dal.NewUnitOfWork(s.db)

	if err := uow.Begin(ctx); err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}

	defer uow.Rollback()

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

	var msgs []any
	for _, order := range orders {
		msg := messages.OrderCreatedMessage{
			ID:                 order.ID,
			CustomerID:         order.CustomerID,
			DeliveryAddress:    order.DeliveryAddress,
			TotalPriceCents:    order.TotalPriceCents,
			TotalPriceCurrency: order.TotalPriceCurrency,
			CreatedAt:          order.CreatedAt,
		}

		for _, item := range order.Items {
			msg.OrderItems = append(msg.OrderItems, messages.OrderItemMessage{
				ID:            item.ID,
				OrderID:       item.OrderID,
				ProductID:     item.ProductID,
				ProductTitle:  item.ProductTitle,
				ProductURL:    item.ProductURL,
				Quantity:      item.Quantity,
				PriceCents:    item.PriceCents,
				PriceCurrency: item.PriceCurrency,
			})
		}

		msgs = append(msgs, msg)
	}

	if err := s.publisher.Publish(ctx, msgs, s.rs.OrderCreateQueue); err != nil {
		return nil, fmt.Errorf("failed to publish order-created messages: %w", err)
	}

	if err := uow.Commit(); err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}

	return orders, nil
}

func (s *OrderService) QueryOrders(
	ctx context.Context,
	req *dto.V1QueryOrdersRequest,
) ([]*core.Order, error) {

	uow := dal.NewUnitOfWork(s.db)

	if err := uow.Begin(ctx); err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}

	defer uow.Rollback()

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

	if err := uow.Commit(); err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}

	return orders, nil
}
