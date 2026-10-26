package user

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	"micro-scrm/internal/modules/customer"
)

// Repository 用户仓库
type Repository interface {
	repos.IRepository
	// EnterpriseUser 获取聚合用户
	EnterpriseUsers(cc ctx.Context, enterpriseID int64, euids []int64) ([]*customer.EnterpriseUser, error)
	QueryByID(c ctx.Context, enterpriseID int64, customerType int, id int64) (*customer.SourceUserBase, error)
	UpdateDescription(ctx.Context, int64, int64, string) error

	MergeUser(c ctx.Context, source *SourceUser) (bool, error)
	GetRobotUserList(c ctx.Context, limit int) ([]*RobotUserEntity, error)
	QueryUserCount(c ctx.Context, cond *customer.ListCond) (int, error)
}

// EventRepository 用户事件仓库
type EventRepository interface {
	repos.IRepository
}

// ShopOrderRepository 商城订单仓库
type ShopOrderRepository interface {
	repos.IRepository
}
