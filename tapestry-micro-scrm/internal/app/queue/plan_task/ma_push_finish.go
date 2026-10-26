package plan_task

import (
	"encoding/json"
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/app"
	planReq "micro-scrm/internal/app/http/request/plan"
	planService "micro-scrm/internal/app/service/plan"
	"strconv"
)

type MaPushFinish struct {
	groupId string
	topics  []string
}

func NewMaPushFinishProduce() *MaPushFinish {
	return &MaPushFinish{
		groupId: "plan_task",
		topics:  []string{"ma_push_finish"},
	}
}

func (q *MaPushFinish) Topics() []string {
	return q.topics
}

func (q *MaPushFinish) GroupID() string {
	return q.groupId
}

func (q *MaPushFinish) Run(c ctx.Context, key string, value *string) error {
	// 处理调整人群包数据
	var params planReq.MaPushFinishQueueData
	err := json.Unmarshal([]byte(*value), &params)
	if err != nil {
		app.Logger().Error(c).Msgf("%s计划人群人群包朋友圈删除标签，错误原因为:%v", key, err)
		return err
	}
	// 针对于同一个子任务的同一个人群包的数据，只能有一个消费者消费
	redisKey := fmt.Sprintf("palntask:ma_push_finish:%s:%s", strconv.Itoa(int(params.EnterpriseId)), strconv.Itoa(int(params.TaskId)))
	lock := app.Redis().SetNX(c, redisKey, 1, 0).Val()
	if !lock {
		return err
	}
	err = planService.NewPlanTask().MaPushFinish(c, params)
	app.Redis().Del(c, redisKey)
	return err
}
