package pkgKafka

import (
	"github.com/Shopify/sarama"
)

type AppKafka struct {
	producer sarama.SyncProducer
}

// NewAppKafka 沟通方法
func NewAppKafka(producer sarama.SyncProducer) *AppKafka {
	return &AppKafka{producer: producer}
}

func (ak *AppKafka) Producer(data string, topic string) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(topic),
		Value: sarama.ByteEncoder(data),
	}

	_, _, err := ak.producer.SendMessage(msg)
	return err
}

func (ak *AppKafka) Close() error {
	return ak.producer.Close()
}
