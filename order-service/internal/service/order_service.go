package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Daty26/order-system/order-service/internal/kafka"
	"github.com/Daty26/order-system/order-service/internal/model"
	"github.com/Daty26/order-system/order-service/internal/repository"
)

type OrderService struct {
	repo      repository.OrderRep
	inventory ProductInventory
	kafka     *kafka.KafkaProducer
}

func NewOrderService(repo repository.OrderRep, producer *kafka.KafkaProducer, inventory ProductInventory) *OrderService {
	return &OrderService{
		repo:      repo,
		kafka:     producer,
		inventory: inventory,
	}
}

const orderCreatedTopic = "order.created"

func (s *OrderService) CreateOrder(ctx context.Context, actor Actor, input CreatedOrderInput) (model.Order, error) {
	if input.UserID <= 0 || len(input.Items) == 0 || len(input.Items) > 100 {
		return model.Order{}, ErrInvalidRequest
	}
	productIDs := make([]int, 0, len(input.Items))

	for _, item := range input.Items {
		if item.ProductID <= 0 || item.Quantity <= 0 {
			return model.Order{}, ErrInvalidRequest
		}
		productIDs = append(productIDs, item.ProductID)
	}
	quotes, err := s.inventory.GetQuotes(ctx, productIDs)
	if err != nil {
		return model.Order{}, fmt.Errorf("get product prices: %w", err)
	}
	orderItems := make([]model.OrderItem, 0, len(input.Items))
	var totalAmountCents int64

	for _, item := range input.Items {
		quote, exists := quotes[item.ProductID]
		if !exists {
			return model.Order{}, fmt.Errorf("product %d: %w", item.ProductID, ErrProductNotFound)
		}
		orderItems = append(orderItems, model.OrderItem{
			ProductID:      item.ProductID,
			Quantity:       item.Quantity,
			UnitPriceCents: quote.PriceCents,
		})
		totalAmountCents += int64(item.Quantity) * quote.PriceCents
	}
	order := model.Order{
		UserID:           input.UserID,
		Status:           model.OrderPending,
		TotalAmountCents: totalAmountCents,
		Items:            orderItems,
	}

	createdOrder, err := s.repo.Create(ctx, order)
	if err != nil {
		return model.Order{}, fmt.Errorf("create order: %w", err)
	}

	if !actor.IsAdmin() && actor.UserID != createdOrder.UserID {
		return model.Order{}, ErrForbiddenOrder
	}

	if err := s.publishCreatedOrder(createdOrder); err != nil {
		return model.Order{}, err
	}

	return createdOrder, nil
}

func (s *OrderService) publishCreatedOrder(order model.Order) error {
	orderCreatedEvent := model.NewOrderCreatedEvent(order)

	payload, err := json.Marshal(orderCreatedEvent)
	if err != nil {
		return fmt.Errorf("marshal order created event: %w", err)
	}

	if err = s.kafka.Publish(orderCreatedTopic, payload); err != nil {
		return fmt.Errorf("kafka publish order.created: %w", err)
	}
	return nil

}

func (s *OrderService) GetOrders(ctx context.Context, actor Actor, limit, offset int) ([]model.Order, error) {
	if !actor.IsAdmin() {
		return s.repo.GetAllByUserID(ctx, actor.UserID, limit, offset)
	}
	return s.repo.GetAll(ctx, limit, offset)
}

func (s *OrderService) GetOrdersByUserId(ctx context.Context, userID, limit, offset int) ([]model.Order, error) {
	if userID <= 0 {
		return []model.Order{}, ErrInvalidRequest
	}
	return s.repo.GetAllByUserID(ctx, userID, limit, offset)
}

func (s *OrderService) GetOrderByID(ctx context.Context, actor Actor, id int) (model.Order, error) {
	if id <= 0 || actor.UserID <= 0 {
		return model.Order{}, ErrInvalidRequest
	}
	order, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return model.Order{}, fmt.Errorf("get order by id: %w", err)
	}
	if !actor.IsAdmin() && order.UserID != actor.UserID {
		return model.Order{}, ErrForbiddenOrder
	}
	return order, nil
}

// func (s *OrderService) UpdateOrder(ctx context.Context, order model.Order) (model.Order, error) {
// 	for _, item := range order.Items {
// 		if item.ProductID <= 0 || item.Quantity <= 0 {
// 			return model.Order{}, errors.New("invalid order data")
// 		}
// 	}
// 	order, err := s.repo.Update(ctx, order)
// 	if err != nil {
// 		return model.Order{}, err
// 	}
// 	return order, nil
// }

func (s *OrderService) DeleteOrder(ctx context.Context, actor Actor, id int) error {
	if !actor.IsAdmin() {
		return ErrForbiddenOrder
	}
	if id <= 0 || actor.UserID <= 0 {
		return ErrInvalidRequest
	}
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get order before delete: %w", err)
	}
	return s.repo.Delete(ctx, id)
}

func (s *OrderService) CancelOrder(ctx context.Context, actor Actor, orderID int) (model.Order, error) {
	if orderID <= 0 || actor.UserID <= 0 {
		return model.Order{}, ErrInvalidRequest
	}
	order, err := s.repo.GetByID(ctx, orderID)
	if err != nil {
		return model.Order{}, fmt.Errorf("get order: %w", err)
	}
	if !actor.IsAdmin() && order.UserID != actor.UserID {
		return model.Order{}, ErrForbiddenOrder
	}
	if order.Status != model.OrderPending {
		return model.Order{}, ErrOrderCannotBeCanceled
	}
	cancelledOrder, err := s.repo.Cancel(ctx, orderID)
	if err != nil {
		return model.Order{}, fmt.Errorf("cancel order: %w", err)
	}
	cancelledOrder.Items = order.Items
	return cancelledOrder, nil
}
