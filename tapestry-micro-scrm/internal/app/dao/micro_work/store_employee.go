package microWorkDao

import (
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	openReq "micro-scrm/internal/app/http/request/open"
	microWorkModel "micro-scrm/internal/app/models/micro_work"
	"micro-scrm/internal/pkg/db/gormx"
)

type StoreEmployee struct {
	dao       dao.Dao
	tableName string
}

func NewStoreEmployee() *StoreEmployee {
	return &StoreEmployee{
		dao:       dao.Dao{DB: gormx.GetWorkInstance()},
		tableName: microWorkModel.StoreEmployeeModel{}.TableName(),
	}
}

func (d *StoreEmployee) Create(c app.Ctx, data microWorkModel.StoreEmployeeModel) error {
	err := d.dao.WithContext(c).Table(d.tableName).Create(&data).Error
	return err
}

func (d *StoreEmployee) CreateInBatches(c app.Ctx, data []microWorkModel.StoreEmployeeModel) error {
	err := d.dao.WithContext(c).Table(d.tableName).CreateInBatches(&data, len(data)).Error
	return err
}

func (d *StoreEmployee) GetInfo(c app.Ctx, cond interface{}) (microWorkModel.StoreEmployeeModel, error) {
	var result microWorkModel.StoreEmployeeModel
	err := d.dao.WithContext(c).Table(d.tableName).Where(cond).Find(&result).Error
	return result, err
}

func (d *StoreEmployee) GetAll(c app.Ctx, cond interface{}) ([]microWorkModel.StoreEmployeeModel, error) {
	var result []microWorkModel.StoreEmployeeModel
	err := d.dao.WithContext(c).Table(d.tableName).Where(cond).Find(&result).Error
	return result, err
}

func (d *StoreEmployee) GetPagination(c app.Ctx, cond openReq.SaListReq) ([]microWorkModel.StoreEmployeeModel, int64, error) {
	var (
		total int64
		resp  []microWorkModel.StoreEmployeeModel
		err   error
	)

	query := d.dao.WithContext(c).Table(d.tableName).Where("enterprise_id = ? and deleted_at = 0", cond.EnterpriseId)
	query.Count(&total)

	offset := (cond.Page - 1) * cond.Size
	err = query.Order("created_at DESC").Limit(cond.Size).Offset(offset).Find(&resp).Error

	return resp, total, err
}

func (d *StoreEmployee) Save(c app.Ctx, data interface{}) error {
	err := d.dao.WithContext(c).Table(d.tableName).Save(&data).Error
	return err
}

func (d *StoreEmployee) Update(c app.Ctx, data interface{}, cond interface{}) error {
	err := d.dao.WithContext(c).Table(d.tableName).Where(cond).UpdateColumns(&data).Error
	return err
}

func (d *StoreEmployee) Delete(c app.Ctx, cond interface{}) error {
	var mod microWorkModel.StoreEmployeeModel
	err := d.dao.WithContext(c).Table(d.tableName).Where(cond).Delete(&mod).Error
	return err
}

func (d *StoreEmployee) ForceDelete(c app.Ctx, cond interface{}) error {
	var mod microWorkModel.StoreEmployeeModel
	err := d.dao.WithContext(c).Table(d.tableName).Where(cond).Unscoped().Delete(&mod).Error
	return err
}
