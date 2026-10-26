package microWorkDao

import (
	"errors"
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	microWorkModel "micro-scrm/internal/app/models/micro_work"
	"micro-scrm/internal/pkg/db/gormx"
)

type StoreWorker struct {
	dao       dao.Dao
	tableName string
}

func NewStoreWorker() *StoreWorker {
	return &StoreWorker{
		dao:       dao.Dao{DB: gormx.GetWorkInstance()},
		tableName: microWorkModel.StoreWorkerModel{}.TableName(),
	}
}

func (d *StoreWorker) UpdateMobileByOldMobile(c app.Ctx, oldMobile, mobile string, enterpriseId int64) error {
	return d.dao.WithContext(c).Table(d.tableName).Where(map[string]interface{}{
		"mobile":        oldMobile,
		"enterprise_id": enterpriseId,
	}).Updates(map[string]interface{}{
		"mobile":     mobile,
		"sms_mobile": oldMobile,
	}).Error
}

func (d *StoreWorker) UpdateDataByMobile(c app.Ctx, data microWorkModel.StoreWorkerModel, mobile string, enterpriseId int64) error {
	return d.dao.WithContext(c).Table(d.tableName).Where(map[string]interface{}{
		"mobile":        mobile,
		"enterprise_id": enterpriseId,
	}).UpdateColumns(data).Error
}

func (d *StoreWorker) UpdateDataByUserId(c app.Ctx, data microWorkModel.StoreWorkerModel, memUserId string, enterpriseId int64) error {
	return d.dao.WithContext(c).Table(d.tableName).Where(map[string]interface{}{
		"mem_userid":    memUserId,
		"enterprise_id": enterpriseId,
	}).UpdateColumns(data).Error
}

func (d *StoreWorker) GetInfo(c app.Ctx, data microWorkModel.StoreWorkerModel) (microWorkModel.StoreWorkerModel, error) {
	var result microWorkModel.StoreWorkerModel
	err := d.dao.WithContext(c).Table(d.tableName).Where(data).First(&result).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	return result, err
}

func (d *StoreWorker) DisableByMemUserId(c app.Ctx, memUserId string, enterpriseId int64) error {
	return d.dao.WithContext(c).Table(d.tableName).
		Where(map[string]interface{}{
			"enterprise_id": enterpriseId,
			"mem_userid":    memUserId,
		}).
		Update("worker_status", microWorkModel.StoreWorkerStatusDisable).Error
}
