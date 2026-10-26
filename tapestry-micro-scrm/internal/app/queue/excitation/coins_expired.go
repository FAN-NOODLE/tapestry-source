package excitation

import (
	"encoding/json"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	excitationReq "micro-scrm/internal/app/http/request/excitation"
	excitationSvc "micro-scrm/internal/app/service/excitation"
)

type CoinsExpired struct {
	groupId string
	topics  []string
}

func NewCoinsExpiredQueue() *CoinsExpired {
	return &CoinsExpired{
		groupId: "excitation",
		topics:  []string{"coins_expired"},
	}
}

func (q *CoinsExpired) Topics() []string {
	return q.topics
}

func (q *CoinsExpired) GroupID() string {
	return q.groupId
}

func (q *CoinsExpired) Run(c ctx.Context, key string, value *string) error {
	var params excitationReq.CoinExpiredReq
	err := json.Unmarshal([]byte(*value), &params)
	if err != nil {
		log.Info(c).Msgf("%s金币过期队列数据解析失败,错误原因：%+v", key, err)
		return err
	}
	err = excitationSvc.NewCoinsService().CoinsExpiredConsume(c, params.EnterpriseId, params.DelaySeconds)
	if err != nil {
		log.Info(c).Msgf("%s金币过期队列执行失败,企业id:%d,错误原因：%+v", key, params.EnterpriseId, err)
		return err
	}
	return nil
}
