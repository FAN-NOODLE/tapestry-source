package activity

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
)

// Repository 存储定义
type Repository interface {
	repos.IRepository
	GetActivity(c ctx.Context, enterpriseID int64, activityID int64) (*ShopActivityEntity, error)
	AddActivity(c ctx.Context, entity *ShopActivityEntity) (int64, error)
	AddActivitySku(c ctx.Context, skuList []*ShopActivitySku) (int64, error)
	UpdateActivity(c ctx.Context, entity *ShopActivityEntity) (int64, error)
	UpdateActivitySku(c ctx.Context, entity *ShopActivitySku) (int64, error)
	GetActivitySkuList(c ctx.Context, enterpriseID int64, activityID int64) ([]*ShopActivitySku, error)
	// 查询所有的活动
	AllActivity(c ctx.Context, request *ShopActivityEntity) ([]*ShopActivityEntity, error)
	// SkuListByActivityIds 通过活动ID集合查询SKU信息
	SkuListByActivityIds(c ctx.Context, request SkuListByActivityIdsReq, pageVo repos.PageVO) (int64, []*GoodsListItem, error)
}
