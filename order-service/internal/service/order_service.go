package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Daty26/order-system/order-service/internal/kafka"
	"github.com/Daty26/order-system/order-service/internal/model"
	"github.com/Daty26/order-system/order-service/internal/repository"
	"log"
	"time"
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

func (s *OrderService) CreateOrder(order model.Orders) (model.Orders, error) {
	for _, item := range order.Items {
		if item.ProductID <= 0 || item.Quantity <= 0 {
			return model.Orders{}, errors.New("invalid order data")
		}
	}
	order.Status = "CREATED"
	order.CreatedAt = time.Now()
	createdOrder, err := s.repo.Create(order)
	if err != nil {
		return model.Orders{}, err
	}
	fmt.Println("created order:")
	fmt.Println(createdOrder)
	items := make([]map[string]int, 0)
	for _, item := range createdOrder.Items {
		items = append(items, map[string]int{
			"product_id": item.ProductID,
			"quantity":   item.Quantity,
		})
	}
	event := map[string]interface{}{
		"order_id": createdOrder.ID,
		"user_id":  createdOrder.UserID,
		"status":   createdOrder.Status,
		"items":    items,
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
		return model.Orders{}, err
	}
	return createdOrder, nil

}

func (s *OrderService) GetOrders() ([]model.Orders, error) {
	return s.repo.GetAll()
}
func (s *OrderService) GetOrdersByUserId(userId int) ([]model.Orders, error) {
	if userId < 0 {
		return []model.Orders{}, errors.New("incorrect user id")
	}
	return s.repo.GetAllByUserID(userId)
}
func (s *OrderService) GetOrderByID(id int) (model.Orders, error) {
	return s.repo.GetByID(id)
}
func (s *OrderService) UpdateOrder(order model.Orders) (model.Orders, error) {
	for _, item := range order.Items {
		if item.ProductID <= 0 || item.Quantity <= 0 {
			return model.Orders{}, errors.New("invalid order data")
		}
	}
	order, err := s.repo.Update(order)
	if err != nil {
		return model.Orders{}, err
	}
	return order, nil
}
func (s *OrderService) DeleteOrder(id int) error {
	if id <= 0 {
		return errors.New("invalid id")
	}
	return s.repo.Delete(id)
}
