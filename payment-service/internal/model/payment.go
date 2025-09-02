package model

type PaymentStatus string

const (
	PaymentPending   PaymentStatus = "PENDING"
	PaymentCompleted PaymentStatus = "COMPLETED"
	PaymentFailed    PaymentStatus = "FAILED"
)

type Payment struct {
	ID      int           `json:"id"`
	OrderID int           `json:"orderID"`
	Status  PaymentStatus `json:"status"`
	Amount  int           `json:"amount"`
}
