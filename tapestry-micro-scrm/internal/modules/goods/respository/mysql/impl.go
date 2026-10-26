package mysql

import (
	"database/sql"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	"micro-scrm/internal/modules/goods"
	"micro-scrm/internal/pkg/db"
	"micro-scrm/internal/pkg/errorIgnore"
)

// GoodsRepo 活动仓库实现
type GoodsRepo struct {
	*repos.SQLRepo
}

// SkuList 查询SPU列表
func (g GoodsRepo) SkuList(c ctx.Context, pageVo repos.PageVO, cond *goods.ListReq) (int64, []*goods.SkuItem, error) {
	option := &repos.FindPageOption{
		TplPrefix: "t_goods/sku_list",
		Fn: func(context ctx.Context, i int64) interface{} {
			return &[]*goods.SkuItem{}
		},
	}

	total, list, err := g.FindPage(c, pageVo, cond, option)

	if errorIgnore.Ignore(err) != nil {
		return 0, nil, nil
	}

	return total, *list.(*[]*goods.SkuItem), err
}

// QuerySpu 查询spu数据
func (g GoodsRepo) QuerySpu(c ctx.Context, enterpriseID int64, spuID int64) (*goods.SpuEntity, error) {
	var spu goods.SpuEntity
	err := db.Get(c, g.SQLRepo, "t_goods_spu/query", &spu, enterpriseID, spuID, 0)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &spu, nil
}

// InsertSpu 添加spu数据
func (g GoodsRepo) InsertSpu(c ctx.Context, entity *goods.SpuEntity) (int64, error) {
	return db.Add(c, g.SQLRepo, "t_goods_spu/insert", entity)
}

// UpdateSpu 更新spu数据
func (g GoodsRepo) UpdateSpu(c ctx.Context, entity *goods.SpuEntity) (int64, error) {
	return db.Update(c, g.SQLRepo, "t_goods_spu/update", entity)
}

// QuerySpu 查询spu数据
func (g GoodsRepo) QuerySku(c ctx.Context, enterpriseID int64, spuID int64) (*goods.SkuEntity, error) {
	var sku goods.SkuEntity
	err := db.Get(c, g.SQLRepo, "t_goods_sku/query", &sku, enterpriseID, spuID, 0)
	return &sku, err
}

// InsertSpu 添加spu数据
func (g GoodsRepo) InsertSku(c ctx.Context, entity *goods.SkuEntity) (int64, error) {
	return db.Add(c, g.SQLRepo, "t_goods_sku/insert", entity)
}

// UpdateSpu 更新spu数据
func (g GoodsRepo) UpdateSku(c ctx.Context, entity *goods.SkuEntity) (int64, error) {
	return db.Update(c, g.SQLRepo, "t_goods_sku/update", entity)
}

// QuerySkuList 查询多个sku数据
func (g GoodsRepo) QuerySkuList(c ctx.Context, skuIDList []int64) ([]*goods.SkuEntity, error) {
	var skuList = make([]*goods.SkuEntity, 0)
	if len(skuIDList) == 0 {
		return skuList, nil
	}
	mapping := make(map[string]interface{})
	mapping["sku_list"] = skuIDList
	mapping["deleted_at"] = 0
	err := db.NamedSelect(c, g.SQLRepo, "t_goods_sku/query-sku-id-list", &skuList, mapping, &mapping)
	return skuList, err
}

// QuerySkuListOfSpu 根据spuID查询sku列表
func (g GoodsRepo) QuerySkuListOfSpu(c ctx.Context, enterpriseID int64, spuID int64) ([]*goods.SkuEntity, error) {
	var skuList = make([]*goods.SkuEntity, 0)
	err := db.Query(c, g.SQLRepo, "t_goods_sku/query-spu-sku-list", &skuList, enterpriseID, spuID, 0)
	return skuList, err
}

// SpuList SPU列表
func (g GoodsRepo) SpuList(c ctx.Context, pageVo repos.PageVO, cond *goods.SpuListRequest) (int64, []*goods.SpuEntity, error) {
	option := &repos.FindPageOption{
		TplPrefix: "t_goods/spu_list",
		Fn: func(context ctx.Context, i int64) interface{} {
			return &[]*goods.SpuEntity{}
		},
	}

	total, list, err := g.FindPage(c, pageVo, cond, option)

	if errorIgnore.Ignore(err) != nil {
		return 0, nil, nil
	}

	return total, *list.(*[]*goods.SpuEntity), err
}

// QuerySpuList 查询spu列表
func (g GoodsRepo) QuerySpuList(c ctx.Context, spuIDList []int64) ([]*goods.SpuEntity, error) {
	var skuList = make([]*goods.SpuEntity, 0)
	if len(spuIDList) == 0 {
		return skuList, nil
	}
	mapping := make(map[string]interface{})
	mapping["spu_list"] = spuIDList
	mapping["deleted_at"] = 0
	err := db.NamedSelect(c, g.SQLRepo, "t_goods_spu/query-spu-id-list", &skuList, mapping, &mapping)
	return skuList, err
}
