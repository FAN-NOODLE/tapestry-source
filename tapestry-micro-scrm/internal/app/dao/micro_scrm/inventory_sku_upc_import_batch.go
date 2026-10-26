package microScrmDao

import (
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
)

type InventorySkuUpcImportBatchDao struct {
	dao.EnterpriseRepo[microScrmModel.InventorySkuUpcImportBatchModel]
	model *microScrmModel.InventorySkuUpcImportBatchModel
}

func NewInventorySkuUpcImportBatchDao() *InventorySkuUpcImportBatchDao {
	model := new(microScrmModel.InventorySkuUpcImportBatchModel)
	return &InventorySkuUpcImportBatchDao{
		EnterpriseRepo: dao.EnterpriseRepo[microScrmModel.InventorySkuUpcImportBatchModel]{
			Dao: dao.Dao{DB: app.DB()},
		},
		model: model,
	}
}

func (d *InventorySkuUpcImportBatchDao) GetEnabledActiveBatch(c app.Ctx, enterpriseId int64) (microScrmModel.InventorySkuUpcImportBatchModel, error) {
	var result microScrmModel.InventorySkuUpcImportBatchModel
	err := d.WithContext(c).
		Scopes(d.BranchDataBase(enterpriseId, d.model.TableName())).
		Model(d.model).
		Where("status = ? AND is_active = ? AND mapping_enabled = ? AND deleted_at = 0",
			microScrmModel.InventorySkuUpcBatchStatusSuccess,
			microScrmModel.InventorySkuUpcBatchActive,
			microScrmModel.InventorySkuUpcMappingEnabled,
		).
		Order("updated_at DESC").
		Find(&result).Error
	return result, err
}

func (d *InventorySkuUpcImportBatchDao) GetActiveSuccessBatch(c app.Ctx, enterpriseId int64) (microScrmModel.InventorySkuUpcImportBatchModel, error) {
	var result microScrmModel.InventorySkuUpcImportBatchModel
	err := d.WithContext(c).
		Scopes(d.BranchDataBase(enterpriseId, d.model.TableName())).
		Model(d.model).
		Where("status = ? AND is_active = ? AND deleted_at = 0",
			microScrmModel.InventorySkuUpcBatchStatusSuccess,
			microScrmModel.InventorySkuUpcBatchActive,
		).
		Order("updated_at DESC").
		Find(&result).Error
	return result, err
}

func (d *InventorySkuUpcImportBatchDao) DisableActiveBatches(c app.Ctx, enterpriseId int64) error {
	return d.UpdateData(c, enterpriseId, map[string]interface{}{
		"is_active":  microScrmModel.InventorySkuUpcBatchActive,
		"deleted_at": 0,
	}, map[string]any{
		"is_active": microScrmModel.InventorySkuUpcBatchInactive,
	})
}
