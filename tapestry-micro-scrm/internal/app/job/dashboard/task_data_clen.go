package dashboardJob

import (
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/threading"
	"micro-scrm/internal/app"
	planDao "micro-scrm/internal/app/dao/micro_scrm"
	appJob "micro-scrm/internal/app/job"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/pkg/xjob"
	"time"
)

type TaskDataClen struct {
	name      string
	startTime int64
	endTime   int64
}

func NewTaskDataClen() *TaskDataClen {
	// var startTime, endTime int64 = helper.GetBetweenByTime(1)
	startTime, endTime := int64(1656777600), int64(1714320000)
	return &TaskDataClen{
		name:      "数据清洗",
		startTime: startTime,
		endTime:   endTime,
	}
}

func (job *TaskDataClen) Run(c ctx.Context, handler *xjob.JobHandler) error {
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
func (job *TaskDataClen) Handle(c ctx.Context, enterpriseId int64) error {
	if enterpriseId != 3360610327843968 {
		return nil
	}
	updatedAt := time.Now().Unix()
	// if enterpriseId != 3363603667487424 {
	// 	return nil
	// }
	// 数据清洗非朋友圈的任务
	contactMany, err := planDao.NewPlanTaskDao().GetMany(c, enterpriseId, "task_category=? and  deleted_at=0", microScrmModel.PLANTASKCONTACT)
	if err != nil {
		return err
	}
	for _, info := range contactMany {
		model := info
		err := planDao.NewPlanTaskAllotLogDao().ExecRowSql(c, enterpriseId, updatedAt, model.PlanID, model.ID)
		if err != nil {
			app.Logger().Error(c).Msgf("获取task_allot_log表失败:%v", err)
			continue
		}

	}
	// 数据清洗朋友圈的任务
	many, err := planDao.NewPlanTaskDao().GetMany(c, enterpriseId, "task_category=? and  deleted_at=0", microScrmModel.PLANTASKMOMENT)
	if err != nil {
		return err
	}
	for _, info := range many {
		model := info
		err := planDao.NewPlanTaskAllotLogDao().ExecMomentRowSql(c, enterpriseId, updatedAt, model.PlanID, model.ID)
		if err != nil {
			app.Logger().Error(c).Msgf("获取task_allot_log表失败:%v", err)
			continue
		}

	}
	fmt.Println("数据清洗完毕")
	return err
}
