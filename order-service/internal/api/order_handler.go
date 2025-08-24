package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/Daty26/order-system/order-service/internal/service"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

type OrderHandler struct {
	service *service.OrderService
}

func NewOrderHandler(s *service.OrderService) *OrderHandler {
	return &OrderHandler{service: s}
}

// GetOrders godoc
// @Summary Get all orders
// @Description List all orders currently stored
// @Produce json
// @Success 200 {array} model.Order
// @Router /orders [get]
func (h *OrderHandler) GetOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.service.GetOrders()
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "couldn't fetch orders")
		return
	}
	SuccessResp(w, http.StatusOK, orders)
}

// GetOrderByID godoc
// @Summary Get order by ID
// @Description Fetch a single order by its ID
// @Param id path int true "Order ID"
// @Produce json
// @Success 200 {object} model.Order
// @Failure 400 {string} string "Invalid ID"
// @Failure 404 {string} string "Order not found"
// @Router /orders/{id} [get]
func (h *OrderHandler) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid ID")
		return
	}
	order, err := h.service.GetOrderByID(id)
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "Couldn't find order with specified id")
		return
	}
	SuccessResp(w, http.StatusOK, order)
}

// UpdateOrder godoc
// @Summary Update an existing order
// @Description Update the item and amount of an order by ID
// @Param id path int true "Order ID"
// @Accept json
// @Produce json
// @Success 200 {object} model.Order
// @Failure 400 {string} string "Invalid ID or request body"
// @Failure 404 {string} string "Order not found"
// @Failure 500 {string} string "Couldn't update order"
// @Router /orders/{id} [put]
func (h *OrderHandler) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Item   string `json:"item"`
		Amount int    `json:"amount"`
	}
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid ID")
		return
	}
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Couldn't convert the req body to specified format")
		return
	}
	order, err := h.service.UpdateOrder(id, req.Item, req.Amount)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Couldn't update order")
		return
	}
	SuccessResp(w, http.StatusOK, order)
}

// DeleteOrder godoc
// @Summary Delete order
// @Description Delete order with specified id
// @Success 20o "Ok"
// @Failure 400 {string} string "Invalid id"
// @Failure 500 {string} string "Couldn't delete order"
// @Router /orders/{id} [delete]
func (h *OrderHandler) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid id")
		return
	}
	err = h.service.DeleteOrder(id)
	if errors.Is(err, sql.ErrNoRows) {
		ErrorResponse(w, http.StatusNotFound, "Couldn't find order with such id")
		return
	}
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Couldn't delete order")
		return
	}
	SuccessResp(w, http.StatusOK, nil)
}

// CreateOrder godoc
// @Summary Create new order
// @Description Create new order with item and amount
// @Accept json
// @Produce json
// @Param order body model.Order true "Order data"
// @Success 201 {object} model.Order
// @Failure 400 {string} string "Invalid request"
// @Router /orders [post]
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Item   string `json:"item"`
		Amount int    `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "invalid request format")
		return
	}
	order, err := h.service.CreateOrder(req.Item, req.Amount)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	SuccessResp(w, http.StatusCreated, order)
}
