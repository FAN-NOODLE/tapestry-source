package pkgKafka

import (
	"github.com/Shopify/sarama"
	"time"
)

type Config struct {
	Producer struct {
		Timeout time.Duration
		Return  struct {
			Successes bool
		}
	}
	Version string
}

// saramaConfig 生成sarama配置
func saramaConfig(config *Config) *sarama.Config {
	sConfig := sarama.NewConfig()
	sConfig.Producer.Return.Successes = config.Producer.Return.Successes
	if config.Producer.Timeout > 0 {
		sConfig.Producer.Timeout = config.Producer.Timeout
	}
	sConfig.Version = getVersion(config.Version)
	return sConfig
}

func getVersion(version string) sarama.KafkaVersion {
	switch version {
	case "V0_8_2_0":
		return sarama.V0_8_2_0
	case "V0_8_2_1":
		return sarama.V0_8_2_1
	case "V0_8_2_2":
		return sarama.V0_8_2_2
	case "V0_9_0_0":
		return sarama.V0_9_0_0
	case "V0_9_0_1":
		return sarama.V0_9_0_1
	case "V0_10_0_0":
		return sarama.V0_10_0_0
	case "V0_10_0_1":
		return sarama.V0_10_0_1
	case "V0_10_1_0":
		return sarama.V0_10_1_0
	case "V0_10_1_1":
		return sarama.V0_10_1_1
	case "V0_10_2_0":
		return sarama.V0_10_2_0
	case "V0_10_2_1":
		return sarama.V0_10_2_1
	case "V0_11_0_0":
		return sarama.V0_11_0_0
	case "V0_11_0_1":
		return sarama.V0_11_0_1
	case "V0_11_0_2":
		return sarama.V0_11_0_2
	case "V1_0_0_0":
		return sarama.V1_0_0_0
	case "V1_1_0_0":
		return sarama.V1_1_0_0
	case "V1_1_1_0":
		return sarama.V1_1_1_0
	case "V2_0_0_0":
		return sarama.V2_0_0_0
	case "V2_0_1_0":
		return sarama.V2_0_1_0
	case "V2_1_0_0":
		return sarama.V2_1_0_0
	case "V2_2_0_0":
		return sarama.V2_2_0_0
	case "V2_3_0_0":
		return sarama.V2_3_0_0
	case "V2_4_0_0":
		return sarama.V2_4_0_0
	case "V2_5_0_0":
		return sarama.V2_5_0_0
	case "V2_6_0_0":
		return sarama.V2_6_0_0
	default:
		return sarama.V1_1_1_0
	}
}
