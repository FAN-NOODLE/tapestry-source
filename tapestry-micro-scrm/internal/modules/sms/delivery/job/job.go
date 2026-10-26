package smsJob

import (
	"strconv"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/modules/enterprise"
	enterpriseMysql "micro-scrm/internal/modules/enterprise/respository/mysql"
	"micro-scrm/internal/modules/sms"
	smsMysql "micro-scrm/internal/modules/sms/respository/mysql"
	"micro-scrm/internal/pkg/xjob"
)

// SynTemplateStatusJob 同步短信签名短信模板状态
func SynTemplateStatusJob(c ctx.Context, jobHandler *xjob.JobHandler) error {

	enterpriseService := enterprise.New(enterpriseMysql.New())
	total, enterprises, err := enterpriseService.ListEnterprise(c, nil, 1, 0)
	if err != nil {
		return err
	}

	if total <= 0 {
		return err
	}

	smsService := sms.New(smsMysql.New())
	for _, vItem := range *enterprises {
		// 执行策略
		jobHandler.SpecDo(c, strconv.FormatInt(vItem.EnterpriseID, 10), func() {
			smsService.SynTemplateStatus(c, vItem.EnterpriseID)
		})
	}

	return nil
}
