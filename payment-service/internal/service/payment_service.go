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

func (s *PaymentService) ProcessPayment(orderID int, amount int) (model.Payment, error) {
	if amount <= 0 {
		return model.Payment{}, errors.New("amount can't be negative")
	}
	payment := model.Payment{
		OrderID: orderID,
		Status:  model.PaymentCompleted,
		Amount:  amount,
	}
	return s.paymentRep.Save(payment)
}

func (s *PaymentService) GetAllPayments() ([]model.Payment, error) {
	return s.paymentRep.GetAll()
}

func (s *PaymentService) GetPaymentByID(id int) (model.Payment, error) {
	if id < 0 {
		return model.Payment{}, errors.New("invalid id")
	}
	return s.paymentRep.GetByID(id)
}
func (s *PaymentService) UpdatePayment(id int, status model.PaymentStatus, amount float64) (model.Payment, error) {
	if id < 0 {
		return model.Payment{}, errors.New("invalid id")
	}
	if status != model.PaymentPending && status != model.PaymentCompleted && status != model.PaymentFailed {
		return model.Payment{}, errors.New("incorrect type of status")
	}
	if amount < 0 {
		return model.Payment{}, errors.New("invalid amount")
	}
	return s.paymentRep.Update(id, status, amount)
}

func (s *PaymentService) DeletePayment(id int) error {
	if id <= 0 {
		return errors.New("invalid id")
	}
	return s.paymentRep.Delete(id)
}
