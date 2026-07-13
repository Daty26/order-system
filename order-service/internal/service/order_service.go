package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/Daty26/order-system/order-service/internal/kafka"
	"github.com/Daty26/order-system/order-service/internal/model"
	"github.com/Daty26/order-system/order-service/internal/repository"
)

type OrderService struct {
	repo      *repository.PostgresOrderRepo
	inventory ProductInventory
	kafka     *kafka.KafkaProducer
}

func NewOrderService(repo *repository.PostgresOrderRepo, producer *kafka.KafkaProducer, inventory ProductInventory) *OrderService {
	return &OrderService{
		repo:      repo,
		kafka:     producer,
		inventory: inventory,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, input CreatedOrderInput) (model.Orders, error) {

	if input.UserID <= 0 || len(input.Items) == 0 || len(input.Items) > 100 {
		return model.Orders{}, ErrInvalidOrder
	}
	productIDs := make([]int, 0, len(input.Items))

	for _, item := range input.Items {
		if item.ProductID <= 0 || item.Quantity <= 0 {
			return model.Orders{}, ErrInvalidOrder
		}
		productIDs = append(productIDs, item.ProductID)
	}
	quotes, err := s.inventory.GetQuotes(ctx, productIDs)
	if err != nil {
		return model.Orders{}, fmt.Errorf("get product prices: %w", err)
	}
	orderItems := make([]model.OrderItem, 0, len(input.Items))
	var totalAmountCents int64

	for _, item := range input.Items {
		quote, exists := quotes[item.ProductID]
		if !exists {
			return model.Orders{}, fmt.Errorf("product %d: %w", item.ProductID, ErrProductNotFound)
		}

		orderItems = append(orderItems, model.OrderItem{
			ProductID:      item.ProductID,
			Quantity:       item.Quantity,
			UnitPriceCents: quote.PriceCents,
		})
		totalAmountCents += int64(item.Quantity) * quote.PriceCents
	}
	order := model.Orders{
		UserID:           input.UserID,
		Status:           model.OrderCreated,
		TotalAmountCents: totalAmountCents,
		Items:            orderItems,
	}

	createdOrder, err := s.repo.Create(ctx, order)
	if err != nil {
		return model.Orders{}, fmt.Errorf("create order: %w", err)
	}
	return createdOrder, nil
	// fmt.Println("created order:" + createdOrder)
	// items := make([]map[string]interface{}, 0)
	// totalAmount := 0.0
	// for _, item := range createdOrder.Items {
	// 	//def price 0
	// 	totalAmount += float64(item.Quantity) * item.Price
	// 	items = append(items, map[string]interface{}{
	// 		"product_id": item.ProductID,
	// 		"quantity":   item.Quantity,
	// 		"price":      item.Price,
	// 	})
	// }
	// event := map[string]interface{}{
	// 	"order_id":     createdOrder.ID,
	// 	"user_id":      createdOrder.UserID,
	// 	"status":       createdOrder.Status,
	// 	"total_amount": totalAmount,
	// 	"items":        items,
	// }
	// createdOrderJson, err := json.Marshal(event)
	// if err != nil {
	// 	return createdOrder, err
	// }
	// log.Println("topic published")
	// log.Println(string(createdOrderJson))
	// err = s.kafka.Publish("order.created", createdOrderJson)
	// if err != nil {
	// 	log.Println("failed to publish topic order.created" + err.Error())
	// 	return model.Orders{}, err
	// }
	// return createdOrder, nil

}

func (s *OrderService) GetOrders(ctx context.Context, limit, offset int) ([]model.Orders, error) {
	if limit <= 0 {
		limit = 20
	}

	if limit > 100 {
		limit = 100
	}

	if offset < 0 {
		return nil, ErrInvalidOrder
	}
	return s.repo.GetAll(ctx, limit, offset)
}

func (s *OrderService) GetOrdersByUserId(ctx context.Context, userId, limit, offset int) ([]model.Orders, error) {
	if limit <= 0 {
		limit = 20
	}

	if limit > 100 {
		limit = 100
	}

	if offset < 0 {
		return nil, ErrInvalidOrder
	}
	if userId <= 0 {
		return nil, ErrInvalidOrder
	}
	return s.repo.GetAllByUserID(ctx, userId, limit, offset)
}

func (s *OrderService) GetOrderByID(ctx context.Context, id int) (model.Orders, error) {
	return s.repo.GetByID(ctx, id)
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
