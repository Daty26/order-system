package repository

import "github.com/Daty26/order-system/payment-service/internal/model"

type ProcessPaymentParams struct {
	OrderID     int
	Status      model.PaymentStatus
	AmountCents int64
	UserID      int
}
