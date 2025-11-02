package service

import (
	"encoding/json"
	"errors"
	"github.com/Daty26/order-system/order-service/internal/kafka"
	"github.com/Daty26/order-system/order-service/internal/model"
	"github.com/Daty26/order-system/order-service/internal/repository"
	"log"
)

type OrderService struct {
	repo  repository.OrderRep
	kafka *kafka.KafkaProducer
}

func NewOrderService(repo repository.OrderRep, producer *kafka.KafkaProducer) *OrderService {
	return &OrderService{
		repo:  repo,
		kafka: producer,
	}
}

func (s *OrderService) CreateOrder(item string, amount int) (model.Order, error) {
	if item == "" || amount <= 0 {
		return model.Order{}, errors.New("invalid order data")
	}
	order := model.Order{
		Item:   item,
		Amount: amount,
	}
	createdOrder, err := s.repo.Create(order)
	if err != nil {
		return model.Order{}, err
	}
	createdOrderJson, err := json.Marshal(createdOrder)
	if err != nil {
		return createdOrder, err
	}
	err = s.kafka.Publish("order.created", createdOrderJson)
	if err != nil {
		log.Println("failed to publish topic order.created" + err.Error())
		return model.Order{}, err
	}
	return createdOrder, nil

}

func (s *OrderService) GetOrders() ([]model.Order, error) {
	return s.repo.GetAll()
}
func (s *OrderService) GetOrderByID(id int) (model.Order, error) {
	return s.repo.GetByID(id)
}
func (s *OrderService) UpdateOrder(id int, item string, amount int) (model.Order, error) {
	if item == "" || amount <= 0 {
		return model.Order{}, errors.New("invalid request value")
	}
	order, err := s.repo.Update(id, item, amount)
	if err != nil {
		return model.Order{}, err
	}
	return order, nil

}
func (s *OrderService) DeleteOrder(id int) error {
	if id <= 0 {
		return errors.New("invalid id")
	}
	return s.repo.Delete(id)
}
