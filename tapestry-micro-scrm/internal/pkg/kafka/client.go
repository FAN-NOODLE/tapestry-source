package pkgKafka

import (
	"github.com/Shopify/sarama"
	"micro-scrm/internal/config"
)

var appKaProxy *AppKafka

type KaConfig struct {
	Addr      []string `json:"addr"`
	Topic     string   `json:"topic"`
	Partition int32    `json:"partition"`
	Version   string   `yaml:"version"`
}

func InitClient(cof config.KafkaCfg) (*AppKafka, error) {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	cfg.Producer.Timeout = 5
	cfg.Producer.Partitioner = sarama.NewHashPartitioner
	cfg.Version = sarama.V0_10_2_0
	producer, err := sarama.NewSyncProducer(cof.Addr, cfg)
	if err != nil {
		return nil, err
	}

	return NewAppKafka(producer), nil
}

// InitAppKafka 初始化kafka
func InitAppKafka(cfg config.KafkaCfg) error {
	var err error
	appKaProxy, err = InitClient(cfg)
	if err != nil {
		return err
	}
	return nil
}

func GetAppKafka() *AppKafka {
	if appKaProxy == nil {
		panic("kafka 模块没有初始化")
	}
	return appKaProxy
}
