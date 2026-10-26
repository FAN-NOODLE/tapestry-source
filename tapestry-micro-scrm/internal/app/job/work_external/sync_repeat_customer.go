package workExternalJob

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"math"
	workExternalSvc "micro-scrm/internal/app/service/work_external"
	"micro-scrm/internal/orm"
	"micro-scrm/internal/orm/model/boss"
	"sync"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"micro-scrm/internal/pkg/xjob"
)

type repeatCustomer struct {
	name string
}

func NewRepeatCustomer() *repeatCustomer {
	return &repeatCustomer{
		name: "同步重复好友数据",
	}
}

func (job repeatCustomer) Run(ctx ctx.Context, _ *xjob.JobHandler) error {
	var (
		err            error
		enterpriseList []boss.EnterpriseModel
	)
	repeatSvc := workExternalSvc.NewRepeatCustomer()
	enterpriseList, err = orm.BossService.Enterprise.GetAllInfo()
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	for _, enterprise := range enterpriseList {
		wg.Add(1)
		go func(enterpriseId int64) {
			defer wg.Done()
			err = job.handleEnterpriseRepeatCustomer(ctx, enterpriseId, repeatSvc)
			if err != nil {
				log.Error(ctx).Msgf("处理企业%d异常,错误:%v", enterpriseId, err)
			}
		}(enterprise.ID)
	}
	wg.Wait()
	return nil
}

func (job repeatCustomer) handleEnterpriseRepeatCustomer(ctx ctx.Context, enterpriseID int64, service *workExternalSvc.RepeatCustomer) error {
	total, err := service.SyncRepeatCustomerTotal(ctx, enterpriseID)
	if err != nil || total == 0 {
		return err
	}
	size := 50
	maxPage := int(math.Ceil(float64(total / int64(size))))
	for i := 0; i <= maxPage; i++ {
		lists, err := service.SyncRepeatCustomerInfo(ctx, enterpriseID, i, size)
		if err != nil || len(lists) == 0 {
			continue
		}
		if err = service.UpdateRepeatCustomer(ctx, enterpriseID, lists); err != nil {
			continue
		}
	}
	return nil
}
