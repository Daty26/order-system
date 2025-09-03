package kafka

import "github.com/IBM/sarama"

type KafkaProducer struct {
	producer sarama.SyncProducer
}

func NewKafkaProducer(broker []string) (*KafkaProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	prod, err := sarama.NewSyncProducer(broker, config)
	if err != nil {
		return nil, err
	}

	return &KafkaProducer{producer: prod}, nil
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
