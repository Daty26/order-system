package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

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
	actor, ok := actorFromContext(r)
	if !ok {
		ErrorResponse(w, http.StatusUnauthorized, "unathorized")
		return
	}
	limit, offset, ok := parsePagination(r)
	if !ok {
		ErrorResponse(w, http.StatusBadRequest, "invalid pagination params")
		return
	}
	orders, err := h.service.GetOrders(r.Context(), actor, limit, offset)
	if err != nil {
		HandleErrors(
			w,
			r,
			h.logger,
			err,
			"failed to get orders",
			"user_id", actor.UserID,
		)
		return
	}
	resp := ToOrderResponses(orders)
	SuccessResp(w, http.StatusOK, resp)
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
	actor, ok := actorFromContext(r)
	if !ok {
		ErrorResponse(w, http.StatusUnauthorized, "unathorized")
		return
	}
	order, err := h.service.GetOrderByID(r.Context(), actor, id)
	if err != nil {
		HandleErrors(
			w,
			r,
			h.logger,
			err,
			"failed to get order by id",
			"id", id,
			"actor_id", actor.UserID,
		)
		return
	}
	resp := ToOrderResponse(order)
	SuccessResp(w, http.StatusOK, resp)
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
// func (h *OrderHandler) UpdateOrder(w http.ResponseWriter, r *http.Request) {
// 	if r.Context().Value("role") != "ADMIN" {
// 		ErrorResponse(w, http.StatusForbidden, "you are not allowed to update orders ")
// 		return
// 	}
// 	uid, ok := r.Context().Value("user_id").(float64)
// 	if !ok {
// 		ErrorResponse(w, http.StatusUnauthorized, "missing user id")
// 		return
// 	}
// 	userId := int(uid)
// 	var req struct {
// 		Items []struct {
// 			ProductId int `json:"product_id"`
// 			Quantity  int `json:"quantity"`
// 		} `json:"items"`
// 	}
// 	id, err := strconv.Atoi(chi.URLParam(r, "id"))
// 	if err != nil {
// 		ErrorResponse(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
// 		return
// 	}
// 	dec := json.NewDecoder(r.Body)
// 	dec.DisallowUnknownFields()
// 	err = dec.Decode(&req)
// 	if err != nil {
// 		ErrorResponse(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
// 		return
// 	}
// 	if len(req.Items) == 0 {
// 		ErrorResponse(w, http.StatusBadRequest, "items cannot be empty")
// 		return
// 	}
// 	items := make([]model.OrderCreatedEventItem, 0)
// 	for _, item := range req.Items {
// 		if item.ProductId <= 0 {
// 			ErrorResponse(w, http.StatusBadRequest, "product_id is required and must be > 0")
// 			return
// 		}
// 		if item.Quantity <= 0 {
// 			ErrorResponse(w, http.StatusBadRequest, "quantity can't be less than 0")
// 			return
// 		}
// 		items = append(items, model.OrderCreatedEventItem{
// 			ProductID: item.ProductId,
// 			Quantity:  item.Quantity,
// 		})
// 	}
// 	order := model.Orders{
// 		OrderID: id,
// 		UserID:  userId,
// 		Status:  "UPDATED",
// 		Items:   items,
// 	}
// 	updatedOrder, err := h.service.UpdateOrder(order)
// 	if err != nil {
// 		log.Println("Couldn't update order: " + err.Error())
// 		ErrorResponse(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
// 		return
// 	}
// 	SuccessResp(w, http.StatusOK, updatedOrder)
// }

// DeleteOrder godoc
// @Summary Delete order
// @Description Delete order with specified id
// @Success 200 "Ok"
// @Failure 400 {string} string "Invalid id"
// @Failure 500 {string} string "Couldn't delete order"
// @Router /orders/{id} [delete]
func (h *OrderHandler) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "invalid order id")
		return
	}
	actor, ok := actorFromContext(r)
	if !ok {
		ErrorResponse(w, http.StatusUnauthorized, "unathorized")
		return
	}

	err = h.service.DeleteOrder(r.Context(), actor, id)
	if err != nil {
		HandleErrors(
			w,
			r,
			h.logger,
			err,
			"failed to delete order",
			"order_id", id,
		)
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
	actor, ok := actorFromContext(r)
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
		UserID: actor.UserID,
		Items:  make([]service.CreatedOrderItemInput, 0, len(req.Items)),
	}
	for _, item := range req.Items {
		input.Items = append(input.Items, service.CreatedOrderItemInput{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		})
	}
	createdOrder, err := h.service.CreateOrder(r.Context(), actor, input)
	if err != nil {
		HandleErrors(
			w,
			r,
			h.logger,
			err,
			"failed to create order",
			"user_id", actor.UserID,
		)
		return
	}
	SuccessResp(w, http.StatusCreated, ToOrderResponse(createdOrder))
}

func (h *OrderHandler) CancelOrder(w http.ResponseWriter, r *http.Request) {
	orderID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || orderID <= 0 {
		ErrorResponse(w, http.StatusBadRequest, "invalid order id")
		return
	}
	actor, ok := actorFromContext(r)
	if !ok {
		ErrorResponse(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	order, err := h.service.CancelOrder(r.Context(), actor, orderID)
	if err != nil {
		HandleErrors(
			w,
			r,
			h.logger,
			err,
			"failed to cancel the order",
			"order_id", orderID,
			"user_id", actor.UserID,
		)
		return
	}
	SuccessResp(w, http.StatusOK, ToOrderResponse(order))
}
