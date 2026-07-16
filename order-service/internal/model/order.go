package model

import "time"

type OrderStatus string

const (
	OrderPending    OrderStatus = "PENDING"
	OrderUpdated    OrderStatus = "CONFIRMED"
	OrderProcessing OrderStatus = "PROCESSING"
	OrderDelivered  OrderStatus = "DELIVERED"
	OrderCancelled  OrderStatus = "CANCELLED"
)

type Orders struct {
	Version          int         `json:"version"`
	OrderID          int         `json:"order_id"`
	UserID           int         `json:"user_id"`
	Status           OrderStatus `json:"status"`
	TotalAmountCents int64       `json:"total_amount_cents"`
	CreatedAt        time.Time   `json:"created_at"`
	Items            []OrderItem `json:"items"`
}
type OrderItem struct {
	ID             int
	OrderID        int
	ProductID      int
	Quantity       int
	UnitPriceCents int64
}
type OrderCreatedEvent struct {
	Version          int                     `json:"version"`
	OrderID          int                     `json:"order_id"`
	UserID           int                     `json:"user_id"`
	Status           OrderStatus             `json:"status"`
	TotalAmountCents int64                   `json:"total_amount_cents"`
	CreatedAt        time.Time               `json:"created_at"`
	Items            []OrderCreatedEventItem `json:"items"`
}
type OrderCreatedEventItem struct {
	ProductID  int   `json:"product_id"`
	Quantity   int   `json:"quantity"`
	PriceCents int64 `json:"price_cents"`
}
