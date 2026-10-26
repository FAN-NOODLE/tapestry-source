package display

import (
	"fmt"
	display2 "micro-scrm/internal/app/dao/display"
	"micro-scrm/internal/app/http/request"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	messageSvc "micro-scrm/internal/app/service/message"
	"micro-scrm/internal/orm"
	"micro-scrm/internal/orm/model/boss"
	helper2 "micro-scrm/internal/pkg/helper"
	qyWork "micro-scrm/internal/pkg/sdk/qy_work"
	"micro-scrm/internal/pkg/xjob"
	"strconv"
	"time"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
)

type TaskStartPush struct {
	Name string
}

func NewTaskStartPush() TaskStartPush {
	return TaskStartPush{
		Name: "VM陈列任务开始提醒",
	}
}

func (job TaskStartPush) Run(c ctx.Context, handler *xjob.JobHandler) error {

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

func (job TaskStartPush) handle(c ctx.Context, enterpriseId int64) error {
	zeroTime := helper2.GetZeroDateByParams(0)
	many, err := display2.NewTaskDao().GetMany(c, enterpriseId, "start_time=? and show_type=? and status=? and deleted_at=0", zeroTime, microScrmModel.ShowStatusUp, microScrmModel.DisplayStatusIng)
	if err != nil {
		return err
	}
	if len(many) == 0 {
		return err
	}
	taskInfoMap := make(map[int64]microScrmModel.DisplayTaskModel)
	for _, model := range many {
		taskInfoMap[model.ID] = model
	}
	title := job.Name
	for taskId, model := range taskInfoMap {
		getMany, err := display2.NewAllotDao().GetMany(c, enterpriseId, map[string]interface{}{
			"display_task_id": taskId,
			"finish_status":   microScrmModel.VMFinishStatusIng,
			"deleted_at":      0,
		})
		if err != nil {
			return err
		}
		content := make([]map[string]string, 0)
		content = append(content, map[string]string{
			"keyname": "陈列任务名称",
			"value":   model.DisplayName,
		})
		startTime := time.Unix(int64(model.StartTime), 0).Format("2006/01/02")
		endTime := time.Unix(int64(model.EndTime), 0).Format("2006/01/02")
		content = append(content, map[string]string{
			"keyname": "任务执行时间",
			"value":   startTime + " - " + endTime,
		})
		for _, info := range getMany {
			storeId := strconv.Itoa(info.StoreID)
			uri := fmt.Sprintf("display-task-detail?id=%d&store_id=%d", taskId, info.StoreID)
			_, err = messageSvc.NewMessageService().SendH5TemplateCardMessage(c, request.SendCardMessageToH5{
				EnterpriseBaseReq: request.EnterpriseBaseReq{
					EnterpriseId: enterpriseId,
				},
				Title:   title,
				Uri:     uri,
				Party:   []string{storeId},
				UserIds: nil,
				Content: content,
				AppSign: qyWork.VM,
			})
			if err != nil {
				fmt.Printf("任务%d下门店%d未开始提醒模板消息发送失败", taskId, info.StoreID)
			}
		}
	}

	return nil
}
