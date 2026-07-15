package api

import "github.com/Daty26/order-system/order-service/internal/model"

type CreatedOrderRequest struct {
	Items []CreatedOrderItemRequest `json:"items"`
}

type CreatedOrderItemRequest struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}
type CreatedOrderResponse struct {
	OrderID          int                 `json:"order_id"`
	Status           model.OrderStatus   `json:"status"`
	TotalAmountCents int64               `json:"total_amount_cents"`
	Items            []OrderItemResponse `json:"items"`
}
type OrderItemResponse struct {

}
