package service

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/Daty26/order-system/payment-service/internal/kafka"
	"github.com/Daty26/order-system/payment-service/internal/model"
	"github.com/Daty26/order-system/payment-service/internal/repository"
)

type PaymentService struct {
	paymentRep repository.PaymentRep
	producer   *kafka.KafkaProducer
}

func NewPaymentService(payRep repository.PaymentRep, prod *kafka.KafkaProducer) *PaymentService {
	return &PaymentService{paymentRep: payRep, producer: prod}
}

func (s *PaymentService) ProcessPayment(ctx context.Context, orderID int, amount int64, userId int) (model.Payment, error) {
	if amount <= 0 || userId < 0 {
		return model.Payment{}, ErrInvalidInput
	}
	payment := model.Payment{
		OrderID:     orderID,
		Status:      model.PaymentCompleted,
		AmountCents: amount,
		UserID:      userId,
	}
	savedPayment, err := s.paymentRep.Save(ctx, payment)
	if err != nil {
		return model.Payment{}, err
	}
	savedPaymentJson, err := json.Marshal(savedPayment)
	if err != nil {
		return model.Payment{}, err
	}
	err = s.producer.Publish("payment.completed", savedPaymentJson)
	if err != nil {
		log.Fatalln("Couldn't publish topic payment.completed: " + err.Error())
		return model.Payment{}, err
	}
	// log.Println(string(savedPaymentJson))
	return savedPayment, nil
}

func (s *PaymentService) GetAllPayments(ctx context.Context, limit, offset int) ([]model.Payment, error) {
	return s.paymentRep.GetAll(ctx, limit, offset)
}

func (s *PaymentService) GetAllByUserId(ctx context.Context, userId int) ([]model.Payment, error) {
	if userId < 0 {
		return []model.Payment{}, ErrInvalidInput
	}
	return s.paymentRep.GetAllByUserId(ctx, userId)
}

func (s *PaymentService) GetPaymentByID(ctx context.Context, id int) (model.Payment, error) {
	if id < 0 {
		return model.Payment{}, errors.New("invalid id")
	}
	return s.paymentRep.GetByID(ctx, id)
}

func (s *PaymentService) UpdatePayment(ctx context.Context, id int, status model.PaymentStatus, amount float64) (model.Payment, error) {
	if id < 0 {
		return model.Payment{}, ErrInvalidInput
	}
	if status != model.PaymentPending && status != model.PaymentCompleted && status != model.PaymentFailed {
		return model.Payment{}, errors.New("incorrect type of status")
	}
	if amount < 0 {
		return model.Payment{}, ErrInvalidInput
	}
	return s.paymentRep.Update(ctx, id, status, amount)
}

func (s *PaymentService) DeletePayment(ctx context.Context, id int) error {
	if id <= 0 {
		return ErrInvalidInput
	}
	return s.paymentRep.Delete(ctx, id)
}
