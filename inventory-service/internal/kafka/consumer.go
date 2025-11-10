package kafka

import (
	"encoding/json"
	"fmt"
	"github.com/Daty26/order-system/inventory-service/internal/service"
	"github.com/IBM/sarama"
)

type OrderCreated struct {
	OrderID int `json:"order_id"`
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

func (kc *KafkaConsumer) Consume(topic string) error {
	consumePartition, err := kc.Consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		return err
	}
	defer consumePartition.Close()
	fmt.Println("Listening for messages on topic; " + topic)
	for msg := range consumePartition.Messages() {
		var event OrderCreated
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			fmt.Println("couldn;t parse the message: " + err.Error())
			continue
		}
		fmt.Printf("received order_id=%d, with %d item(s)\n", event.OrderID, len(event.Items))
		for _, items := range event.Items {
			err = kc.service.ReduceStock(items.ProductID, items.Quantity)
			fmt.Printf("product_id: %d, quantity: %d\n", items.ProductID, items.Quantity)
		}
	}
	return nil
}
