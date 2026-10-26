package excitation

import (
	"encoding/json"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	excitationReq "micro-scrm/internal/app/http/request/excitation"
	excitationSvc "micro-scrm/internal/app/service/excitation"
	"time"
)

type PointsProduce struct {
	groupId string
	topics  []string
}

func NewPointsProduceQueue() *PointsProduce {
	return &PointsProduce{
		groupId: "excitation",
		topics:  []string{"points_produce"},
	}
}

func (q *PointsProduce) Topics() []string {
	return q.topics
}

func (q *PointsProduce) GroupID() string {
	return q.groupId
}

func (q *PointsProduce) Run(c ctx.Context, key string, value *string) error {
	var params excitationReq.PointsProduceSave
	err := json.Unmarshal([]byte(*value), &params)
	if err != nil {
		log.Info(c).Msgf("%s新增积分数据累计失败,积分类型:%d,错误原因：%+v", key, params.Type, err)
		return err
	}
	nowUnix := time.Now().Unix()
	if params.AddTime > 0 {
		nowUnix = params.AddTime
	}
	err = excitationSvc.NewPointsProduceService().Produce(c, nowUnix, params)
	if err != nil {
		log.Info(c).Msgf("%s新增积分数据累计失败,积分类型:%d,错误原因：%+v", key, params.Type, err)
		return err
	}
	return nil
}
