package kafka

import (
	"github.com/IBM/sarama"
)

type MessageHandler func(value []byte)

type KafkaConsumer struct {
	consumer sarama.Consumer
	handler  MessageHandler
}

func NewKafkaConsumer(broker []string, handler MessageHandler) (*KafkaConsumer, error) {
	c, err := sarama.NewConsumer(broker, nil)
	if err != nil {
		return nil, err
	}
	return &KafkaConsumer{consumer: c, handler: handler}, nil
}

func (kc *KafkaConsumer) Consume(topic string) error {
	partition, err := kc.consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		return err
	}

	go func() {
		for msg := range partition.Messages() {
			kc.handler(msg.Value)
		}
	}()
	return err
}
func (kc *KafkaConsumer) Close() error {
	return kc.consumer.Close()
}
