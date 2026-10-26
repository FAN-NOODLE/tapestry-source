package online_warehouse

import (
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/threading"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/http/request"
	appJob "micro-scrm/internal/app/job"
	onlineWarehouseModel "micro-scrm/internal/app/models/micro_scrm/online_warehouse"
	messageSvc "micro-scrm/internal/app/service/message"
	onlineWarehouseSvc "micro-scrm/internal/app/service/online_warehouse"
	"micro-scrm/internal/pkg/helper"
	qyWork "micro-scrm/internal/pkg/sdk/qy_work"
	"micro-scrm/internal/pkg/xjob"
	"strconv"
)

type OrderWaitNotice struct {
	name string
}

func NewOrderWaitNotice() *OrderWaitNotice {
	return &OrderWaitNotice{
		name: "网仓订单待填写通知",
	}
}

func (job *OrderWaitNotice) Run(c ctx.Context, handler *xjob.JobHandler) error {
	enterpriseIds, err := appJob.GetEnterpriseIds(handler)
	if err != nil {
		return err
	}

	var wg threading.RoutineGroup

	for _, id := range enterpriseIds {
		tmpId := id
		wg.Run(func() {
			err := job.Handle(c, tmpId)
			if err != nil {
				app.Logger().Error(c).Msgf("error : %v", err)
			}
		})
	}

	wg.Wait()

	return nil
}

func (job *OrderWaitNotice) Handle(c ctx.Context, enterpriseId int64) error {

	startTime, endTime := helper.PeriodOfTimestamp()
	list, err := onlineWarehouseSvc.NewOrderSvc().GetUpdateNumByMember(c, enterpriseId, startTime, endTime)
	if err != nil {
		return err
	}
	for _, model := range list {
		content := make([]map[string]string, 0)
		content = append(content, map[string]string{
			"keyname": "订单数",
			"value":   strconv.FormatInt(model.Num, 10),
		})
		content = append(content, map[string]string{
			"keyname": "更新时间",
			"value":   helper.GetFormatDateByParams(model.UpdateAt),
		})
		uri := fmt.Sprintf("online-order-list?status=%d", onlineWarehouseModel.OrderStatusWait)
		_, err = messageSvc.NewMessageService().SendH5TemplateCardMessage(c, request.SendCardMessageToH5{
			EnterpriseBaseReq: request.EnterpriseBaseReq{
				EnterpriseId: enterpriseId,
			},
			Title:   "您有客人未提交收货信息订单",
			Uri:     uri,
			Party:   nil,
			UserIds: []string{model.UserId},
			Content: content,
			AppSign: qyWork.OnlineWareHouse,
		})
		if err != nil {
			fmt.Printf("导购%s:网仓订单消息通知失败,错误：%v", model.UserId, err)
		}
	}
	return err

}
