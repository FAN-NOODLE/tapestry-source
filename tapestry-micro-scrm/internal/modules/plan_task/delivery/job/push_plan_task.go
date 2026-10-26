package job

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/modules/enterprise"
	enterpriseMysql "micro-scrm/internal/modules/enterprise/respository/mysql"
	"micro-scrm/internal/modules/members"
	memberMysql "micro-scrm/internal/modules/members/respository/mysql"
	"micro-scrm/internal/modules/plan_task"
	planTaskMysql "micro-scrm/internal/modules/plan_task/respository/mysql"
	"micro-scrm/internal/modules/ws_trace"
	"micro-scrm/internal/modules/ws_trace/respository/mysql"
	"micro-scrm/internal/pkg/xjob"
)

func StatisticsOrderData(c ctx.Context, jobHandler *xjob.JobHandler) error {
	enterpriseService := enterprise.New(enterpriseMysql.New())
	pService := plan_task.NewService(planTaskMysql.NewPlanTask())
	wsTraceService := ws_trace.NewService(mysql.NewWsTrace())
	memberservice := members.New(memberMysql.New())
	total, enterprises, err := enterpriseService.ListEnterprise(c, nil, 1, 0)
	if err != nil {
		return err
	}
	if total <= 0 {
		return err
	}
	for _, vItem := range *enterprises {
		value := vItem
		// jobHandler.SpecDo(c, strconv.FormatInt(vItem.EnterpriseID, 10), func() {
		pService.StatisticsOrderData(c, value, wsTraceService, memberservice)
		// })
	}
	return nil
}
