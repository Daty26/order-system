package kafka

import (
	"github.com/Daty26/order-system/payment-service/internal/service"
	"github.com/IBM/sarama"
)

type KafkaConsumer struct {
	consumer sarama.ConsumerGroup
	service  *service.PaymentService
}
