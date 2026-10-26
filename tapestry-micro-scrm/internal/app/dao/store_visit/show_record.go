package storeVisitDao

import (
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	storeVisitModel "micro-scrm/internal/app/models/store_visit"
)

type ShowRecordDao struct {
	dao.EnterpriseRepo[storeVisitModel.StoreVisitShowRecordModel]
	model     *storeVisitModel.StoreVisitShowRecordModel
	tableName string
}

func NewShowRecordDao() *ShowRecordDao {
	return &ShowRecordDao{
		EnterpriseRepo: dao.EnterpriseRepo[storeVisitModel.StoreVisitShowRecordModel]{
			Dao: dao.Dao{DB: app.DB()},
		},
		model:     new(storeVisitModel.StoreVisitShowRecordModel),
		tableName: storeVisitModel.StoreVisitShowRecordModel{}.TableName(),
	}
}

func (d *ShowRecordDao) GetListPreloadMedia(c app.Ctx, enterpriseId int64, query any, args ...any) ([]storeVisitModel.StoreVisitShowRecordModel, error) {
	var result []storeVisitModel.StoreVisitShowRecordModel
	err := d.WithContext(c).Scopes(d.BranchDataBase(enterpriseId, d.model.TableName())).
		Model(d.model).
		Preload("MediaList", func(db *gorm.DB) *gorm.DB {
			return db.Scopes(d.BranchTable(storeVisitModel.StoreVisitMediaModel{}.TableName(), enterpriseId)).
				Where("rel_type = ? AND deleted_at = 0", storeVisitModel.ShowRecord)
		}).
		Where(query, args...).Order("id desc").Find(&result).Error
	return result, err
}

func (d *ShowRecordDao) GetListPreloadMediaOrderByType(c app.Ctx, enterpriseId int64, query any, args ...any) ([]storeVisitModel.StoreVisitShowRecordModel, error) {
	var result []storeVisitModel.StoreVisitShowRecordModel
	err := d.WithContext(c).Scopes(d.BranchDataBase(enterpriseId, d.model.TableName())).
		Model(d.model).
		Preload("MediaList", func(db *gorm.DB) *gorm.DB {
			return db.Scopes(d.BranchTable(storeVisitModel.StoreVisitMediaModel{}.TableName(), enterpriseId)).
				Where("rel_type = ? AND deleted_at = 0", storeVisitModel.ShowRecord)
		}).
		Where(query, args...).Order("id desc").Find(&result).Error
	return result, err
}

func (d *ShowRecordDao) GetVisitIdByTypeName(c app.Ctx, enterpriseId int64, query any, args ...any) ([]int64, error) {

	var visitIds []int64
	err := d.WithContext(c).Scopes(d.BranchDataBase(enterpriseId, d.model.TableName())).
		Model(d.model).Where(query, args...).Group("visit_id").Pluck("visit_id", &visitIds).Error
	return visitIds, err
}
