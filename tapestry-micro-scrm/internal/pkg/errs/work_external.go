package errs

import "gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"

var (
	// 天数必传
	DayIsMust = errorx.NewErrorMsg(130001, "参数day必传")
	// 社群不存在
	GroupChartNotFound = errorx.NewErrorMsg(130002, "社群不存在")
)
