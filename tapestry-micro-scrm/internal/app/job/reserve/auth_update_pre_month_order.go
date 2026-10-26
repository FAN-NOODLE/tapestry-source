package reserveJob

import (
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	microBossDao "micro-scrm/internal/app/dao/micro_boss"
	reserveDao "micro-scrm/internal/app/dao/micro_scrm"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/orm"
	"micro-scrm/internal/orm/model/boss"
	"micro-scrm/internal/pkg/xjob"
	"sync"
	"time"
)

type AuthUpdatePreMonthOrder struct {
	Name string
}

func NewAuthUpdatePreMonthOrder() *AuthUpdatePreMonthOrder {
	return &AuthUpdatePreMonthOrder{
		Name: "自动完成一个月前的订单",
	}
}

func (job AuthUpdatePreMonthOrder) Run(ctx ctx.Context, handler *xjob.JobHandler) error {

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

func (job AuthUpdatePreMonthOrder) getTask(ctx ctx.Context, enterpriseID int64) error {
	res, err := microBossDao.NewEnterpriseConfig().GetInfo(ctx, map[string]interface{}{"enterprise_id": enterpriseID, "key": "auth_overtime_reserve_order"})
	if err != nil {
		return err
	}
	if res.ID == 0 {
		return nil
	}
	nowTime := time.Now().Unix()
	if res.Name != "coach" {
		overTime := nowTime - int64(res.Value)*86400
		err = reserveDao.NewReservesOrder().UpdateOrdersOverTime(ctx, enterpriseID, "reserve_time <=? and status = ?", overTime, microScrmModel.RUNNING)
	} else {
		overTime := nowTime - int64(res.Value)*86400
		err = reserveDao.NewReservesOrder().UpdateOrdersOverTime(ctx, enterpriseID, "created_at<=? and status in (?)", overTime, []int{microScrmModel.RUNNING, microScrmModel.WAITFINISH})
	}
	fmt.Printf("处理%s的邀约单改为过期,报错:%v", res.Name, err)
	// p.ReleaseTimeout(1)
	return err
}
