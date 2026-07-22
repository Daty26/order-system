package service

type ProcessPaymentInput struct {
	OrderID        int
	AmountCents    int64
	UserID         int
	IdempotencyKey string
}
