package service

import (
	"context"
	"errors"
	"strings"

	"github.com/Daty26/order-system/notification-service/internal/model"
	"github.com/Daty26/order-system/notification-service/internal/repository"
)

type NotificationService struct {
	repo repository.NotificationRepo
}

func NewNotificationService(repo repository.NotificationRepo) *NotificationService {
	return &NotificationService{repo: repo}
}

func (s *NotificationService) Insert(ctx context.Context, input InsertInput) (model.Notification, error) {
	if !input.Status.IsValid() {
		return model.Notification{}, ErrInvalidStatus
	}
	if strings.TrimSpace(input.Message) == "" {
		return model.Notification{}, ErrInvalidMessage
	}
	if input.OrderID <= 0 || input.PaymentID <= 0 || input.UserID <= 0 {
		return model.Notification{}, ErrInvalidID
	}
	params := repository.InsertParams{
		OrderID:   input.OrderID,
		PaymentID: input.PaymentID,
		Status:    input.Status,
		UserID:    input.UserID,
		Message:   input.Message,
	}
	notification, err := s.repo.Insert(ctx, params)
	if err != nil {
		if errors.Is(err, repository.ErrNotificationAlreadyExists) {
			return model.Notification{}, ErrNotificationAlreadyExists
		}
		return model.Notification{}, err
	}
	return notification, nil
}

func (s *NotificationService) GetByID(ctx context.Context, id int) (model.Notification, error) {
	if id <= 0 {
		return model.Notification{}, ErrInvalidID
	}
	return s.repo.GetByID(ctx, id)
}

func (s *NotificationService) GetByStatus(ctx context.Context, input GetByStatusInput) ([]model.Notification, error) {
	if !input.Status.IsValid() {
		return []model.Notification{}, ErrInvalidStatus
	}
	if input.UserID <= 0 {
		return []model.Notification{}, ErrInvalidID
	}

	params := repository.GetByStatusParams{
		UserID: input.UserID,
		Status: input.Status,
		Limit:  input.Limit,
		Offset: input.Offset,
	}
	return s.repo.GetByStatus(ctx, params)
}

func (s *NotificationService) GetAllByUserID(ctx context.Context, input GetAllByUserIDInput) ([]model.Notification, error) {
	if input.UserID <= 0 {
		return []model.Notification{}, ErrInvalidID
	}
	params := repository.GetAllByUserIDParams{
		UserID: input.UserID,
		Limit:  input.Limit,
		Offset: input.Offset,
	}
	return s.repo.GetAllByUserID(ctx, params)
}

func (s *NotificationService) GetAll(ctx context.Context, limit, offset int) ([]model.Notification, error) {
	return s.repo.GetAll(ctx, limit, offset)
}

func (s *NotificationService) UpdateStatus(ctx context.Context, input UpdateStatusInput) (model.Notification, error) {
	if input.ID <= 0 {
		return model.Notification{}, ErrInvalidID
	}
	if !input.Status.IsValid() {
		return model.Notification{}, ErrInvalidStatus
	}
	params := repository.UpdateStatusParams{
		ID:     input.ID,
		Status: input.Status,
	}
	return s.repo.UpdateStatus(ctx, params)
}

func (s *NotificationService) DeleteByID(ctx context.Context, id int) error {
	if id <= 0 {
		return ErrInvalidID
	}
	return s.repo.DeleteByID(ctx, id)
}
