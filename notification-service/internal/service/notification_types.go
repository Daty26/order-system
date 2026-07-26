package service

import "github.com/Daty26/order-system/notification-service/internal/model"

type InsertInput struct {
	OrderID   int
	PaymentID int
	Status    model.NotificationStatus
	UserID    int
	Message   string
}
type GetByStatusInput struct {
	UserID int
	Status model.NotificationStatus
	Limit  int
	Offset int
}
type GetAllByUserIDInput struct {
	UserID int
	Limit  int
	Offset int
}
type UpdateStatusInput struct {
	ID     int
	Status model.NotificationStatus
}
