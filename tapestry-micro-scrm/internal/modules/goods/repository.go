package goods

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
)

// Repository 存储定义
type Repository interface {
	repos.IRepository
	// SkuList 查询SKU列表
	SkuList(c ctx.Context, pageVo repos.PageVO, cond *ListReq) (int64, []*SkuItem, error)
	// SpuList 查询SPU列表
	SpuList(c ctx.Context, pageVo repos.PageVO, cond *SpuListRequest) (int64, []*SpuEntity, error)
	QuerySpu(c ctx.Context, enterpriseID int64, spuID int64) (*SpuEntity, error)
	InsertSpu(c ctx.Context, entity *SpuEntity) (int64, error)
	UpdateSpu(c ctx.Context, entity *SpuEntity) (int64, error)
	QuerySku(c ctx.Context, enterpriseID int64, spuID int64) (*SkuEntity, error)
	InsertSku(c ctx.Context, entity *SkuEntity) (int64, error)
	UpdateSku(c ctx.Context, entity *SkuEntity) (int64, error)
	QuerySkuList(c ctx.Context, skuIDList []int64) ([]*SkuEntity, error)
	QuerySpuList(c ctx.Context, spuIDList []int64) ([]*SpuEntity, error)
	QuerySkuListOfSpu(c ctx.Context, enterpriseID int64, spuID int64) ([]*SkuEntity, error)
}
