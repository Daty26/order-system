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

func (ns *NotificationService) Insert(orderId int, paymentId int, status model.NotificationStatus, userId int, message string) (model.Notification, error) {
	if status != "PENDING" && status != "SENT" && status != "FAILED" {
		return model.Notification{}, errors.New("enter right status")
	}
	if len(message) <= 0 {
		return model.Notification{}, errors.New("message body can't be empty")
	}

	var notification model.Notification = model.Notification{
		OrderID:   orderId,
		PaymentID: paymentId,
		Status:    status,
		Message:   message,
		UserID:    userId,
	}
	return ns.repo.Insert(notification)
}
func (ns *NotificationService) GetByID(id int) (model.Notification, error) {
	if id <= 0 {
		return model.Notification{}, errors.New("incorrect id")
	}
	return ns.repo.GetByID(id)
}
func (ns *NotificationService) GetByStatus(status model.NotificationStatus, userId int) ([]model.Notification, error) {
	if status != "PENDING" && status != "SENT" && status != "FAILED" {
		return []model.Notification{}, errors.New("Enter right status")
	}
	return ns.repo.GetByStatus(status, userId)
}

func (ns *NotificationService) GetAllByUserID(userId int) ([]model.Notification, error) {
	if userId < 0 {
		return []model.Notification{}, errors.New("invalid user id")
	}
	return ns.repo.GetAllByUserID(userId)
}
func (ns *NotificationService) UpdateStatusByID(id int, status model.NotificationStatus) (model.Notification, error) {
	if id <= 0 {
		return model.Notification{}, errors.New("incorrect id")
	}
	if status != "PENDING" && status != "SENT" && status != "FAILED" {
		return model.Notification{}, errors.New("enter right status")
	}
	return ns.repo.UpdateStatusByID(id, status)
}

func (ns *NotificationService) DeleteByID(id int) error {
	if id <= 0 {
		return errors.New("incorrect id")
	}
	return ns.repo.DeleteByID(id)
}
