package service

import (
	"encoding/json"
	"errors"
	"github.com/Daty26/order-system/payment-service/internal/kafka"
	"github.com/Daty26/order-system/payment-service/internal/model"
	"github.com/Daty26/order-system/payment-service/internal/repository"
	"log"
)

type PaymentService struct {
	paymentRep repository.PaymentRep
	producer   *kafka.KafkaProducer
}

func NewPaymentService(payRep repository.PaymentRep, prod *kafka.KafkaProducer) *PaymentService {
	return &PaymentService{paymentRep: payRep, producer: prod}
}

func (s *PaymentService) ProcessPayment(orderID int, amount int, userId int) (model.Payment, error) {
	if amount <= 0 {
		return model.Payment{}, errors.New("amount can't be negative")
	}
	if userId < 0 {
		return model.Payment{}, errors.New("incorrect userID")
	}
	payment := model.Payment{
		OrderID: orderID,
		Status:  model.PaymentCompleted,
		Amount:  amount,
		UserID:  userId,
	}
	savedPayment, err := s.paymentRep.Save(payment)
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
	log.Println(string(savedPaymentJson))
	return savedPayment, nil
}

func (s *PaymentService) GetAllPayments() ([]model.Payment, error) {
	return s.paymentRep.GetAll()
}
func (s *PaymentService) GetAllByUserId(userId int) ([]model.Payment, error) {
	if userId < 0 {
		return []model.Payment{}, errors.New("incorrect userid")
	}
	return s.paymentRep.GetAllByUserId(userId)
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
