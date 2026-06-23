package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Daty26/order-system/inventory-service/internal/service"
	"github.com/IBM/sarama"
)

type OrderCreated struct {
	OrderID int `json:"order_id"`
	UserID  int `json:"user_id"`
	Items   []struct {
		ProductID int `json:"product_id"`
		Quantity  int `json:"quantity"`
	} `json:"items"`
}

type KafkaConsumer struct {
	Consumer sarama.Consumer
	service  *service.InventoryService
}

func NewKafkaConsumer(broker []string, inventoryService *service.InventoryService) (*KafkaConsumer, error) {
	consumer, err := sarama.NewConsumer(broker, nil)
	if err != nil {
		return nil, err
	}
	return &KafkaConsumer{Consumer: consumer, service: inventoryService}, nil
}

func (kc *KafkaConsumer) Consume(ctx context.Context,topic string) error {
	partition, err := kc.Consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil{
		return fmt.Errorf("consume partition: %w", err)
	}
	defer partition.Close()
	log.Printf("Listening for messages on topic %s", topic)
	for{
		select{
		case <- ctx.Done():
			log.Printf("stopping consumer for topic: %s", topic)
			return nil
		case msg, ok := <-partition.Messages():
			if !ok{
				return nil
			}
			msgCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			err := kc.handleOrderCreated(msgCtx, msg.Value); 
			cancel()
			if err != nil{
				log.Printf("failed to process order.created message: %v", err)
			}
		}
	}
}
func (kc * KafkaConsumer) handleOrderCreated(ctx context.Context, value []byte) error{
	var event OrderCreated
	if err := json.Unmarshal(value, &event); err != nil{
		return fmt.Errorf("decode order.created event: %w", err)
	}
	log.Printf(
		"received order_id=%d with %d items for user_id=%d",
		event.OrderID,
		len(event.Items),
		event.UserID,
	)
	if event.OrderID <= 0 || len(event.Items) == 0{
		return errors.New("invalid order.created event")
	}
	for _, item := range event.Items{
		if _, err :=kc.service.ReduceStock(ctx, item.ProductID, item.Quantity); err != nil{
			return fmt.Errorf("reduce stock for product_id=%d: %w", item.ProductID, err)
		}
	}
	return nil
}
