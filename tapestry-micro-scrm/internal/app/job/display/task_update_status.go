package display

import (
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	display2 "micro-scrm/internal/app/dao/display"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/orm"
	"micro-scrm/internal/orm/model/boss"
	"micro-scrm/internal/pkg/xjob"
	"time"
)

type TaskUpdateStatus struct {
	Name string
}

func NewTaskUpdateStatus() TaskUpdateStatus {
	return TaskUpdateStatus{
		Name: "VM任务任务结束和开始状态更新",
	}
}

func (job TaskUpdateStatus) Run(c ctx.Context, handler *xjob.JobHandler) error {

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
	return err
}

func (job TaskUpdateStatus) handle(c ctx.Context, enterpriseId int64) error {
	nowTime := time.Now().Unix()
	err := job.updateStart(c, enterpriseId, nowTime)
	if err != nil {
		fmt.Printf("执行更新企业%d的VM陈列任务开始状态失败%v\n", enterpriseId, err)
	}
	errs := job.updateEnd(c, enterpriseId, nowTime)
	if errs != nil {
		fmt.Printf("执行更新企业%d的VM陈列任务结束状态失败%v\n", enterpriseId, errs)
	}
	return nil
}

func (job TaskUpdateStatus) updateStart(c ctx.Context, enterpriseId, nowTime int64) error {
	// 查询所有开始时间大于当前时间，状态为未开始的任务，更新任务状态为进行中
	many, err := display2.NewTaskDao().GetMany(c, enterpriseId, "start_time<=? and show_type=? and status=? and deleted_at=0", nowTime, microScrmModel.ShowStatusUp, microScrmModel.DisplayStatus)
	if err != nil {
		return err
	}
	if len(many) == 0 {
		return err
	}
	var taskIds []int64
	for _, model := range many {
		taskIds = append(taskIds, model.ID)
	}
	err = display2.NewTaskDao().UpdateData(c, enterpriseId, map[string]interface{}{
		"id": taskIds,
	}, map[string]interface{}{
		"status": microScrmModel.DisplayStatusIng,
	})
	if err != nil {
		return err
	}
	return err
}

func (job TaskUpdateStatus) updateEnd(c ctx.Context, enterpriseId, nowTime int64) error {
	// 查询所有结束时间小于当前时间，状态为进行中的任务，更新任务状态为已结束
	endMany, err := display2.NewTaskDao().GetMany(c, enterpriseId, "end_time<=? and show_type=? and status=? and deleted_at=0", nowTime, microScrmModel.ShowStatusUp, microScrmModel.DisplayStatusIng)
	if err != nil {
		return err
	}
	if len(endMany) == 0 {
		return err
	}
	var endTaskIds []int64
	for _, model := range endMany {
		endTaskIds = append(endTaskIds, model.ID)
	}
	err = display2.NewTaskDao().UpdateData(c, enterpriseId, map[string]interface{}{
		"id": endTaskIds,
	}, map[string]interface{}{
		"status": microScrmModel.DisplayStatusEnd,
	})
	if err != nil {
		return err
	}
	return nil
}
