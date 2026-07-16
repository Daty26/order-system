package api

import "github.com/Daty26/order-system/order-service/internal/model"

type CreatedOrderRequest struct {
	Items []CreatedOrderItemRequest `json:"items"`
}

type CreatedOrderItemRequest struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}
type OrderResponse struct {
	OrderID          int                 `json:"order_id"`
	Status           model.OrderStatus   `json:"status"`
	TotalAmountCents int64               `json:"total_amount_cents"`
	Items            []OrderItemResponse `json:"items"`
}
type OrderItemResponse struct {
	ProductID      int   `json:"product_id"`
	Quantity       int   `json:"quantity"`
	UnitPriceCents int64 `json:"unit_price_cents"`
}

func ToOrderResponse(order model.Orders) OrderResponse {
	items := make([]OrderItemResponse, 0, len(order.Items))
	for _, item := range order.Items {
		items = append(items, OrderItemResponse{
			ProductID:      item.OrderID,
			Quantity:       item.Quantity,
			UnitPriceCents: item.UnitPriceCents,
		})
	}
	return OrderResponse{
		OrderID:          order.OrderID,
		Status:           order.Status,
		TotalAmountCents: order.TotalAmountCents,
		Items:            items,
	}
}

func ToOrderResponses(orders []model.Orders) []OrderResponse {
	responses := make([]OrderResponse, 0, len(orders))
	for _, order := range orders {
		responses = append(responses, ToOrderResponse(order))
	}
	return responses
}
