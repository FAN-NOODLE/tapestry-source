package mysql

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/sqlplus"
	"micro-scrm/internal/modules/activity"
	"micro-scrm/internal/pkg/db"
	"micro-scrm/internal/pkg/errorIgnore"
)

// ActivityRepo 活动仓库实现
type ActivityRepo struct {
	*repos.SQLRepo
}

// GetActivity 获取活动信息
func (a *ActivityRepo) GetActivity(c ctx.Context, enterpriseID int64, activityID int64) (*activity.ShopActivityEntity, error) {
	var data activity.ShopActivityEntity
	err := a.get(c, "t_shop_activity/query", &data, enterpriseID, activityID, 0)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

// AddActivity 添加活动数据
func (a *ActivityRepo) AddActivity(c ctx.Context, entity *activity.ShopActivityEntity) (int64, error) {
	return a.add(c, "t_shop_activity/insert", entity)
}

// AddActivitySku 添加活动spu数据
func (a *ActivityRepo) AddActivitySku(c ctx.Context, spuList []*activity.ShopActivitySku) (int64, error) {
	return a.add(c, "t_shop_activity_sku/insert", spuList)
}

// UpdateActivity 更新活动信息
func (a *ActivityRepo) UpdateActivity(c ctx.Context, entity *activity.ShopActivityEntity) (int64, error) {
	return a.update(c, entity, "t_shop_activity/update")
}

// UpdateActivitySku 更新活动sku信息
func (a *ActivityRepo) UpdateActivitySku(c ctx.Context, entity *activity.ShopActivitySku) (int64, error) {
	return a.update(c, entity, "t_shop_activity_sku/update")
}

// get 查询封装
func (a *ActivityRepo) get(c ctx.Context, tpl string, dest interface{}, args ...interface{}) error {
	return db.Get(c, a.SQLRepo, tpl, dest, args...)
}

// add 插入封装
func (a *ActivityRepo) add(c ctx.Context, tpl string, entity interface{}) (int64, error) {
	return db.Add(c, a.SQLRepo, tpl, entity)
}

// update 更新封装
func (a *ActivityRepo) update(c ctx.Context, entity interface{}, tpl string) (int64, error) {
	return db.Update(c, a.SQLRepo, tpl, entity)
}

// GetActivitySkuList 获取活动spu列表
func (a *ActivityRepo) GetActivitySkuList(c ctx.Context, enterpriseID int64, activityID int64) ([]*activity.ShopActivitySku, error) {
	var list = make([]*activity.ShopActivitySku, 0)
	err := db.Query(c, a.SQLRepo, "t_shop_activity_sku/query", &list, enterpriseID, activityID, 0)
	return list, err
}

// AllActivity 查询所有的活动
func (a *ActivityRepo) AllActivity(c ctx.Context, request *activity.ShopActivityEntity) ([]*activity.ShopActivityEntity, error) {
	stmt, err := a.DB.Template("t_shop_activity/all", a.Cond2Args(request))
	if err != nil {
		return nil, err
	}
	list := make([]*activity.ShopActivityEntity, 0)
	res, err := stmt.NamedQuery(c, request)
	if err != nil {
		return nil, err
	}
	defer res.Close()
	sqlplusErr := sqlplus.ScanAll(res, &list)
	if sqlplusErr != nil {
		return nil, err
	}
	return list, nil

}

// SkuListByActivityIds 通过活动ID集合查询SKU信息
func (a *ActivityRepo) SkuListByActivityIds(c ctx.Context, request activity.SkuListByActivityIdsReq, pageVo repos.PageVO) (int64, []*activity.GoodsListItem, error) {
	option := &repos.FindPageOption{
		TplPrefix: "t_shop_activity_sku/sku_list_ids",
		Fn: func(context ctx.Context, i int64) interface{} {
			return &[]*activity.GoodsListItem{}
		},
	}

	total, list, err := a.FindPage(c, pageVo, request, option)

	if errorIgnore.Ignore(err) != nil {
		return 0, nil, nil
	}

	return total, *list.(*[]*activity.GoodsListItem), err
}
