package scripts

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
)

// Repository 数据库操作
type Repository interface {
	repos.IRepository
	// user_category 分表SQL
	SpiltTable(cc ctx.Context, enterpriseID int64) error
	//分库sql
	SpiltDatabase(cc ctx.Context, enterpriseID int64) error
	GetRepeatData(cc ctx.Context) ([]Data, error)
	DelRepeatData(cc ctx.Context, id int64, number int64) error
	GetHasRemarkData(cc ctx.Context) ([]ExternalUserRelationData, error)
	UpdateRemarkData(cc ctx.Context, enterpriseID int64, id int64, remark string) error
}
