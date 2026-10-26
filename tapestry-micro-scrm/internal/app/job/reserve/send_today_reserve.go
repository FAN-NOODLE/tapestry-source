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
	"time"

	"micro-scrm/internal/pkg/xjob"
)

type sendTodayReserveData struct {
	name string
}

func NewSendTodayReserveData() *sendTodayReserveData {
	return &sendTodayReserveData{
		name: "推送当天的邀约提醒",
	}
}

func (job sendTodayReserveData) Run(ctx ctx.Context, handler *xjob.JobHandler) error {
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

func (job sendTodayReserveData) getTask(ctx ctx.Context, enterpriseId int64) error {

	// IsCoach := false
	// res, err := microBossDao.NewEnterpriseConfig().GetInfo(ctx, map[string]interface{}{"enterprise_id": enterpriseId, "key": "store_name"})
	// if err != nil {
	//	return err
	// }
	// if res.Name == "coach" {
	//	IsCoach = true
	// }
	// startTime, endTime := helper.PeriodOfTimestamp()
	// if IsCoach {
	// 每天9:45提醒导购明日的预约单（邀约单上的到店时间是明天的，且邀约单状态为进行中） sw&&ks&&coach调整为一致逻辑
	startTime, endTime := helper.YesterdayPeriodOfTimestamp(-1)
	// }
	result, err := reserveSvc.NewOrderService().GetTodayOrderList(ctx, enterpriseId, startTime, endTime)
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
	todayStr := helper.GetNowDate()

	title := "明日客户到店邀约单提醒"
	remark := "您明日有需要邀约到店的客户,点击可提前查看详情哦~"
	tomorrowTime, _ := helper.YesterdayPeriodOfTimestamp(-1)
	todayStr = time.Unix(tomorrowTime, 0).Format("2006-01-02")

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
		content = append(content, map[string]string{"key": "预约日期", "value": todayStr})
		content = append(content, map[string]string{"key": "预约门店", "value": departmentName})
		content = append(content, map[string]string{"key": "温馨提示", "value": remark})
		params := microwork.SendMiniMessageRequest{
			EnterpriseID: enterpriseId,
			MemUserIds:   []string{userid},
			Title:        title,
			Content:      content,
			Page:         fmt.Sprintf("/pages/store-invite/record-list/record-list?time=%s&statusIndex=%d&member_id=%d", todayStr, microScrmModel.RUNNING, memberIds[userid]),
		}
		_, err := messageSvc.NewMessageService().SendMiniMessage(ctx, params)
		app.Logger().Info(ctx).Err(err).Interface("params", params).Msg("sendPlanTaskMiniMessage.SendMiniMessageRequest")
	}
	return nil
}
