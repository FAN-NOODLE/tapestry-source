// Package callback 对接Open 平台
package enterprise

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
)

// Repository 存储定义
type Repository interface {
	repos.IRepository
	FindEnterpriseDetail(ctx ctx.Context, c Corp) (*Corp, error)
	ListEnterprise(cc ctx.Context, enterpriseIds []int64, page int, size int) (int64, *[]*Corp, error)
	UpdateEnterprise(cc ctx.Context, c Corp) error
	OfficialInfo(context ctx.Context, id int64) (*OfficialInfo, error)
	OfficialUser(context ctx.Context, req *OfficialUserReq) (*OfficialUser, error)
	CreateOfficialUser(context ctx.Context, user *OfficialUser) error
	CreateOfficialInfo(context ctx.Context, data *EpConfig) error

	QueryAppTypeCount(c ctx.Context, enterpriseID int64) (map[int]int, error)
	AddWxAppConfig(c ctx.Context, data *WxAppConfigEntity) (int64, error)
	QueryOneApp(c ctx.Context, enterpriseID int64, apptype int) (string, error)
}
