package errs

import "gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"

var (
	MarketingLaunchZero    = errorx.NewErrorMsg(190002, "营销策略已开始不能被编辑")
	MarketingEndZero       = errorx.NewErrorMsg(190003, "营销策略已结束不能被编辑")
	NotWithinExecutionTime = errorx.NewErrorMsg(190004, "营销策略无需执行，不在策略执行执行时间范围内不需要执行")
	MarketingCannotDelete  = errorx.NewErrorMsg(190005, "已开始的营销策略不能被删除")
)
