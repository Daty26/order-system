package service

import (
	"errors"
	"github.com/Daty26/order-system/payment-service/internal/model"
	"github.com/Daty26/order-system/payment-service/internal/repository"
)

type PaymentService struct {
	paymentRep repository.PaymentRep
}

func NewPaymentService(payRep repository.PaymentRep) *PaymentService {
	return &PaymentService{paymentRep: payRep}
}
func (s *PaymentService) ProcessPayment(orderId int, amount float64) (model.Payment, error) {
	if amount <= 0 {
		return model.Payment{}, errors.New("amount can't be negative")
	}
	payment := model.Payment{
		OrderID: orderId,
		Status:  model.PaymentCompleted,
		Amount:  amount,
	}
	return s.paymentRep.Save(payment)
}

func (s *PaymentService) GetAllPayments() ([]model.Payment, error) {
	return s.paymentRep.GetAll()
}
