package kafka

import (
	"encoding/json"
	"fmt"
	"github.com/Daty26/order-system/inventory-service/internal/service"
	"github.com/IBM/sarama"
)

type OrderCreated struct {
	OrderID int `json:"order_id"`
	//Items   []struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
	//} `json:'items"`
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
			fmt.Println("coulnd;t parse the message: " + err.Error())
			continue
		}

		err = kc.service.ReduceStock(event.ProductID, event.Quantity)
		if err != nil {
			fmt.Printf("error reducing stock: %s", err.Error())
		}
		fmt.Printf("product_id: %d, quantity: %d", event.ProductID, event.Quantity)
	}
	return nil
}
