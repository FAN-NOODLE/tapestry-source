package drainageContactDao

import (
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	"micro-scrm/internal/app/models/drainage"
	"micro-scrm/internal/pkg/db/gormx"
)

type DrainageContactMemberDao struct {
	dao       dao.Dao
	tableName string
	model     *drainageContactModel.DrainageContactMemberModel
}

func NewDrainageContactMemberDao() *DrainageContactMemberDao {
	return &DrainageContactMemberDao{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		tableName: drainageContactModel.DrainageContactMemberModel{}.TableName(),
		model:     new(drainageContactModel.DrainageContactMemberModel),
	}
}

func (d *DrainageContactMemberDao) GetOne(c app.Ctx, query interface{}, args ...interface{}) (drainageContactModel.DrainageContactMemberModel, error) {
	var result drainageContactModel.DrainageContactMemberModel
	err := d.dao.WithContext(c).Model(d.model).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *DrainageContactMemberDao) GetMany(c app.Ctx, query interface{}, args ...interface{}) ([]drainageContactModel.DrainageContactMemberModel, error) {
	var result []drainageContactModel.DrainageContactMemberModel
	err := d.dao.WithContext(c).Model(d.model).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *DrainageContactMemberDao) CreateData(c app.Ctx, data drainageContactModel.DrainageContactMemberModel) error {
	err := d.dao.WithContext(c).Model(d.model).Create(&data).Error
	return err
}

func (d *DrainageContactMemberDao) UpdateData(c app.Ctx, query interface{}, data map[string]interface{}) error {
	err := d.dao.WithContext(c).Model(d.model).Where(query).UpdateColumns(data).Error
	return err
}

func (d *DrainageContactMemberDao) CreateInBatchesData(c app.Ctx, data []drainageContactModel.DrainageContactMemberModel) error {
	err := d.dao.WithContext(c).Model(d.model).CreateInBatches(&data, len(data)).Error
	return err
}
