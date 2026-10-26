package pkgKafka

import "github.com/Shopify/sarama"

// SyncProducer 生产者接口定义
type SyncProducer interface {
	SendMessage(*ProducerMessage) (partition int32, offset int64, err error)
	SendMessages([]*ProducerMessage) error
	Close() error
}

type syncProducer struct {
	producer sarama.SyncProducer
}

// SendMessages 批量发送消息
func (s *syncProducer) SendMessages(messages []*ProducerMessage) error {
	var msgs []*sarama.ProducerMessage
	for _, v := range messages {
		msgs = append(msgs, v.Transfer())
	}
	return s.producer.SendMessages(msgs)
}

// SendMessage 发送消息
func (s *syncProducer) SendMessage(message *ProducerMessage) (partition int32, offset int64, err error) {
	return s.producer.SendMessage(message.Transfer())
}

// Close 关闭生产者
func (s *syncProducer) Close() error {
	return s.producer.Close()
}

type ProducerMessage struct {
	Topic string
	Key   string
	Value string
}

// Transfer 数据格式转换
func (p *ProducerMessage) Transfer() *sarama.ProducerMessage {
	return &sarama.ProducerMessage{
		Topic: p.Topic,
		Key:   sarama.ByteEncoder(p.Key),
		Value: sarama.ByteEncoder(p.Value),
	}
}
