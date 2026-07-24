package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Daty26/order-system/payment-service/internal/kafka"
	"github.com/Daty26/order-system/payment-service/internal/model"
	"github.com/Daty26/order-system/payment-service/internal/repository"
)

type PaymentService struct {
	paymentRep repository.PaymentRep
	producer   *kafka.KafkaProducer
	client     OrderClient
}

func NewPaymentService(payRep repository.PaymentRep, prod *kafka.KafkaProducer, client OrderClient) *PaymentService {
	return &PaymentService{paymentRep: payRep, producer: prod, client: client}
}

func (s *PaymentService) ProcessPayment(ctx context.Context, input ProcessPaymentInput) (model.Payment, error) {
	if input.UserID <= 0 || input.OrderID <= 0 {
		return model.Payment{}, ErrInvalidInput
	}
	order, err := s.client.GetOrder(ctx, input.OrderID, input.AuthHeader)
	if err != nil {
		return model.Payment{}, fmt.Errorf("get order client: %w", err)
	}
	payment := repository.ProcessPaymentParams{
		OrderID:     input.OrderID,
		Status:      model.PaymentCompleted,
		UserID:      input.UserID,
		AmountCents: order.TotalAmountCents,
	}
	savedPayment, err := s.paymentRep.Save(ctx, payment)
	if err != nil {
		return model.Payment{}, fmt.Errorf("save payment: %w", err)
	}
	savedPaymentJson, err := json.Marshal(savedPayment)
	if err != nil {
		return model.Payment{}, fmt.Errorf("marshal payment completed event: %w", err)
	}
	if err = s.producer.Publish("payment.completed", savedPaymentJson); err != nil {
		return model.Payment{}, fmt.Errorf("publish payment completed event: %w", err)
	}
	return savedPayment, nil
}

func (s *PaymentService) GetAllPayments(ctx context.Context, limit, offset int) ([]model.Payment, error) {
	return s.paymentRep.GetAll(ctx, limit, offset)
}

func (s *PaymentService) GetAllByUserId(ctx context.Context, input GetAllByUserIDInput) ([]model.Payment, error) {
	if input.ID <= 0 {
		return []model.Payment{}, ErrInvalidInput
	}
	params := repository.GetAllByUserIDParams{
		ID:     input.ID,
		Limit:  input.Limit,
		Offset: input.Offset,
	}
	return s.paymentRep.GetAllByUserId(ctx, params)
}

func (s *PaymentService) GetPaymentByID(ctx context.Context, id int) (model.Payment, error) {
	if id <= 0 {
		return model.Payment{}, ErrInvalidInput
	}
	return s.paymentRep.GetByID(ctx, id)
}

func (s *PaymentService) UpdatePayment(ctx context.Context, input UpdatePaymentInput) (model.Payment, error) {
	if input.ID < 0 {
		return model.Payment{}, ErrInvalidInput
	}
	if input.Status != model.PaymentPending && input.Status != model.PaymentCompleted && input.Status != model.PaymentFailed {
		return model.Payment{}, fmt.Errorf("incorrect status: %w", ErrInvalidInput)
	}
	params := repository.UpdatePaymentParams{
		ID:     input.ID,
		Status: input.Status,
	}
	return s.paymentRep.UpdateStatus(ctx, params)
}

func (s *PaymentService) DeletePayment(ctx context.Context, id int) error {
	if id <= 0 {
		return ErrInvalidInput
	}
	return s.paymentRep.Delete(ctx, id)
}
