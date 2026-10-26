package drainagejob

import (
	"strconv"

	"micro-scrm/internal/modules/drainage/respository/mysql"

	"micro-scrm/internal/modules/drainage"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/modules/enterprise"
	enterpriseMysql "micro-scrm/internal/modules/enterprise/respository/mysql"
	"micro-scrm/internal/pkg/xjob"
)

// 修改欢迎语的内容结构
func FixDrainageWelcomeContentJob(c ctx.Context, jobHandler *xjob.JobHandler) error {
	enterpriseService := enterprise.New(enterpriseMysql.New())
	total, enterprises, err := enterpriseService.ListEnterprise(c, nil, 1, 0)
	if err != nil {
		return err
	}
	if total <= 0 {
		return err
	}
	s := drainage.NewDrainage(mysql.NewDrainage())
	for _, vItem := range *enterprises {
		jobHandler.SpecDo(c, strconv.FormatInt(vItem.EnterpriseID, 10), func() {
			s.FixDrainageWelcomeContent(c, vItem.EnterpriseID)
		})
	}
	return nil
}
