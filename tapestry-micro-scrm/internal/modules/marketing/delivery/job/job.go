package marketingJob

import (
	"time"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	userEs "micro-scrm/internal/modules/customer/respository/elastic"
	"micro-scrm/internal/modules/enterprise"
	enterpriseMysql "micro-scrm/internal/modules/enterprise/respository/mysql"
	"micro-scrm/internal/modules/marketing"
	marketingMysql "micro-scrm/internal/modules/marketing/respository/mysql"
	"micro-scrm/internal/modules/sms"
	smsMysql "micro-scrm/internal/modules/sms/respository/mysql"
	"micro-scrm/internal/pkg/xjob"
)

// MarketingCronJob 执行策略可执行记录
func MarketingCronJob(c ctx.Context, jobHandler *xjob.JobHandler) error {
	enterpriseService := enterprise.New(enterpriseMysql.New())
	total, enterprises, err := enterpriseService.ListEnterprise(c, nil, 1, 0)
	if err != nil {
		return err
	}

	if total <= 0 {
		return err
	}

	smsService := sms.New(smsMysql.New())
	marketingService := marketing.New(marketingMysql.New(), userEs.NewEsRepo())
	marketingService.SmsService = smsService
	// 当前时间
	currentTime := time.Now().Unix()
	for _, vItem := range *enterprises {
		// 执行策略
		marketingService.MarketingCron(c, vItem.EnterpriseID, vItem.Domain, currentTime, smsService)
	}
	return nil
}

// GetGroupMsgResultJob 获取群发发送结果
func GetGroupMsgResultJob(c ctx.Context, jobHandler *xjob.JobHandler) error {
	marketingService := marketing.New(marketingMysql.New(), userEs.NewEsRepo())

	err := marketingService.GetGroupMsgResult(c)

	return err
}

// AutoStartMarketing 自动开始营销策略
func AutoStartMarketing(c ctx.Context, jobHandler *xjob.JobHandler) error {
	enterpriseService := enterprise.New(enterpriseMysql.New())
	total, enterprises, err := enterpriseService.ListEnterprise(c, nil, 1, 0)
	if err != nil {
		return err
	}

	if total <= 0 {
		return err
	}

	marketingService := marketing.New(marketingMysql.New(), userEs.NewEsRepo())

	for _, vItem := range *enterprises {
		// 执行策略
		marketingService.DealWaitStartStrategy(c, vItem.EnterpriseID)

	}

	return nil
}

// AutoEndMarketing 自动结束营销策略
func AutoEndMarketing(c ctx.Context, jobHandler *xjob.JobHandler) error {
	enterpriseService := enterprise.New(enterpriseMysql.New())
	total, enterprises, err := enterpriseService.ListEnterprise(c, nil, 1, 0)
	if err != nil {
		return err
	}

	if total <= 0 {
		return err
	}

	marketingService := marketing.New(marketingMysql.New(), userEs.NewEsRepo())

	for _, vItem := range *enterprises {
		// 执行策略
		marketingService.DealWaitEndStrategy(c, vItem.EnterpriseID)

	}

	return nil
}

// CheckGoOnStrategyTemplateJob 检测投放中的策略模版是否正常
func CheckGoOnStrategyTemplateJob(c ctx.Context, jobHandler *xjob.JobHandler) error {
	enterpriseService := enterprise.New(enterpriseMysql.New())
	total, enterprises, err := enterpriseService.ListEnterprise(c, nil, 1, 0)
	if err != nil {
		return err
	}

	if total <= 0 {
		return err
	}

	smsService := sms.New(smsMysql.New())
	marketingService := marketing.New(marketingMysql.New(), userEs.NewEsRepo())
	marketingService.SmsService = smsService
	for _, vItem := range *enterprises {
		// 执行策略
		marketingService.CheckGoOnStrategyTemplate(c, vItem.EnterpriseID, smsService)

	}
	return nil
}
