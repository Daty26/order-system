package repository

import "github.com/Daty26/order-system/payment-service/internal/model"

type ProcessPaymentParams struct {
	OrderID     int
	Status      model.PaymentStatus
	AmountCents int64
	UserID      int
}
type UpdatePaymentParams struct {
	ID     int
	Status model.PaymentStatus
}
type GetAllByUserIDParams struct {
	ID     int
	Limit  int
	Offset int
}
