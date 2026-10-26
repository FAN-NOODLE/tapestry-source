package reserveJob

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"micro-scrm/internal/app/dao/micro_scrm"
	"micro-scrm/internal/orm"
	"micro-scrm/internal/orm/model/boss"
	"micro-scrm/internal/pkg/xjob"
	"sync"
	"time"
)

type UpdateReserveServiceStatus struct {
	name string
}

func NewUpdateReserveServiceStatus() *UpdateReserveServiceStatus {
	return &UpdateReserveServiceStatus{
		name: "更新服务配置的有效期状态",
	}
}

func (job *UpdateReserveServiceStatus) Run(ctx ctx.Context, handler *xjob.JobHandler) error {

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

func (job *UpdateReserveServiceStatus) getTask(c ctx.Context, enterpriseId int64) error {
	nowTime := time.Now().Unix()
	err := microScrmDao.NewReserveServiceDao().UpdateServeStatus(c, enterpriseId, nowTime)
	if err != nil {
		log.Error(c).Msgf("更新服务配置状态,异常：%v", err)
	}
	return nil
}
