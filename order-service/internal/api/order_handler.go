package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Daty26/order-system/order-service/internal/model"
	"github.com/Daty26/order-system/order-service/internal/service"
	"github.com/go-chi/chi/v5"
)

type OrderHandler struct {
	service *service.OrderService
	logger  *slog.Logger
}

func NewOrderHandler(s *service.OrderService, logger *slog.Logger) *OrderHandler {
	return &OrderHandler{service: s, logger: logger}
}

// GetOrders godoc
// @Summary Get all orders
// @Description List all orders currently stored
// @Produce json
// @Success 200 {array} model.Order
// @Router /orders [get]
func (h *OrderHandler) GetOrders(w http.ResponseWriter, r *http.Request) {
	role, ok := r.Context().Value("role").(string)
	if !ok {
		ErrorResponse(w, http.StatusUnauthorized, "unathorized")
		return
	}

	userIDRaw, ok := r.Context().Value("user_id").(float64)
	if !ok {
		ErrorResponse(w, http.StatusUnauthorized, "unathorized")
		return
	}
	userID := int(userIDRaw)

	limit, offset, ok := parsePagination(r)
	if !ok {
		ErrorResponse(w, http.StatusBadRequest, "invalid pagination params")
		return
	}

	var (
		orders []model.Orders
		err    error
	)
	if role == "ADMIN" {
		orders, err = h.service.GetOrders(r.Context(), limit, offset)
	} else {
		orders, err = h.service.GetOrdersByUserId(r.Context(), userID, limit, offset)
	}

	if err != nil {
		if errors.Is(err, service.ErrInvalidOrder) {
			ErrorResponse(w, http.StatusBadRequest, "invalid order request")
			return
		}
		h.logger.ErrorContext(r.Context(), "failed to get orders",
			"error", err,
			"user_id", userID,
			"role", role,
		)

		ErrorResponse(w, http.StatusInternalServerError, "something went wrong")
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
		ErrorResponse(w, http.StatusBadRequest, "invalid order id")
		return
	}
	role, ok := r.Context().Value("role").(string)
	if !ok {
		ErrorResponse(w, http.StatusUnauthorized, "unathorized")
		return
	}
	userIDRaw, ok := r.Context().Value("user_id").(float64)
	if !ok {
		ErrorResponse(w, http.StatusUnauthorized, "unathorized")
		return
	}
	userID := int(userIDRaw)
	order, err := h.service.GetOrderByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ErrorResponse(w, http.StatusNotFound, "order not found")
			return
		}
		h.logger.ErrorContext(r.Context(), "failed to get order by id",
			"error", err,
			"order_id", id,
		)

		ErrorResponse(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	if role != "ADMIN" && order.UserID != userID {
		ErrorResponse(w, http.StatusForbidden, "you are not allowed to view this order")
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
		ErrorResponse(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err = dec.Decode(&req)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}
	if len(req.Items) == 0 {
		ErrorResponse(w, http.StatusBadRequest, "items cannot be empty")
		return
	}
	items := make([]model.OrderCreatedEventItem, 0)
	for _, item := range req.Items {
		if item.ProductId <= 0 {
			ErrorResponse(w, http.StatusBadRequest, "product_id is required and must be > 0")
			return
		}
		if item.Quantity <= 0 {
			ErrorResponse(w, http.StatusBadRequest, "quantity can't be less than 0")
			return
		}
		items = append(items, model.OrderCreatedEventItem{
			ProductID: item.ProductId,
			Quantity:  item.Quantity,
		})
	}
	order := model.Orders{
		OrderID: id,
		UserID:  userId,
		Status:  "UPDATED",
		Items:   items,
	}
	updatedOrder, err := h.service.UpdateOrder(order)
	if err != nil {
		log.Println("Couldn't update order: " + err.Error())
		ErrorResponse(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
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
		ErrorResponse(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
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
		ErrorResponse(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	var req CreatedOrderRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}
	input := service.CreatedOrderInput{
		UserID: int(uid),
		Items:  make([]service.CreatedOrderItemInput, 0, len(req.Items)),
	}
	for _, item := range req.Items {
		input.Items = append(input.Items, service.CreatedOrderItemInput{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		})
	}
	createdOrder, err := h.service.CreateOrder(r.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidOrder):
			ErrorResponse(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, service.ErrProductNotFound):
			ErrorResponse(w, http.StatusNotFound, err.Error())
		default:
			h.logger.ErrorContext(r.Context(), "failed to create order", err)
			ErrorResponse(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		}
		return
	}
	SuccessResp(w, http.StatusCreated, createdOrder)
}
