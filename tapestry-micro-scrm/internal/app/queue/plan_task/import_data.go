package plan_task

import (
	"encoding/json"
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/app"
	planReq "micro-scrm/internal/app/http/request/plan"
	cacheSvc "micro-scrm/internal/app/service/cache"
	planService "micro-scrm/internal/app/service/plan"
	"strconv"
)

type ImportDataProduce struct {
	groupId string
	topics  []string
}

func NewImportDataProduce() *ImportDataProduce {
	return &ImportDataProduce{
		groupId: "plan_task",
		topics:  []string{"import_data"},
	}
}

func (q *ImportDataProduce) Topics() []string {
	return q.topics
}

func (q *ImportDataProduce) GroupID() string {
	return q.groupId
}

func (q *ImportDataProduce) Run(c ctx.Context, key string, value *string) error {
	// 处理调整人群包数据
	var params planReq.ImportDataQueueData
	err := json.Unmarshal([]byte(*value), &params)
	if err != nil {
		app.Logger().Error(c).Msgf("%s计划人群人群包导入的数据解析失败，错误原因为:%v", key, err)
		return err
	}
	// 针对于同一个子任务的同一个人群包的数据，只能有一个消费者消费
	redisKey := fmt.Sprintf("palntask:importdata:%s:%s:%s", strconv.Itoa(int(params.EnterpriseId)), strconv.Itoa(int(params.ConsumerFileId)), params.TaskName)
	lock := app.Redis().SetNX(c, redisKey, 1, 0).Val()
	defer app.Redis().Del(c, redisKey)
	if !lock {
		return err
	}
	config, err := cacheSvc.NewRedisData().GetEnterpriseConfig(c, params.EnterpriseId, "store_name")
	if err != nil {
		return err
	}
	storeName := config.Name
	planService.NewPlanTaskImport().SyncTaskImportData(c, params.EnterpriseId, params.ConsumerFileId, params.SendType, params.TaskName, params.TaskId, storeName)
	return err
}
