package api

type InsertNotificationRequest struct {
	OrderID   int    `json:"order_id"`
	PaymentID int    `json:"payment_id"`
	Status    string `json:"status"`
	Message   string `json:"message"`
}
type UpdateNotificationByStatus struct {
	Status string `json:"status"`
}
