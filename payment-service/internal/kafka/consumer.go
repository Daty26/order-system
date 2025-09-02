package kafka

import (
	"encoding/json"
	"fmt"
	"github.com/Daty26/order-system/payment-service/internal/service"
	"github.com/IBM/sarama"
)

type KafkaConsumer struct {
	consumer sarama.Consumer
	service  *service.PaymentService
}

func NewKafkaConsumer(broker []string, svc *service.PaymentService) (*KafkaConsumer, error) {
	c, err := sarama.NewConsumer(broker, nil)
	if err != nil {
		return nil, err
	}
	return &KafkaConsumer{consumer: c, service: svc}, nil
}

func (kc *KafkaConsumer) Consume(topic string) error {
	partition, err := kc.consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		return err
	}

	go func() {
		for msg := range partition.Messages() {
			fmt.Println("Received: ", string(msg.Value))
			var order struct {
				ID     int    `json:"id"`
				Item   string `json:"item"`
				Amount int    `json:"amount"`
			}
			if err := json.Unmarshal(msg.Value, &order); err != nil {
				fmt.Println("failed td parse order: ", order)
				continue
			}

			_, err := kc.service.ProcessPayment(order.ID, order.Amount)
			if err != nil {
				fmt.Println("failed to process payment:", err)
			}
		}
	}()
	return err
}
