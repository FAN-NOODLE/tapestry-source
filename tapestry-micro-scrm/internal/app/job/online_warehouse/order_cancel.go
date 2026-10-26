package online_warehouse

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/threading"
	"micro-scrm/internal/app"
	onlineWarehouseDao "micro-scrm/internal/app/dao/online_warehouse"
	appJob "micro-scrm/internal/app/job"
	onlineWarehouseModel "micro-scrm/internal/app/models/micro_scrm/online_warehouse"
	"micro-scrm/internal/pkg/xjob"
)

type OrderCancel struct {
	name string
}

func NewOrderCancel() *OrderCancel {
	return &OrderCancel{
		name: "网仓订单作废",
	}
}

func (job *OrderCancel) Run(c ctx.Context, handler *xjob.JobHandler) error {
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

func (job *OrderCancel) Handle(c ctx.Context, enterpriseId int64) error {

	err := onlineWarehouseDao.NewOrderDao().UpdateData(c, enterpriseId, map[string]interface{}{
		"status": onlineWarehouseModel.OrderStatusWait,
	}, map[string]interface{}{
		"status": onlineWarehouseModel.OrderStatusCancel,
	})
	if err != nil {
		return err
	}
	err = onlineWarehouseDao.NewOrderSkuDao().UpdateData(c, enterpriseId, map[string]interface{}{
		"status": onlineWarehouseModel.OrderStatusWait,
	}, map[string]interface{}{
		"status": onlineWarehouseModel.OrderStatusCancel,
	})

	return err

}
