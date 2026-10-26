// 定义存储使用的接口和辅助方法
package shop_member

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
)

type Repository interface {
	repos.IRepository
	// 查询会员等级
	GetShopMemberLevel(c ctx.Context, enterpriseID int64, levelID int64) (*ShopMemberLevelEntity, error)
	// 添加会员等级
	AddShopMemberLevel(c ctx.Context, params *ShopMemberLevelEntity) (int64, error)
	// 更新会员等级
	UpdateShopMemberLevel(c ctx.Context, params *ShopMemberLevelEntity) (int64, error)
	// GetLevelList 获取等级列表
	GetLevelList(c ctx.Context, enterpriseID int64, levelIDs []int64) ([]*ShopMemberLevelEntity, error)
}
