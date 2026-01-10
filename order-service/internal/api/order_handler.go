package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/Daty26/order-system/order-service/internal/model"
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
	role := r.Context().Value("role").(string)
	userId := int(r.Context().Value("user_id").(float64))
	if role == "ADMIN" {
		orders, err := h.service.GetOrders()
		if err != nil {
			ErrorResponse(w, http.StatusInternalServerError, "couldn't fetch orders")
			return
		}
		SuccessResp(w, http.StatusOK, orders)
		return
	}
	orders, err := h.service.GetOrdersByUserId(userId)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "couldn't fetch orders of specified user")
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
// @Description Update the item and quantity of an order by ID
// @Param id path int true "Order ID"
// @Accept json
// @Produce json
// @Success 200 {object} model.Order
// @Failure 400 {string} string "Invalid ID or request body"
// @Failure 404 {string} string "Order not found"
// @Failure 500 {string} string "Couldn't update order"
// @Router /orders/{id} [put]
func (h *OrderHandler) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	if r.Context().Value("role") != "ADMIN" {
		ErrorResponse(w, http.StatusForbidden, "you are not allowed to update orders ")
		return
	}
	uid, ok := r.Context().Value("user_id").(float64)
	if !ok {
		ErrorResponse(w, http.StatusUnauthorized, "missing user id")
		return
	}
	userId := int(uid)
	var req struct {
		Items []struct {
			ProductId int `json:"product_id"`
			Quantity  int `json:"quantity"`
		} `json:"items"`
	}
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid ID")
		return
	}
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err = dec.Decode(&req)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Couldn't convert the req body to specified format")
		return
	}
	if len(req.Items) == 0 {
		ErrorResponse(w, http.StatusBadRequest, "items cannot be empty")
		return
	}
	items := make([]model.OrderItems, 0)
	for _, item := range req.Items {
		if item.ProductId <= 0 {
			ErrorResponse(w, http.StatusBadRequest, "product_id is required and must be > 0")
			return
		}
		if item.Quantity <= 0 {
			ErrorResponse(w, http.StatusBadRequest, "quantity can't be less than 0")
			return
		}
		items = append(items, model.OrderItems{
			ProductID: item.ProductId,
			Quantity:  item.Quantity,
		})
	}
	order := model.Orders{
		ID:     id,
		UserID: userId,
		Status: "UPDATED",
		Items:  items,
	}
	updatedOrder, err := h.service.UpdateOrder(order)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Couldn't update order")
		return
	}
	SuccessResp(w, http.StatusOK, updatedOrder)
}

// DeleteOrder godoc
// @Summary Delete order
// @Description Delete order with specified id
// @Success 200 "Ok"
// @Failure 400 {string} string "Invalid id"
// @Failure 500 {string} string "Couldn't delete order"
// @Router /orders/{id} [delete]
func (h *OrderHandler) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	if r.Context().Value("role") != "ADMIN" {
		ErrorResponse(w, http.StatusForbidden, "you are not allowed to delete orders")
		return
	}
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
// @Description Create new order with item and quantity
// @Accept json
// @Produce json
// @Param order body model.Order true "Order data"
// @Success 201 {object} model.Order
// @Failure 400 {string} string "Invalid request"
// @Router /orders [post]
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	uid, ok := r.Context().Value("user_id").(float64)
	if !ok {
		ErrorResponse(w, http.StatusUnauthorized, "missing user id")
		return
	}
	userId := int(uid)

	var req struct {
		Items []struct {
			ProductId int `json:"product_id"`
			Quantity  int `json:"quantity"`
		} `json:"items"`
	}
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "invalid request format: "+err.Error())
		return
	}
	if len(req.Items) == 0 {
		ErrorResponse(w, http.StatusBadRequest, "can't create empty order")
		return
	}
	items := make([]model.OrderItems, 0)
	for _, item := range req.Items {
		if item.ProductId <= 0 {
			ErrorResponse(w, http.StatusBadRequest, "product_id is required and must be > 0")
			return
		}
		if item.Quantity <= 0 {
			ErrorResponse(w, http.StatusBadRequest, "quantity can't be less than 0")
			return
		}
		items = append(items, model.OrderItems{
			ProductID: item.ProductId,
			Quantity:  item.Quantity,
		})
	}
	order := model.Orders{
		UserID: userId,
		Items:  items,
	}
	createdOrder, err := h.service.CreateOrder(order)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	SuccessResp(w, http.StatusCreated, createdOrder)
}
