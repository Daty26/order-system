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

func (ns *NotificationService) Insert(notification model.Notification) (model.Notification, error) {
	if notification.Status != "PENDING" && notification.Status != "SENT" && notification.Status != "FAILED" {
		return model.Notification{}, errors.New("enter right status")
	}
	if len(notification.Message) <= 0 {
		return model.Notification{}, errors.New("message body can't be empty")
	}
	return ns.repo.Insert(notification)
}
