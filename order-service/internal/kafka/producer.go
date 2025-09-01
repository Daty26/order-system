package kafka

import "github.com/IBM/sarama"

type KafkaProducer struct {
	producer sarama.SyncProducer
}

func NewKafkaProducer(brokers []string) (*KafkaProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	pd, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &KafkaProducer{producer: pd}, nil
}
func (k *KafkaProducer) Publish(topic string, message []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(message),
	}
	_, _, err := k.producer.SendMessage(msg)
	return err
}

func (k *KafkaProducer) Close() error {
	return k.producer.Close()
}
