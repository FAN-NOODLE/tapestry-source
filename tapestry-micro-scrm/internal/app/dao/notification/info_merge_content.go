package notificationDao

import (
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/pkg/db/gormx"
)

type InfoMergeContentDao struct {
	dao       dao.Dao
	model     *microScrmModel.NotificationInfoMergeContentModel
	tableName string
}

func NewInfoMergeContentDao() *InfoMergeContentDao {
	return &InfoMergeContentDao{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		model:     new(microScrmModel.NotificationInfoMergeContentModel),
		tableName: microScrmModel.NotificationInfoMergeContentModel{}.TableName(),
	}
}

func (d InfoMergeContentDao) GetOne(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) (microScrmModel.NotificationInfoMergeContentModel, error) {
	var result microScrmModel.NotificationInfoMergeContentModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).Where(query, args...).Find(&result).Error
	return result, err
}

func (d InfoMergeContentDao) GetMany(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) ([]microScrmModel.NotificationInfoMergeContentModel, error) {
	var result []microScrmModel.NotificationInfoMergeContentModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).Where(query, args...).Find(&result).Error
	return result, err
}

func (d InfoMergeContentDao) CreateData(c app.Ctx, enterpriseId int64, data microScrmModel.NotificationInfoMergeContentModel) error {
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).Create(&data).Error
	return err
}

func (d InfoMergeContentDao) CreateInBatchesData(c app.Ctx, enterpriseId int64, data []microScrmModel.NotificationInfoMergeContentModel) error {
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).CreateInBatches(&data, len(data)).Error
	return err
}

func (d InfoMergeContentDao) UpdateData(c app.Ctx, enterpriseId int64, query interface{}, data map[string]interface{}) error {
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).Where(query).UpdateColumns(data).Error
	return err
}

func (d InfoMergeContentDao) DeleteData(c app.Ctx, enterpriseId int64, query interface{}) error {
	var model microScrmModel.NotificationInfoMergeContentModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).Where(query).Delete(&model).Error
	return err
}
func (d InfoMergeContentDao) UnscopedDeleteData(c app.Ctx, enterpriseId int64, query interface{}) error {
	var model microScrmModel.NotificationInfoMergeContentModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).Where(query).Unscoped().Delete(&model).Error
	return err
}
