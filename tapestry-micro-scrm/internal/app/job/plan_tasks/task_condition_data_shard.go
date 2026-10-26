package plan_tasks

import (
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	planDao "micro-scrm/internal/app/dao/micro_scrm"
	"micro-scrm/internal/orm"
	"micro-scrm/internal/orm/model/boss"
	"micro-scrm/internal/pkg/xjob"
)

type TaskConditionDataShard struct {
	Name string
}

func NewTaskConditionDataShard() TaskConditionDataShard {
	return TaskConditionDataShard{
		Name: "计划任务task_condition_data表初始化分表",
	}
}

func (job TaskConditionDataShard) Run(c ctx.Context, handler *xjob.JobHandler) error {
	var (
		err            error
		enterpriseList []boss.EnterpriseModel
	)
	enterpriseList, err = orm.BossService.Enterprise.GetAllInfo()
	if err != nil {
		return err
	}

	for _, v := range enterpriseList {
		enterpriseId := v.ID
		err = job.handle(c, enterpriseId)
		if err != nil {
			return err
		}
	}
	return nil
}

func (job TaskConditionDataShard) handle(c ctx.Context, enterpriseId int64) error {
	var err error

	planMany, err := planDao.NewPlanTaskDao().GetPlanMany(c, enterpriseId, map[string]interface{}{
		"deleted_at": 0,
	})
	if err != nil {
		return err
	}

	for _, model := range planMany {
		err := planDao.NewPlanTaskConditionDataDao().CreateShardTable(c, enterpriseId, model.ID)
		if err != nil {
			return err
		}
		fmt.Printf("task_condition_data_%d创建成功", model.ID)
		err = planDao.NewPlanTaskConditionDataDao().InsertDataToShard(c, enterpriseId, model.ID)
		if err != nil {
			return err
		}
		fmt.Printf("task_condition_data_%d数据插入成功", model.ID)
	}
	return nil
}
