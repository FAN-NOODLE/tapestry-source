package pkgKafka

import (
	"context"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"micro-scrm/internal/config"
)

// DefaultConsumeRetryMax 默认消费最大次数
const DefaultConsumeRetryMax = 5

// ConsumeHandle 消费回调
type ConsumeHandle Handle

// Consumer 消费组接口
type Consumer interface {
	Close() error
	Consume() error
}

// GroupServer 消费组服务
type GroupServer struct {
	ctx       ctx.Context
	groupID   string
	topics    []string
	handle    ConsumeHandle
	consumers []Consumer
	// 消费者数量
	consumerNum int
	// 消费异步重试
	consumeRetryMax int
	// 异步消费
	asyncConsume bool
}

// ListenAndServe 启动监听消费服务
func (cgs *GroupServer) ListenAndServe() error {
	for i := 0; i < cgs.consumerNum; i++ {
		kafkaConfig := &Config{
			Version: config.CfgData.Kafka.Version,
		}
		kafkaConfig.Producer.Return.Successes = true
		cg, err := NewConsumerGroup(
			cgs.ctx,
			config.CfgData.Kafka.Addr,
			kafkaConfig,
			cgs.groupID,
			cgs.topics,
			Handle(cgs.handle),
			cgs.consumeRetryMax,
			cgs.asyncConsume,
		)
		if err != nil {
			return err
		}
		err = cg.Consume()
		if err != nil {
			return err
		}
		cgs.consumers = append(cgs.consumers, cg)
	}
	return nil
}

// Shutdown 关闭消费组
func (cgs *GroupServer) Shutdown(c context.Context) (err error) {
	for _, cm := range cgs.consumers {
		if er := cm.Close(); er != nil {
			err = er
			log.Error(ctx.Wrap(c)).Msgf("close consumer err:%v", err)
		}
	}
	return
}

// Name 获取消费组服务名称
func (cgs *GroupServer) Name() string {
	return cgs.groupID
}

// NewGroupServer 初始化消费组服务
func NewGroupServer(c ctx.Context,
	groupID string,
	handle ConsumeHandle,
	topics []string,
	consumeRetryMax int,
	consumeNum int,
	asyncConsume bool,
) *GroupServer {
	if consumeRetryMax < 1 || consumeRetryMax > 20 {
		consumeRetryMax = DefaultConsumeRetryMax
	}
	return &GroupServer{
		ctx:             c,
		groupID:         groupID,
		handle:          handle,
		topics:          topics,
		consumerNum:     consumeNum,
		consumeRetryMax: consumeRetryMax,
		asyncConsume:    asyncConsume,
	}
}
