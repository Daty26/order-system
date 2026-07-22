package model

type PaymentStatus string

const (
	PaymentPending   PaymentStatus = "PENDING"
	PaymentCompleted PaymentStatus = "COMPLETED"
	PaymentFailed    PaymentStatus = "FAILED"
)

type Payment struct {
	ID          int           `json:"payment_id"`
	OrderID     int           `json:"order_id"`
	Status      PaymentStatus `json:"status"`
	AmountCents int64         `json:"amount_cents"`
	UserID      int           `json:"user_id"`
}
