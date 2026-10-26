package pkgKafka

import (
	"errors"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/panjf2000/ants/v2"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/threading"
	"micro-scrm/internal/pkg/helper"
	"strings"
	"sync"
	"time"
)

// Handle 消费处理器
type Handle func(c ctx.Context, key string, value *string) error

// Group 消费组
type Group struct {
	ctx           ctx.Context
	groupID       string
	consumerGroup sarama.ConsumerGroup
	topics        []string
	handle        Handle
	// 消费异步重试
	consumeRetryMax int
	// 异步消费
	asyncConsume bool
}

// Close 关闭
func (c *Group) Close() error {
	return c.consumerGroup.Close()
}

// Consume 消费
func (c *Group) Consume() error {
	// 消费
	threading.GoSafe(func() error {
		for {
			err := c.consumerGroup.Consume(c.ctx, c.topics, c)
			if err != nil {
				c.error(c.ctx, "kafka consume error: %v", err)
				if errors.Is(err, sarama.ErrClosedConsumerGroup) {
					return err
				}
			}
		}
	})

	// 监听消费组错误
	threading.GoSafe(func() error {
		for cgErr := range c.consumerGroup.Errors() {
			c.error(c.ctx, "got kafka consume group error: %v", cgErr)
		}
		return nil
	})

	return nil
}

// Setup 创建
func (c *Group) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup 清理
func (c *Group) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim 消费
func (c *Group) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	fmt.Println("启动啦")
	limitCh := make(chan interface{}, 1024)
	threading.GoSafe(func() error {
		defer close(limitCh)
		for {
			select {
			case <-session.Context().Done():
				c.log(c.ctx, "session context done, member id: %v", session.MemberID())
				return nil
			case cm := <-claim.Messages():
				limitCh <- cm
			}
		}
	})

	cc := helper.ChannelRateLimit(limitCh, time.Second, 1000)

	if c.asyncConsume {
		pool, _ := ants.NewPool(10)
		defer pool.Release()
		wg := sync.WaitGroup{}
		for cm := range cc {
			wg.Add(1)
			_ = pool.Submit(func() {
				defer wg.Done()
				c.msgHandler(session, cm.(*sarama.ConsumerMessage))
			})
		}
		wg.Wait()
	} else {
		for cm := range cc {
			c.msgHandler(session, cm.(*sarama.ConsumerMessage))
		}
	}

	return nil
}

// msgHandler 消息处理器
func (c *Group) msgHandler(session sarama.ConsumerGroupSession, cm *sarama.ConsumerMessage) {
	context := ctx.Request(fmt.Sprintf("queue-%s-%d", string(cm.Key), helper.Rand(9999999999)))
	valueStr := string(cm.Value)
	defer func() {
		if r := recover(); r != nil {
			c.error(context, fmt.Sprintf("consume message error: %v", r))
		}
		session.MarkMessage(cm, "")
	}()
	c.log(context, "receive message: %s", helper.MarshalToString(cm))

	for i := 0; i < c.consumeRetryMax; i++ {
		err := c.handle(context, string(cm.Key), &valueStr)
		c.log(context, "try: %d/%d, consume end, msgId: %s, err: %v, value: %s", i+1, c.consumeRetryMax, cm.Key, err, valueStr)
		if err == nil {
			break
		}
	}
}

// log 记录日志
func (c *Group) log(context ctx.Context, format string, values ...interface{}) {
	format = fmt.Sprintf("[%s@%s] %s", c.groupID, strings.Join(c.topics, ","), format)
	log.Info(context).Msgf(format, values...)
}

// error 记录错误日志
func (c *Group) error(context ctx.Context, format string, values ...interface{}) {
	format = fmt.Sprintf("[%s@%s] %s", c.groupID, strings.Join(c.topics, ","), format)
	log.Error(context).Msgf(format, values...)
}

// NewConsumerGroup 创建新消费组
func NewConsumerGroup(c ctx.Context,
	addrs []string,
	config *Config,
	groupID string,
	topics []string,
	handle Handle,
	consumeRetryMax int,
	asyncConsume bool,
) (*Group, error) {
	consumerGroup, err := sarama.NewConsumerGroup(addrs, groupID, saramaConfig(config))
	if err != nil {
		return nil, err
	}
	return &Group{
		ctx:             c,
		consumerGroup:   consumerGroup,
		groupID:         groupID,
		topics:          topics,
		handle:          handle,
		consumeRetryMax: consumeRetryMax,
		asyncConsume:    asyncConsume,
	}, nil
}
