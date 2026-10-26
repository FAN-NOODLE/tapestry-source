package microScrmDao

import (
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	"micro-scrm/internal/orm/model/micro_scrm_enterprise"
	"micro-scrm/internal/pkg/db/gormx"
)

type WxUserUnion struct {
	dao       dao.Dao
	tableName string
}

func NewWxUserUnion() *WxUserUnion {
	return &WxUserUnion{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		tableName: micro_scrm_enterprise.WxUserUnionModel{}.TableName(),
	}
}

func (d *WxUserUnion) DeleteData(c app.Ctx, enterpriseId int64, query interface{}) error {
	var model []micro_scrm_enterprise.WxUserUnionModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where(query).Unscoped().Delete(&model).Error
	return err
}

func (d *WxUserUnion) CreateInBatchesData(c app.Ctx, enterpriseId int64, data []micro_scrm_enterprise.WxUserUnionModel) error {
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).CreateInBatches(&data, len(data)).Error
	return err
}

func (d *WxUserUnion) GetMany(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) ([]micro_scrm_enterprise.WxUserUnionModel, error) {
	var result []micro_scrm_enterprise.WxUserUnionModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where(query, args...).Find(&result).Error
	return result, err
}
