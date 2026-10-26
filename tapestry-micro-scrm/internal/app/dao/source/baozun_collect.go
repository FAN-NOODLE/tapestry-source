package source

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	"micro-scrm/internal/app/models/source"
)

type BaozunCollectDao struct {
	dao.Repo[source.BaozunCollectModel]
	model     *source.BaozunCollectModel
	tableName string
}

func NewBaozunCollectDao() *BaozunCollectDao {
	return &BaozunCollectDao{
		Repo: dao.Repo[source.BaozunCollectModel]{
			DB: app.DB(),
		},
		model:     new(source.BaozunCollectModel),
		tableName: source.BaozunCollectModel{}.TableName(),
	}
}

func (d BaozunCollectDao) GetCollectNumByProductCode(c ctx.Context, enterpriseId int64, productCodes []string) ([]source.SourceProductCollectNumInfo, error) {
	var res []source.SourceProductCollectNumInfo
	err := d.WithContext(c).Table(d.tableName).Where(map[string]interface{}{
		"enterprise_id": enterpriseId,
		"product_code":  productCodes,
	}).Select("product_code,collect_num").Group("product_code").Scan(&res).Error
	return res, err
}
