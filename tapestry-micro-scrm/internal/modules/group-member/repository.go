package groupmember

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
)

// Repository 存储方法
type Repository interface {
	repos.IRepository
	SelectGroupNum(c ctx.Context, param GroupParam) (int64, error)
	SelectReGroupUv(c ctx.Context, param GroupParam) (int64, error)
	SelectAddGroupUv(c ctx.Context, param GroupParam) (int64, error)
	SelectLossGroupUv(c ctx.Context, param GroupParam) (int64, error)
}
