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

func (s *OrderService) CreateOrder(order model.Order) (model.Order, error) {
	if order.ProductID < 0 || order.Quantity < 0 || order.UserID < 0 {
		return model.Order{}, errors.New("invalid order data")
	}

	createdOrder, err := s.repo.Create(order)
	if err != nil {
		return model.Order{}, err
	}
	event := map[string]interface{}{
		"order_id": createdOrder.ID,
		"user_id":  createdOrder.UserID,
		"items": []map[string]int{{
			"product_id": createdOrder.ProductID,
			"quantity":   createdOrder.Quantity,
		},
		},
	}
	createdOrderJson, err := json.Marshal(event)
	if err != nil {
		return createdOrder, err
	}
	log.Println("topic published")
	log.Println(string(createdOrderJson))
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
func (s *OrderService) GetOrdersByUserId(userId int) ([]model.Order, error) {
	if userId < 0 {
		return []model.Order{}, errors.New("incorrect user id")
	}
	return s.repo.GetAllByUserID(userId)
}
func (s *OrderService) GetOrderByID(id int) (model.Order, error) {
	return s.repo.GetByID(id)
}
func (s *OrderService) UpdateOrder(id int, productId int, quantity int, userId int) (model.Order, error) {
	if productId <= 0 || quantity <= 0 || userId < 0 {
		return model.Order{}, errors.New("invalid request value")
	}
	order, err := s.repo.Update(id, productId, quantity, userId)
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
