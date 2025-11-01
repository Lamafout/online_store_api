package v1

import (
	"encoding/json"
	"net/http"

	"github.com/Lamafout/online-store-api/core/models/common"
	"github.com/Lamafout/online-store-api/core/models/dto"
	"github.com/Lamafout/online-store-api/internal/bll/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
)

type OrderHandler struct {
	service *services.OrderService
}

func NewOrderHandler(db *sqlx.DB, service *services.OrderService) *OrderHandler {
	return &OrderHandler{
		service: service,
	}
}

func (h *OrderHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/batch-create", h.BatchCreateOrders)
	r.Post("/query", h.QueryOrders)
	return r
}

// @Summary Batch create orders
// @Description Creates multiple orders in batch
// @Tags Orders
// @Accept json
// @Produce json
// @Param request body dto.V1CreateOrderRequest true "Orders data"
// @Success 201 {object} dto.V1CreateOrderResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /orders/batch-create [post]
func (h *OrderHandler) BatchCreateOrders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req dto.V1CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	validate := validator.New()

	if err := validate.Struct(req); err != nil {
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

	createdOrders, err := h.service.BatchCreateOrders(ctx, orders)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}


	responseOrders := make([]common.Order, len(createdOrders))
	for i, order := range createdOrders {
		responseOrders[i] = *order
	}

	response := dto.V1CreateOrderResponse{
		Orders: responseOrders,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(response)
}

// @Summary Query orders
// @Description Query orders with filters
// @Tags Orders
// @Accept json
// @Produce json
// @Param request body dto.V1QueryOrdersRequest true "Query filters"
// @Success 200 {object} dto.V1QueryOrdersResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /orders/query [post]
func (h *OrderHandler) QueryOrders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req dto.V1QueryOrdersRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	if req.Page != nil && *req.Page < 1 {
		http.Error(w, `{"error": "Page must be greater than 0"}`, http.StatusBadRequest)
		return
	}

	if req.PageSize != nil && *req.PageSize < 1 {
		http.Error(w, `{"error": "PageSize must be greater than 0"}`, http.StatusBadRequest)
		return
	}

	validate := validator.New()

	if err := validate.Struct(req); err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	orders, err := h.service.QueryOrders(ctx, &req)
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
	_ = json.NewEncoder(w).Encode(response)
}
