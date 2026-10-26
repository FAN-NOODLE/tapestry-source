package job

import (
	"strconv"

	"micro-scrm/internal/modules/customer"
	userEs "micro-scrm/internal/modules/customer/respository/elastic"
	"micro-scrm/internal/modules/customer/respository/mysql"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/modules/enterprise"
	enterpriseMysql "micro-scrm/internal/modules/enterprise/respository/mysql"
	"micro-scrm/internal/pkg/xjob"
)

// 计算重复客户
func CountRepeatCustomer(c ctx.Context, jobHandler *xjob.JobHandler) error {
	enterpriseService := enterprise.New(enterpriseMysql.New())
	total, enterprises, err := enterpriseService.ListEnterprise(c, nil, 1, 0)
	if err != nil {
		return err
	}
	if total <= 0 {
		return err
	}
	s := customer.New(mysql.NewCustomer(), mysql.NewWorkUser(), userEs.NewEsRepo())
	for _, vItem := range *enterprises {
		jobHandler.SpecDo(c, strconv.FormatInt(vItem.EnterpriseID, 10), func() {
			s.RepeatCustomerCount(c, vItem.EnterpriseID)
		})
	}
	return nil
}

// 计算流失客户
func LossCustomerAddRelationTimeJob(c ctx.Context, jobHandler *xjob.JobHandler) error {
	enterpriseService := enterprise.New(enterpriseMysql.New())
	total, enterprises, err := enterpriseService.ListEnterprise(c, nil, 1, 0)
	if err != nil {
		return err
	}
	if total <= 0 {
		return err
	}
	s := customer.New(mysql.NewCustomer(), mysql.NewWorkUser(), userEs.NewEsRepo())
	for _, vItem := range *enterprises {
		jobHandler.SpecDo(c, strconv.FormatInt(vItem.EnterpriseID, 10), func() {
			s.LossCustomerAddRelationTimeJob(c, vItem.EnterpriseID, vItem.CorpID)
		})
	}
	return nil
}
