package api

import "github.com/Daty26/order-system/payment-service/internal/model"

type ProcessPaymentRequest struct {
	OrderID int `json:"order_id"`
}
type UpdatePaymentRequest struct {
	Status model.PaymentStatus `json:"status"`
}
