package thirdDao

import (
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	thirdModel "micro-scrm/internal/app/models/third"
	"micro-scrm/internal/pkg/db/gormx"
)

type ReserveRecordDao struct {
	dao       dao.Dao
	tableName string // 表名
}

func NewReserveRecordDao() *ReserveRecordDao {
	return &ReserveRecordDao{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		tableName: thirdModel.ReserveRecordModel{}.TableName(),
	}
}

// WithTx 切换db连接,用于事物
func (d *ReserveRecordDao) WithTx(tx *gorm.DB) *ReserveRecordDao {
	d.dao.DB = tx
	return d
}

func (d *ReserveRecordDao) Create(c app.Ctx, data thirdModel.ReserveRecordModel) error {
	err := d.dao.WithContext(c).Table(d.tableName).Create(&data).Error
	return err
}

func (d *ReserveRecordDao) CreateInBatches(c app.Ctx, data []thirdModel.ReserveRecordModel) error {
	err := d.dao.WithContext(c).Table(d.tableName).CreateInBatches(&data, len(data)).Error
	return err
}

func (d *ReserveRecordDao) GetInfo(c app.Ctx, data interface{}) (thirdModel.ReserveRecordModel, error) {
	var result thirdModel.ReserveRecordModel
	err := d.dao.WithContext(c).Table(d.tableName).Where(data).Find(&result).Error
	return result, err
}

func (d *ReserveRecordDao) GetList(c app.Ctx, data thirdModel.ReserveRecordModel) ([]thirdModel.ReserveRecordModel, error) {
	var result []thirdModel.ReserveRecordModel
	err := d.dao.WithContext(c).Table(d.tableName).Where(data).Find(&result).Error
	return result, err
}

func (d *ReserveRecordDao) Save(c app.Ctx, data thirdModel.ReserveRecordModel) error {
	err := d.dao.WithContext(c).Table(d.tableName).Save(&data).Error
	return err
}

func (d *ReserveRecordDao) Update(c app.Ctx, data thirdModel.ReserveRecordModel, cond interface{}) error {
	err := d.dao.WithContext(c).Table(d.tableName).Where(cond).UpdateColumns(&data).Error
	return err
}

func (d *ReserveRecordDao) Delete(c app.Ctx, cond interface{}) error {
	var mod thirdModel.ReserveRecordModel
	err := d.dao.WithContext(c).Table(d.tableName).Where(cond).Delete(&mod).Error
	return err
}

func (d *ReserveRecordDao) ForceDelete(c app.Ctx, cond interface{}) error {
	var mod thirdModel.ReserveRecordModel
	err := d.dao.WithContext(c).Table(d.tableName).Where(cond).Unscoped().Delete(&mod).Error
	return err
}

func (d *ReserveRecordDao) GetMany(c app.Ctx, query interface{}, args ...interface{}) ([]thirdModel.ReserveRecordModel, error) {
	var result []thirdModel.ReserveRecordModel
	err := d.dao.WithContext(c).Table(d.tableName).Where(query, args).Find(&result).Error
	return result, err
}
