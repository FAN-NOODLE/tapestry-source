package settings

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
)

//Repository 数据库操作
type Repository interface {
	repos.IRepository
	GetSettingsInfo(c ctx.Context, request SettingsInfoRequest) (*ChatSettings, error)
	EnterpriseList(c ctx.Context) ([]ChatSettings, error)
}
