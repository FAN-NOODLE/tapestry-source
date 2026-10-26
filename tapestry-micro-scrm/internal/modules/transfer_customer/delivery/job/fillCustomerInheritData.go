package job

import (
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	customerSvc "micro-scrm/internal/app/service/customer"
	"micro-scrm/internal/config"
	"micro-scrm/internal/orm"
	"micro-scrm/internal/orm/model/boss"
	"micro-scrm/internal/pkg/xjob"
)

// 离职待继承人员信息完善
type fillCustomerInheritData struct {
	inheritService *customerSvc.Inherit
	name           string
}

func FillCustomerInheritData() *fillCustomerInheritData {
	return &fillCustomerInheritData{
		name:           fmt.Sprintf("%s::离职待继承外部联系人数据冗余字段更新Job", config.CfgData.App),
		inheritService: customerSvc.NewInherit(),
	}
}
func (job fillCustomerInheritData) Run(ctx ctx.Context, jobHandler *xjob.JobHandler) error {

	// 获取所有品牌的企业ID
	enterpriseInfoList, enterpriseInfoErr := orm.BossService.Enterprise.GetAllInfo()
	if enterpriseInfoErr != nil {
		log.Info(ctx).Msgf("%s::获取所有企业的信息时报错，错误为：%v", job.name, enterpriseInfoErr)
		return enterpriseInfoErr
	}

	// var wg sync.WaitGroup
	for _, enterpriseInfo := range enterpriseInfoList {
		// wg.Add(1)
		job.exec(ctx, enterpriseInfo)

	}
	// wg.Wait()

	fmt.Printf("%s,任务执行结束\n", job.name)
	return nil
}

func (job fillCustomerInheritData) exec(ctx ctx.Context, data boss.EnterpriseModel) {
	if err := job.inheritService.UpdateIncompleteData(ctx, data.ID); err != nil {
		log.Info(ctx).Msgf("%s::离职待继承外部联系人数据补全报错，错误为：%v", job.name, err)
		return
	}

}
