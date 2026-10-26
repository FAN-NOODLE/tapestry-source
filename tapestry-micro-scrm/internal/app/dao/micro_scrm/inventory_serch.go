package microScrmDao

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
)

type InventorySearchDao struct {
	dao.Repo[microScrmModel.InventorySearchModel]
	model     *microScrmModel.InventorySearchModel
	tableName string
}

func NewInventorySearchDao() *InventorySearchDao {
	return &InventorySearchDao{
		Repo: dao.Repo[microScrmModel.InventorySearchModel]{
			DB: app.DB(),
		},
		model:     new(microScrmModel.InventorySearchModel),
		tableName: microScrmModel.InventorySearchModel{}.TableName(),
	}
}

func (d *InventorySearchDao) GetKeywords(c ctx.Context, enterpriseId, memberID int64) ([]microScrmModel.InventorySearchModel, error) {
	res := make([]microScrmModel.InventorySearchModel, 0)
	err := d.WithContext(c).Model(d.model).Where(map[string]interface{}{
		"enterprise_id": enterpriseId,
		"member_id":     memberID,
	}).Order("updated_at DESC").Find(&res).Error
	return res, err
}

func (d *InventorySearchDao) GetLimitInfo(c app.Ctx, query any, args ...any) (microScrmModel.InventorySearchModel, error) {
	var res microScrmModel.InventorySearchModel
	err := d.WithContext(c).Model(d.model).Where(query, args...).Order("updated_at DESC").Offset(microScrmModel.LimitNum).Limit(1).Find(&res).Error
	return res, err
}
