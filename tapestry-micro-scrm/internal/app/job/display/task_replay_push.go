package display

import (
	"fmt"
	"micro-scrm/internal/app"
	display2 "micro-scrm/internal/app/dao/display"
	"micro-scrm/internal/app/http/request"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	messageSvc "micro-scrm/internal/app/service/message"
	"micro-scrm/internal/orm"
	"micro-scrm/internal/orm/model/boss"
	"micro-scrm/internal/pkg/helper"
	qyWork "micro-scrm/internal/pkg/sdk/qy_work"
	"micro-scrm/internal/pkg/xjob"
	"time"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
)

type TaskReplayPush struct {
	Name string
}

func NewTaskReplayPush() TaskReplayPush {
	return TaskReplayPush{
		Name: "VM任务special在有效期5天后未审核需转移至VM审核",
	}
}

func (job TaskReplayPush) Run(c ctx.Context, handler *xjob.JobHandler) error {

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

func (job TaskReplayPush) handle(c ctx.Context, enterpriseId int64) error {
	err := job.delSpecialUserId(c, enterpriseId)
	if err != nil {
		fmt.Printf("执行更新企业%d的VM陈列任务删除过期未审核的special人员%v\n", enterpriseId, err)
	}
	return nil
}

func (job TaskReplayPush) delSpecialUserId(c ctx.Context, enterpriseId int64) error {
	// 查询所有结束时间小于当前时间，状态为进行中的任务，更新任务状态为已结束
	now := time.Now()
	todayZero := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	// 获取5天前的0点
	sixDaysAgoZero := todayZero.AddDate(0, 0, -6)
	fiveDaysAgoZero := todayZero.AddDate(0, 0, -5)
	// 转换为时间戳（秒级）
	sixDaysAgoZeroTimestamp := sixDaysAgoZero.Unix()
	fiveDaysAgoZeroTimestamp := fiveDaysAgoZero.Unix()
	many, err := display2.NewTaskDao().GetMany(c, enterpriseId, "end_time>=? and end_time<? and deleted_at=0", sixDaysAgoZeroTimestamp, fiveDaysAgoZeroTimestamp)
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
	allotList, err := display2.NewAllotDao().GetMany(c, enterpriseId, "display_task_id in ? and finish_status=? and special_user_id <> ? and special_check <=? and deleted_at=0", taskIds, microScrmModel.VMFinishStatusCheck, "", microScrmModel.VMSpecialCheckWait)
	if err != nil {
		return err
	}
	deleteIds := make([]int64, 0)
	for _, model := range allotList {
		deleteIds = append(deleteIds, model.ID)
	}
	err = display2.NewAllotDao().UpdateData(c, enterpriseId, map[string]interface{}{
		"id": deleteIds,
	}, map[string]interface{}{
		"special_user_id": "",
		"special_remark":  "",
	})

	for _, allotInfo := range allotList {
		// VM发送审核通知
		uri := fmt.Sprintf("display-task-audit?id=%d&store_id=%d", allotInfo.DisplayTaskID, allotInfo.StoreID)
		content := make([]map[string]string, 0)
		content = append(content, map[string]string{
			"keyname": "陈列任务名称",
			"value":   allotInfo.DisplayName,
		})
		content = append(content, map[string]string{
			"keyname": "提交门店",
			"value":   allotInfo.StoreName,
		})
		submitTime := helper.GetFormatDateByParams(allotInfo.SubmitTime)
		content = append(content, map[string]string{
			"keyname": "提交时间",
			"value":   submitTime,
		})
		sendUserId := allotInfo.OwnerUserID
		//go func() {
		_, err := messageSvc.NewMessageService().SendH5TemplateCardMessage(c, request.SendCardMessageToH5{
			EnterpriseBaseReq: request.EnterpriseBaseReq{
				EnterpriseId: enterpriseId,
			},
			Title:   "VM陈列任务待审核提醒",
			Uri:     uri,
			Party:   nil,
			UserIds: []string{sendUserId},
			Content: content,
			AppSign: qyWork.VM,
		})
		if err != nil {
			app.Logger().Error(c).Msgf("VM审核通知消息发送失败：%+v", err)
		}
		//}()
	}

	return err
}
