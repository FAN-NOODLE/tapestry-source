package thirdDao

import (
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	thirdModel "micro-scrm/internal/app/models/third"
	"micro-scrm/internal/pkg/db/gormx"
)

type ReserveRecordAssignDao struct {
	dao       dao.Dao
	tableName string // 表名
}

func NewReserveRecordAssignDao() *ReserveRecordAssignDao {
	return &ReserveRecordAssignDao{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		tableName: thirdModel.ReserveRecordAssignModel{}.TableName(),
	}
}

// WithTx 切换db连接,用于事物
func (d *ReserveRecordAssignDao) WithTx(tx *gorm.DB) *ReserveRecordAssignDao {
	d.dao.DB = tx
	return d
}

func (d *ReserveRecordAssignDao) Create(c app.Ctx, data thirdModel.ReserveRecordAssignModel) error {
	err := d.dao.WithContext(c).Table(d.tableName).Create(&data).Error
	return err
}

func (d *ReserveRecordAssignDao) CreateInBatches(c app.Ctx, data []thirdModel.ReserveRecordAssignModel) error {
	err := d.dao.WithContext(c).Table(d.tableName).CreateInBatches(&data, len(data)).Error
	return err
}

func (d *ReserveRecordAssignDao) GetInfo(c app.Ctx, data thirdModel.ReserveRecordAssignModel) (thirdModel.ReserveRecordAssignModel, error) {
	var result thirdModel.ReserveRecordAssignModel
	err := d.dao.WithContext(c).Table(d.tableName).Where(data).Find(&result).Error
	return result, err
}

func (d *ReserveRecordAssignDao) GetAll(c app.Ctx, data thirdModel.ReserveRecordAssignModel) ([]thirdModel.ReserveRecordAssignModel, error) {
	var result []thirdModel.ReserveRecordAssignModel
	err := d.dao.WithContext(c).Table(d.tableName).Where(data).Find(&result).Error
	return result, err
}

func (d *ReserveRecordAssignDao) Save(c app.Ctx, data thirdModel.ReserveRecordAssignModel) error {
	err := d.dao.WithContext(c).Table(d.tableName).Save(&data).Error
	return err
}

func (d *ReserveRecordAssignDao) Update(c app.Ctx, data thirdModel.ReserveRecordAssignModel, cond interface{}) error {
	err := d.dao.WithContext(c).Table(d.tableName).Where(cond).UpdateColumns(&data).Error
	return err
}

func (d *ReserveRecordAssignDao) Delete(c app.Ctx, cond interface{}) error {
	var mod thirdModel.ReserveRecordAssignModel
	err := d.dao.WithContext(c).Table(d.tableName).Where(cond).Delete(&mod).Error
	return err
}

func (d *ReserveRecordAssignDao) ForceDelete(c app.Ctx, cond interface{}) error {
	var mod thirdModel.ReserveRecordAssignModel
	err := d.dao.WithContext(c).Table(d.tableName).Where(cond).Unscoped().Delete(&mod).Error
	return err
}
