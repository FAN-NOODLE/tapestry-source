package enterprise

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/helper"
	"gw-scrm/internal/pkg/sdk/microboss"
)

// Info 获取企业信息
func Info(ctx ctx.Context, id int64) interface{} {
	var enterpriseParams Enterprise
	enterpriseParams.EnterpriseId = id
	res, err := microboss.EnterpriseInfo(ctx, enterpriseParams)
	helper.Must(err)
	return res
}
