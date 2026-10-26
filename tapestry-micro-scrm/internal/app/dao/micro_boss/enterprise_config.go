package microBossDao

import (
	"fmt"
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	microBossModel "micro-scrm/internal/app/models/micro_boss"
	"micro-scrm/internal/pkg/db/gormx"
)

type EnterpriseConfig struct {
	dao       dao.Dao
	tableName string
}

func NewEnterpriseConfig() *EnterpriseConfig {
	return &EnterpriseConfig{
		dao:       dao.Dao{DB: gormx.GetBossInstance()},
		tableName: microBossModel.EnterpriseConfigModel{}.TableName(),
	}
}

func (d *EnterpriseConfig) Create(c app.Ctx, data microBossModel.EnterpriseConfigModel) error {
	err := d.dao.WithContext(c).Table(d.tableName).Create(&data).Error
	return err
}

func (d *EnterpriseConfig) CreateInBatches(c app.Ctx, data []microBossModel.EnterpriseConfigModel) error {
	err := d.dao.WithContext(c).Table(d.tableName).CreateInBatches(&data, len(data)).Error
	return err
}

func (d *EnterpriseConfig) GetInfo(c app.Ctx, data interface{}) (microBossModel.EnterpriseConfigModel, error) {
	var result microBossModel.EnterpriseConfigModel
	err := d.dao.WithContext(c).Table(d.tableName).Where(data).Find(&result).Error
	return result, err
}

func (d *EnterpriseConfig) GetAll(c app.Ctx, data interface{}) ([]microBossModel.EnterpriseConfigModel, error) {
	var result []microBossModel.EnterpriseConfigModel
	err := d.dao.WithContext(c).Table(d.tableName).Where(data).Find(&result).Error
	return result, err
}

func (d *EnterpriseConfig) Save(c app.Ctx, data interface{}) error {
	err := d.dao.WithContext(c).Table(d.tableName).Save(&data).Error
	return err
}

func (d *EnterpriseConfig) Update(c app.Ctx, data interface{}, cond interface{}) error {
	err := d.dao.WithContext(c).Table(d.tableName).Where(cond).UpdateColumns(&data).Error
	return err
}

func (d *EnterpriseConfig) Delete(c app.Ctx, cond interface{}) error {
	var mod microBossModel.EnterpriseConfigModel
	err := d.dao.WithContext(c).Table(d.tableName).Where(cond).Delete(&mod).Error
	return err
}

func (d *EnterpriseConfig) ForceDelete(c app.Ctx, cond interface{}) error {
	var mod microBossModel.EnterpriseConfigModel
	err := d.dao.WithContext(c).Table(d.tableName).Where(cond).Unscoped().Delete(&mod).Error
	return err
}

func (d *EnterpriseConfig) Increment(c app.Ctx, cond interface{}, column string, num int64) error {
	err := d.dao.WithContext(c).Table(d.tableName).
		Where(cond).
		UpdateColumn(column, gorm.Expr(fmt.Sprintf("%s + ?", column), num)).Error
	return err
}

func (d *EnterpriseConfig) Decrement(c app.Ctx, cond interface{}, column string, num int64) error {
	err := d.dao.WithContext(c).Table(d.tableName).
		Where(cond).
		UpdateColumn(column, gorm.Expr(fmt.Sprintf("%s - ?", column), num)).Error
	return err
}
