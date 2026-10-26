package dashboardJob

import (
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/threading"
	"micro-scrm/internal/app"
	planDao "micro-scrm/internal/app/dao/micro_scrm"
	appJob "micro-scrm/internal/app/job"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/xjob"
	"time"
)

type TaskContactData struct {
	name      string
	startTime int64
	endTime   int64
}

func NewTaskContactData() *TaskContactData {
	var startTime, endTime int64 = helper.GetBetweenByTime(1)
	return &TaskContactData{
		name:      "contacted数据更新",
		startTime: startTime,
		endTime:   endTime,
	}
}

func (job *TaskContactData) Run(c ctx.Context, handler *xjob.JobHandler) error {
	enterpriseIds, err := appJob.GetEnterpriseIds(handler)
	if err != nil {
		return err
	}
	var wg threading.RoutineGroup
	for _, id := range enterpriseIds {
		tmpId := id
		wg.Run(func() {
			err = job.Handle(c, tmpId)
			if err != nil {
				app.Logger().Error(c).Msgf("error : %v", err)
			}
		})
	}
	wg.Wait()
	return nil
}

// Handle 同步dashboard数据
func (job *TaskContactData) Handle(c ctx.Context, enterpriseId int64) error {
	// 获取结束时间大于昨天的0点时间，查询到的主任务为进行中的任务，统计已联系的客户，同步数据
	many, err := planDao.NewPlanTaskDao().GetMany(c, enterpriseId, "task_end_time>=? and deleted_at=0", job.startTime)
	if err != nil {
		return err
	}
	updatedAt := time.Now().Unix()
	for _, info := range many {
		model := info
		planId := model.PlanID
		taskId := model.ID
		if model.TaskCategory != microScrmModel.PLANTASKMOMENT { // 沟通任务
			err := planDao.NewPlanTaskAllotLogDao().ExecRowSql(c, enterpriseId, updatedAt, planId, taskId)
			if err != nil {
				app.Logger().Error(c).Msgf("task_allot_log表触达好友数失败:%v", err)
				continue
			}
		} else {
			err := planDao.NewPlanTaskAllotLogDao().ExecMomentRowSql(c, enterpriseId, updatedAt, planId, taskId)
			if err != nil {
				app.Logger().Error(c).Msgf("task_allot_log表更新朋友圈失败:%v", err)
				continue
			}
		}
	}
	fmt.Println("contacted数据更新完毕")
	return err
}
