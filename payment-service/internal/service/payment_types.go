package service

import "context"

type ProcessPaymentInput struct {
	OrderID    int
	UserID     int
	AuthHeader string
}
type OrderSummary struct {
	OrderID          int   `json:"order_id"`
	TotalAmountCents int64 `json:"total_amount_cents"`
}

type OrderClient interface {
	GetOrder(ctx context.Context, orderID int, authHeader string) (OrderSummary, error)
}
