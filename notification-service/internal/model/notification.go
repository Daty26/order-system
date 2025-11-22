package model

import "time"

type NotificationStatus string

const (
	NotificationPending NotificationStatus = "PENDING"
	NotificationSent    NotificationStatus = "SENT"
	NotificationFailed  NotificationStatus = "FAILED"
)

type Notification struct {
	ID        int                `json:"id"`
	OrderID   int                `json:"orderID"`
	PaymentID int                `json:"paymentID"`
	Status    NotificationStatus `json:"status"`
	Message   string             `json:"message"`
	UserID    int                `json:"user_id"`
	CreatedAt time.Time          `json:"created_at"`
}
