package service

import (
	"errors"
	"github.com/Daty26/order-system/order-service/internal/model"
	"github.com/Daty26/order-system/order-service/internal/repository"
)

type OrderService struct {
	repo repository.OrderRep
}

func NewOrderService(repo repository.OrderRep) *OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) CreateOrder(item string, amount int) (model.Order, error) {
	if item == "" || amount <= 0 {
		return model.Order{}, errors.New("invalid order data")
	}
	order := model.Order{
		Item:   item,
		Amount: amount,
	}
	return s.repo.Create(order)
}

func (s *OrderService) GetOrders() ([]model.Order, error) {
	return s.repo.GetAll()
}
func (s *OrderService) GetOrderByID(id int) (model.Order, error) {
	order, err := s.repo.GetByID(id)
	if err != nil {
		return model.Order{}, err
	}
	return order, nil
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
