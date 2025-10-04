package v1

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Lamafout/online-store-api/core/models/common"
	"github.com/Lamafout/online-store-api/core/models/dto"
	"github.com/Lamafout/online-store-api/internal/bll/services"
	"github.com/go-chi/chi/v5"
)

type OrderHandler struct {
	service *services.OrderService
}

func NewOrderHandler(service *services.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

func (h *OrderHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/", h.CreateOrder)
	r.Post("/batch-create", h.BatchCreateOrders) // Новая ручка
	r.Post("/query", h.QueryOrders) // Новая ручка
	r.Get("/{id}", h.GetOrder)
	return r
}

// @Summary Create a new order
// @Description Creates a new order with items
// @Tags Orders
// @Accept json
// @Produce json 
// @Param order body models.Order true "Order data"
// @Success 201 {object} models.Order
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /orders [post]
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var order common.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	if err := h.service.CreateOrder(r.Context(), &order); err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

// @Summary Batch create orders
// @Description Creates multiple orders in batch
// @Tags Orders
// @Accept json
// @Produce json
// @Param request body models.V1CreateOrderRequest true "Orders data"
// @Success 201 {object} models.V1CreateOrderResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /orders/batch-create [post]
func (h *OrderHandler) BatchCreateOrders(w http.ResponseWriter, r *http.Request) {
	var req dto.V1CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}


	orders := make([]*common.Order, len(req.Orders))
	for i, orderReq := range req.Orders {
		items := make([]common.OrderItem, len(orderReq.OrderItems))
		for j, itemReq := range orderReq.OrderItems {
			items[j] = common.OrderItem{
				ProductID:     itemReq.ProductID,
				Quantity:      itemReq.Quantity,
				ProductTitle:  itemReq.ProductTitle,
				ProductURL:    itemReq.ProductURL,
				PriceCents:    itemReq.PriceCents,
				PriceCurrency: itemReq.PriceCurrency,
			}
		}
		
		orders[i] = &common.Order{
			CustomerID:         orderReq.CustomerID,
			DeliveryAddress:    orderReq.DeliveryAddress,
			TotalPriceCents:    orderReq.TotalPriceCents,
			TotalPriceCurrency: orderReq.TotalPriceCurrency,
			Items:              items,
		}
	}

	createdOrders, err := h.service.BatchCreateOrders(r.Context(), orders)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	// Конвертируем обратно в response
	responseOrders := make([]common.Order, len(createdOrders))
	for i, order := range createdOrders {
		responseOrders[i] = *order
	}

	response := dto.V1CreateOrderResponse{
		Orders: responseOrders,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// @Summary Query orders
// @Description Query orders with filters
// @Tags Orders
// @Accept json
// @Produce json
// @Param request body models.V1QueryOrdersRequest true "Query filters"
// @Success 200 {object} models.V1QueryOrdersResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /orders/query [post]
func (h *OrderHandler) QueryOrders(w http.ResponseWriter, r *http.Request) {
	var req dto.V1QueryOrdersRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	// Валидация базовых параметров
	if req.Page != nil && *req.Page < 1 {
		http.Error(w, `{"error": "Page must be greater than 0"}`, http.StatusBadRequest)
		return
	}
	
	if req.PageSize != nil && *req.PageSize < 1 {
		http.Error(w, `{"error": "PageSize must be greater than 0"}`, http.StatusBadRequest)
		return
	}

	orders, err := h.service.QueryOrders(r.Context(), &req)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	response := dto.V1QueryOrdersResponse{
		Orders: make([]common.Order, len(orders)),
	}
	for i, order := range orders {
		response.Orders[i] = *order
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// @Summary Get an order by ID
// @Description Retrieves an order with its items by ID
// @Tags Orders
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} models.Order
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /orders/{id} [get]
func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, `{"error": "Invalid order ID"}`, http.StatusBadRequest)
		return
	}

	order, err := h.service.GetOrder(r.Context(), id)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}