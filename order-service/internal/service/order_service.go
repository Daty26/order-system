package service

import (
	"context"
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

func (s *OrderService) CreateOrder(ctx context.Context, actor Actor, input CreatedOrderInput) (model.Orders, error) {
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
		Status:           model.OrderPending,
		TotalAmountCents: totalAmountCents,
		Items:            orderItems,
	}

	createdOrder, err := s.repo.Create(ctx, order)
	if err != nil {
		return model.Orders{}, fmt.Errorf("create order: %w", err)
	}
	if !actor.IsAdmin() && actor.UserID != createdOrder.UserID {
		return model.Orders{}, ErrForbiddenOrder
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

func (s *OrderService) GetOrders(ctx context.Context, actor Actor, limit, offset int) ([]model.Orders, error) {
	if !actor.IsAdmin() {
		return s.repo.GetAllByUserID(ctx, actor.UserID, limit, offset)
	}
	return s.repo.GetAll(ctx, limit, offset)
}

func (s *OrderService) GetOrdersByUserId(ctx context.Context, userID, limit, offset int) ([]model.Orders, error) {
	if userID <= 0 {
		return []model.Orders{}, ErrInvalidOrder
	}
	return s.repo.GetAllByUserID(ctx, userID, limit, offset)
}

func (s *OrderService) GetOrderByID(ctx context.Context, actor Actor, id int) (model.Orders, error) {
	if id <= 0 || actor.UserID <= 0 {
		return model.Orders{}, ErrInvalidOrder
	}
	order, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return model.Orders{}, fmt.Errorf("get order by id: %w", err)
	}
	if !actor.IsAdmin() && order.UserID != actor.UserID {
		return model.Orders{}, ErrForbiddenOrder
	}
	return order, nil
}

// func (s *OrderService) UpdateOrder(ctx context.Context, order model.Orders) (model.Orders, error) {
// 	for _, item := range order.Items {
// 		if item.ProductID <= 0 || item.Quantity <= 0 {
// 			return model.Orders{}, errors.New("invalid order data")
// 		}
// 	}
// 	order, err := s.repo.Update(ctx, order)
// 	if err != nil {
// 		return model.Orders{}, err
// 	}
// 	return order, nil
// }

func (s *OrderService) DeleteOrder(ctx context.Context, actor Actor, id int) error {
	if !actor.IsAdmin() {
		return ErrForbiddenOrder
	}
	if id <= 0 || actor.UserID <= 0 {
		return ErrInvalidOrder
	}
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get order before delete: %w", err)
	}
	return s.repo.Delete(ctx, id)
}

func (s *OrderService) CancelOrder(ctx context.Context, actor Actor, orderID int) (model.Orders, error) {
	if orderID <= 0 || actor.UserID <= 0 {
		return model.Orders{}, ErrInvalidOrder
	}
	order, err := s.repo.GetByID(ctx, orderID)
	if err != nil {
		return model.Orders{}, fmt.Errorf("get order: %w", err)
	}
	if !actor.IsAdmin() && order.UserID != actor.UserID {
		return model.Orders{}, ErrForbiddenOrder
	}
	if order.Status != model.OrderPending {
		return model.Orders{}, ErrOrderCannotBeCanceled
	}
	cancelledOrder, err := s.repo.Cancel(ctx, orderID)
	if err != nil {
		return model.Orders{}, fmt.Errorf("cancel order: %w", err)
	}
	cancelledOrder.Items = order.Items
	return cancelledOrder, nil
}
