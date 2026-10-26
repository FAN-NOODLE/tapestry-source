package queue

import (
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/uid"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/xmq"
	re "micro-download/internal/pkg/redis"
	"strconv"
	"time"
)

type redisListener struct {
	redis *redis.Client
}

type RedisMessage struct {
	MsgId    string `json:"msg_id"`
	MsgBody  string `json:"msg_body"`
	PushTime int64  `json:"push_time"`
}

// ReceiveQueueBatchMsg 批量接受消息
func (r *redisListener) ReceiveQueueBatchMsg(c ctx.Context, queueName string, numOfMsg int, pollingWaitSeconds int) (string, []*xmq.Message, error) {
	var messageList []*xmq.Message
	item := &xmq.Message{}
	redisMessage := &RedisMessage{}
	// 重置数量为1 只取一条数据  因为并没有实现ack 所以多条数据如果队列崩溃  数据将丢失
	numOfMsg = 1
	for i := 0; i < numOfMsg; i++ {
		cmd := r.redis.LPop(c, queueName)
		if cmd.Val() != "" {
			if err := json.Unmarshal([]byte(cmd.Val()), redisMessage); err == nil {
				if redisMessage.MsgBody != "" {
					item.MsgID = &redisMessage.MsgId
					item.MsgBody = &redisMessage.MsgBody
					messageList = append(messageList, item)
				}
			} else {
				log.Info(c).Msgf("xmqRedisListener接受到非法结构消息(不是RedisMessage{}类型)[%s]:%s", queueName, cmd.Val())
			}
		}
	}
	// todo 等待时间实现  业务与定时等待时间两个协程  select 谁先出来return谁
	return queueName, messageList, nil
}

// DelayAckMsg 延迟确认消息
func (r *redisListener) DelayAckMsg(c ctx.Context, queueName string, receiptHandle *string) error {
	return nil
}

// ReplayMsg 重新投递消息
func (r *redisListener) ReplayMsg(c ctx.Context, queueName string, body *string) (*xmq.SendMsgResponse, error) {
	res := r.redis.LPush(c, queueName, *body)
	_, err := res.Result()
	if err != nil {
		return nil, err
	}
	code := 200
	requestId := ""
	msgId := ""
	return &xmq.SendMsgResponse{
		Code:      &code,
		Message:   body,
		RequestID: &requestId,
		MsgID:     &msgId,
	}, nil
}

func NewRedisListener(redis *redis.Client) *redisListener {
	return &redisListener{redis: redis}
}

// dispatch 队列任务
func Dispatch(ctx ctx.Context, queueName string, msg string) error {
	message := &RedisMessage{
		MsgId:    strconv.FormatInt(uid.ID(ctx), 10),
		MsgBody:  msg,
		PushTime: time.Now().Unix(),
	}
	messageStr, err := json.Marshal(message)
	if err != nil {
		return err
	}
	cmd := re.Client().LPush(ctx, queueName, string(messageStr))
	if _, err := cmd.Result(); err != nil {
		return err
	}
	return nil
}
