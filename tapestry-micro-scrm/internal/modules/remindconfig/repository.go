package remindconfig

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
)

//Repository 数据库操作
type Repository interface {
	repos.IRepository
	GetRemindConfig(c ctx.Context, request ConfigInfoRequest) (*ConfigInfoResponse, error)
}
