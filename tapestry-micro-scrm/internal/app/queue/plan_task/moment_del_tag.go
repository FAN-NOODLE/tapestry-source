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

type MomentDelTagProduce struct {
	groupId string
	topics  []string
}

func NewMomentDelTagProduce() *MomentDelTagProduce {
	return &MomentDelTagProduce{
		groupId: "plan_task",
		topics:  []string{"moment_del_tag"},
	}
}

func (q *MomentDelTagProduce) Topics() []string {
	return q.topics
}

func (q *MomentDelTagProduce) GroupID() string {
	return q.groupId
}

func (q *MomentDelTagProduce) Run(c ctx.Context, key string, value *string) error {
	// 处理调整人群包数据
	var params planReq.MomentDelTagQueueData
	err := json.Unmarshal([]byte(*value), &params)
	if err != nil {
		app.Logger().Error(c).Msgf("%s计划人群人群包朋友圈删除标签，错误原因为:%v", key, err)
		return err
	}
	// 针对于同一个子任务的同一个人群包的数据，只能有一个消费者消费
	redisKey := fmt.Sprintf("palntask:moment_del_tag:%s:%s", strconv.Itoa(int(params.EnterpriseId)), strconv.Itoa(int(params.TaskId)))
	lock := app.Redis().SetNX(c, redisKey, 1, 0).Val()
	if !lock {
		return err
	}
	err = planService.NewPlanTask().TaskMomentDelTag(c, params)
	app.Redis().Del(c, redisKey)
	return err
}
