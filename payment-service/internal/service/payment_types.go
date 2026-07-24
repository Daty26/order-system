package service

import (
	"context"

	"github.com/Daty26/order-system/payment-service/internal/model"
)

type OrderSummary struct {
	OrderID          int   `json:"order_id"`
	TotalAmountCents int64 `json:"total_amount_cents"`
}

type OrderClient interface {
	GetOrder(ctx context.Context, orderID int, authHeader string) (OrderSummary, error)
}
type ProcessPaymentInput struct {
	OrderID    int
	UserID     int
	AuthHeader string
}
type UpdatePaymentInput struct {
	ID     int
	Status model.PaymentStatus
}
type GetAllByUserIDInput struct {
	ID     int
	Limit  int
	Offset int
}
