package reserveJob

import (
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/app"
	departmentsDao "micro-scrm/internal/app/dao/departments"
	"micro-scrm/internal/app/models/micro_scrm"
	messageSvc "micro-scrm/internal/app/service/message"
	reserveSvc "micro-scrm/internal/app/service/reserve"
	"micro-scrm/internal/orm"
	"micro-scrm/internal/orm/model/boss"
	"micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/sdk/microwork"
	"strconv"
	"sync"

	"micro-scrm/internal/pkg/xjob"
)

type sendPreTodayReserveData struct {
	name string
}

func NewSendPreTodayReserveData() *sendPreTodayReserveData {
	return &sendPreTodayReserveData{
		name: "推送前一天到店的未完成邀约单",
	}
}

func (job sendPreTodayReserveData) Run(ctx ctx.Context, handler *xjob.JobHandler) error {

	var (
		err            error
		enterpriseList []boss.EnterpriseModel
	)
	enterpriseList, err = orm.BossService.Enterprise.GetAllInfo()
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	for _, enterprise := range enterpriseList {
		wg.Add(1)
		go func(enterpriseID int64) {
			defer wg.Done()
			job.getTask(ctx, enterpriseID)
		}(enterprise.ID)
	}
	wg.Wait()
	return nil

}

func (job sendPreTodayReserveData) getTask(ctx ctx.Context, enterpriseId int64) error {
	// 需要验证ks和sw，调用enterprise_config,PM要求三端企业都要有未完成提醒
	status := microScrmModel.RUNNING // 待完成的订单
	title := "昨日未完成到店邀约单提醒"
	remark := "您昨日有未确认顾客到店的邀约单,点击可查看详情哦~"
	startTime, endTime := helper.YesterdayPeriodOfTimestamp(1)
	result, err := reserveSvc.NewOrderService().GetPreTodayOrderList(ctx, enterpriseId, startTime, endTime, status)
	if err != nil {
		return err
	}
	var sendUsers []string
	memberOrders := make(map[string]int64, 0)
	memberIds := make(map[string]int64, 0)
	for _, v := range result {
		if _, ok := memberOrders[v.MemberUserID]; !ok {
			memberOrders[v.MemberUserID] = v.Num
			memberIds[v.MemberUserID] = v.MemberId
		}
		sendUsers = append(sendUsers, v.MemberUserID)
	}

	for _, userid := range sendUsers {
		if _, found := memberOrders[userid]; !found {
			continue
		}
		num := int64(0)
		if _, found := memberOrders[userid]; found {
			num = memberOrders[userid]
		}
		departmentName, _ := departmentsDao.NewDepartments().GetDepartmentsNameByMemberId(ctx, enterpriseId, memberIds[userid])
		content := make([]map[string]string, 0)
		content = append(content, map[string]string{"key": "预约服务", "value": strconv.FormatInt(num, 10) + "单"})
		content = append(content, map[string]string{"key": "预约门店", "value": departmentName})
		content = append(content, map[string]string{"key": "温馨提示", "value": remark})
		params := microwork.SendMiniMessageRequest{
			EnterpriseID: enterpriseId,
			MemUserIds:   []string{userid},
			Title:        title,
			Content:      content,
			Page:         fmt.Sprintf("/pages/store-invite/record-list/record-list?time=%s&statusIndex=%d&member_id=%d", helper.GetFormatDate(startTime, 0), status, memberIds[userid]),
		}
		_, err = messageSvc.NewMessageService().SendMiniMessage(ctx, params)
		app.Logger().Info(ctx).Err(err).Interface("params", params).Msg("sendPlanTaskMiniMessage.SendMiniMessageRequest")
	}

	return nil
}
