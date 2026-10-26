package reserveJob

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	microBossDao "micro-scrm/internal/app/dao/micro_boss"
	reserveSvc "micro-scrm/internal/app/service/reserve"
	"micro-scrm/internal/orm"
	"micro-scrm/internal/orm/model/boss"
	"micro-scrm/internal/pkg/helper"
	"sync"

	"micro-scrm/internal/pkg/xjob"
)

type preTodayReserveOrderChange struct {
	name string
}

func NewPreTodayReserveOrderChange() *preTodayReserveOrderChange {
	return &preTodayReserveOrderChange{
		name: "前一天的到店时间订单状态更改",
	}
}

func (job preTodayReserveOrderChange) Run(ctx ctx.Context, handler *xjob.JobHandler) error {

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

func (job preTodayReserveOrderChange) getTask(ctx ctx.Context, enterpriseId int64) error {
	// 需要验证ks和sw，调用enterprise_config
	res, err := microBossDao.NewEnterpriseConfig().GetInfo(ctx, map[string]interface{}{"enterprise_id": enterpriseId, "key": "pre_today_order_change"})
	if err != nil {
		return err
	}
	if res.ID == 0 {
		return nil
	}
	startTime, endTime := helper.YesterdayPeriodOfTimestamp(1)
	return reserveSvc.NewOrderService().UpdatePreTodayOrders(ctx, enterpriseId, startTime, endTime)
}
