package customer

import (
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/threading"
	appJob "micro-scrm/internal/app/job"
	cacheSvc "micro-scrm/internal/app/service/cache"
	customerSvc "micro-scrm/internal/app/service/customer"
	"micro-scrm/internal/orm"
	"micro-scrm/internal/pkg/xjob"
	"strconv"
)

type ExternalExtraSwSync struct {
	name string
	days int
}

func NewExternalExtraSwSync() *ExternalExtraSwSync {
	return &ExternalExtraSwSync{
		name: "顾客列表任务邀约统计信息更新",
		days: 30,
	}
}

func (job *ExternalExtraSwSync) Run(c ctx.Context, handler *xjob.JobHandler) error {
	specsArgs := handler.Runner().Specs()
	enterpriseIds := make([]int64, 0)
	if len(specsArgs) > 0 {
		for _, arg := range specsArgs {
			enterpriseId, _ := strconv.ParseInt(arg, 10, 64)
			enterpriseIds = append(enterpriseIds, enterpriseId)
		}
	} else {
		enterprises, err := orm.BossService.Enterprise.GetAllInfo()
		if err != nil || len(enterprises) == 0 {
			return err
		}
		for _, enterprise := range enterprises {
			enterpriseIds = append(enterpriseIds, enterprise.ID)
		}
	}

	var wg threading.RoutineGroup

	for _, id := range enterpriseIds {
		enterpriseId := id
		cacheRes, err := cacheSvc.NewRedisData().GetEnterpriseConfig(c, enterpriseId, "store_name")
		if err != nil {
			continue
		}
		if cacheRes.Name != "sw" {
			continue
		}
		wg.Run(func() {
			job.Handle(c, enterpriseId)
		})
	}

	wg.Wait()

	return nil
}

func (job *ExternalExtraSwSync) Handle(c ctx.Context, enterpriseId int64) {
	// 用户数据处理
	err := customerSvc.NewExternalSearchSw().SetAllUserExtraData(c, enterpriseId)
	if err != nil {
		return
	}
	// 邀约数据处理
	reserveTotal, err := customerSvc.NewExternalSearchSw().SetReserveExtraData(c, enterpriseId)
	if err != nil {
		return
	}
	// 任务数据处理
	taskTotal, err := customerSvc.NewExternalSearchSw().SetTaskExtraData(c, enterpriseId, job.days)
	if err != nil {
		return
	}
	taskHaveTotal, err := customerSvc.NewExternalSearchSw().SetTaskHaveExtraData(c, enterpriseId)
	if err != nil {
		return
	}
	interactTimeTotal, err := customerSvc.NewExternalSearchSw().SetInteractTimeExtraData(c, enterpriseId)
	if err != nil {
		return
	}
	// ck数据合并
	err = customerSvc.NewExternalSearchSw().ExtraDataMerge(c, enterpriseId, false)
	if err != nil {
		return
	}
	appJob.JobNotice(c, appJob.JobData{
		EnterpriseName: "sw",
		JobName:        job.name,
		Content:        fmt.Sprintf("%s,执行成功.邀约更新:%d,任务触达:%d,有无任务:%d,活跃沟通:%d", job.name, reserveTotal, taskTotal, taskHaveTotal, interactTimeTotal),
	})

}
