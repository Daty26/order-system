package repository

import "github.com/Daty26/order-system/notification-service/internal/model"

type InsertParams struct {
	OrderID   int
	PaymentID int
	Status    model.NotificationStatus
	Message   string
	UserID    int
}
type GetAllByUserIDParams struct {
	UserID int
	Limit  int
	Offset int
}
type GetByStatusParams struct {
	Status model.NotificationStatus
	UserID int
	Limit  int
	Offset int
}
type UpdateStatusByIDParams struct {
	ID     int
	Status model.NotificationStatus
}
