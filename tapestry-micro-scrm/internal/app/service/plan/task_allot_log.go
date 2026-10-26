package planService

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/app"
	planDao "micro-scrm/internal/app/dao/micro_scrm"
	planModels "micro-scrm/internal/app/models/micro_scrm"
)

type taskAllotLog struct {
	taskAllotLogDao *planDao.PlanTaskAllotLogDao
}

func NewTaskAllotLog() *taskAllotLog {
	return &taskAllotLog{
		taskAllotLogDao: planDao.NewPlanTaskAllotLogDao(),
	}
}

func (p *taskAllotLog) CreateTaskAllotLogData(c app.Ctx, enterpriseID int64, info planModels.TaskAllotLogModel) (id int64, err error) {
	return p.taskAllotLogDao.CreateTaskAllotLogData(c, enterpriseID, info)

}

// GetListsByTaskId 根据子任务id获取群发的导购信息
func (p *taskAllotLog) GetListsByTaskId(c ctx.Context, enterpriseId, taskId int64) ([]planModels.TaskAllotLogModel, error) {
	return p.taskAllotLogDao.GetListsByTaskId(c, enterpriseId, taskId)
}
