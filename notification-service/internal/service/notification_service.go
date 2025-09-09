package service

import (
	"errors"
	"github.com/Daty26/order-system/notification-service/internal/model"
	"github.com/Daty26/order-system/notification-service/internal/repository"
)

type NotificationService struct {
	repo repository.NotificationRepo
}

func NewNotificationService(repo repository.NotificationRepo) *NotificationService {
	return &NotificationService{repo: repo}
}

func (ns *NotificationService) Insert(order_id int, payment_id int, status model.NotificationStatus, message string) (model.Notification, error) {
	if status != "PENDING" && status != "SENT" && status != "FAILED" {
		return model.Notification{}, errors.New("enter right status")
	}
	if len(message) <= 0 {
		return model.Notification{}, errors.New("message body can't be empty")
	}

	var notification model.Notification = model.Notification{
		OrderID:   order_id,
		PaymentID: payment_id,
		Status:    status,
		Message:   message,
	}
	return ns.repo.Insert(notification)
}
