package mysql

import (
	"database/sql"
	"time"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	"micro-scrm/internal/modules/shop_member"
	"micro-scrm/internal/pkg/db"
	"micro-scrm/internal/pkg/uid"
)

// ShopMemberRepository
type ShopMemberRepository struct {
	*repos.SQLRepo
}

// GetShopMemberLevel 通过levelID查询会员等级信息
func (u *ShopMemberRepository) GetShopMemberLevel(c ctx.Context, enterpriseID int64, levelID int64) (*shop_member.ShopMemberLevelEntity, error) {
	var entity shop_member.ShopMemberLevelEntity
	err := db.Get(c, u.SQLRepo, "t_shop_member_level/get-member-level-by-level-id", &entity, enterpriseID, levelID)
	if err != nil && err != sql.ErrNoRows {
		return &entity, err
	}
	return &entity, nil
}

// AddMemberLevel 添加商城会员等级
func (u *ShopMemberRepository) AddShopMemberLevel(c ctx.Context, level *shop_member.ShopMemberLevelEntity) (int64, error) {
	level.ID = uid.ID(c)
	level.CreatedAt = time.Now().Unix()
	level.UpdatedAt = time.Now().Unix()
	return db.Add(c, u.SQLRepo, "t_shop_member_level/insert", level)
}

// updateMemberLevel 更新商城会员等级
func (u *ShopMemberRepository) UpdateShopMemberLevel(c ctx.Context, level *shop_member.ShopMemberLevelEntity) (int64, error) {
	level.UpdatedAt = time.Now().Unix()
	return db.Update(c, u.SQLRepo, "t_shop_member_level/update", level)
}

// 获取等级列表
func (u *ShopMemberRepository) GetLevelList(c ctx.Context, enterpriseID int64, levelIDs []int64) ([]*shop_member.ShopMemberLevelEntity, error) {
	mapping := make(map[string]interface{})
	mapping["cond"] = struct {
		EnterpriseID int64   `json:"enterprise_id" desc:"企业ID"`
		LevelIDs     []int64 `json:"level_ids"`
		LevelCount   int
	}{
		EnterpriseID: enterpriseID,
		LevelIDs:     levelIDs,
		LevelCount:   len(levelIDs),
	}
	var list = make([]*shop_member.ShopMemberLevelEntity, 0)
	err := db.NamedSelect(c, u.SQLRepo, "t_shop_member_level/list", &list, mapping, &mapping)
	if err != nil && err != sql.ErrNoRows {
		return list, err
	}
	return list, nil
}
