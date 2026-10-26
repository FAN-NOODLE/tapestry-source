package microScrmDao

import (
	"strings"

	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/pkg/helper"
)

type InventorySkuUpcMappingDao struct {
	dao.EnterpriseRepo[microScrmModel.InventorySkuUpcMappingModel]
	model *microScrmModel.InventorySkuUpcMappingModel
}

func NewInventorySkuUpcMappingDao() *InventorySkuUpcMappingDao {
	model := new(microScrmModel.InventorySkuUpcMappingModel)
	return &InventorySkuUpcMappingDao{
		EnterpriseRepo: dao.EnterpriseRepo[microScrmModel.InventorySkuUpcMappingModel]{
			Dao: dao.Dao{DB: app.DB()},
		},
		model: model,
	}
}

func (d *InventorySkuUpcMappingDao) GetUpcCodesBySkuItem(c app.Ctx, enterpriseId int64, batchNo string, skuItem string) ([]string, error) {
	var upcCodes []string
	err := d.WithContext(c).
		Scopes(d.BranchDataBase(enterpriseId, d.model.TableName())).
		Model(d.model).
		Where("batch_no = ? AND sku_item = ? AND deleted_at = 0", batchNo, strings.ToUpper(strings.TrimSpace(skuItem))).
		Pluck("upc_code", &upcCodes).Error
	if err != nil {
		return nil, err
	}
	return helper.ArrayUniq(upcCodes), nil
}
