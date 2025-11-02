package model

type PaymentStatus string

const PaymentPending PaymentStatus = "PENDING"
const PaymentCompleted PaymentStatus = "CREATED"
const PaymentFailed PaymentStatus = "FAILED"

type PaymentCreated struct {
	PaymentID int           `json:"payment_id"`
	OrderID   int           `json:"order_id"`
	Status    PaymentStatus `json:"status"`
}
