package services

import (
	"context"
	"fmt"
	"time"

	core "github.com/Lamafout/online-store-api/core/models"
	"github.com/Lamafout/online-store-api/internal/dal/unit_of_work"
	"github.com/Lamafout/online-store-api/internal/dal/models"
	"github.com/go-playground/validator/v10"
)

type OrderService struct {
	uow      *dal.UnitOfWork
	validate *validator.Validate
}

func NewOrderService(uow *dal.UnitOfWork) *OrderService {
	return &OrderService{
		uow:      uow,
		validate: validator.New(),
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, order *core.Order) error {
	if err := s.validate.Struct(order); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	if err := s.uow.Begin(ctx); err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer s.uow.Rollback()

	dalOrder := &models.V1OrderDal{
		CustomerID:        order.CustomerID,
		DeliveryAddress:   order.DeliveryAddress,
		TotalPriceCents:   order.TotalPriceCents,
		TotalPriceCurrency: order.TotalPriceCurrency,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
	if err := s.uow.GetOrderRepo().CreateOrder(ctx, dalOrder); err != nil {
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
		if err := s.uow.GetOrderItemRepo().CreateOrderItem(ctx, dalItem); err != nil {
			return fmt.Errorf("failed to create order item: %w", err)
		}
		item.ID = dalItem.ID
	}

	if err := s.uow.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func (s *OrderService) GetOrder(ctx context.Context, id int64) (*core.Order, error) {
	dalOrder, err := s.uow.GetOrderRepo().GetOrderByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}
	dalItems, err := s.uow.GetOrderItemRepo().GetOrderItemsByOrderID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get order items: %w", err)
	}

	order := &core.Order{
		ID:                 dalOrder.ID,
		CustomerID:        dalOrder.CustomerID,
		DeliveryAddress:   dalOrder.DeliveryAddress,
		TotalPriceCents:   dalOrder.TotalPriceCents,
		TotalPriceCurrency: dalOrder.TotalPriceCurrency,
		CreatedAt:         dalOrder.CreatedAt,
		UpdatedAt:         dalOrder.UpdatedAt,
		Items:              make([]core.OrderItem, len(dalItems)),
	}
	for i, dalItem := range dalItems {
		order.Items[i] = core.OrderItem{
			ID:            dalItem.ID,
			OrderID:       dalItem.OrderID,
			ProductID:     dalItem.ProductID,
			Quantity:      dalItem.Quantity,
			ProductTitle:  dalItem.ProductTitle,
			ProductURL:    dalItem.ProductURL,
			PriceCents:    dalItem.PriceCents,
			PriceCurrency: dalItem.PriceCurrency,
			CreatedAt:     dalItem.CreatedAt,
			UpdatedAt:     dalItem.UpdatedAt,
		}
	}
	return order, nil
}