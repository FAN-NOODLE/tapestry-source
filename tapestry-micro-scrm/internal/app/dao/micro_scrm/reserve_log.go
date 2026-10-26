package microScrmDao

import (
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/pkg/db/gormx"
)

type ReservesLog struct {
	dao       dao.Dao
	tableName string
}

func NewReserveLog() *ReservesLog {
	return &ReservesLog{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		tableName: microScrmModel.ReserveLogModel{}.TableName(),
	}
}

func (d *ReservesLog) Create(c app.Ctx, enterpriseId int64, data microScrmModel.ReserveLogModel) error {
	return d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Create(&data).Error
}

func (d *ReservesLog) List(c app.Ctx, enterpriseId int64, reserveOrderID int64, logTypes []int) (list []microScrmModel.ReserveLogModel, err error) {

	err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).
		Where("reserve_order_id = ?", reserveOrderID).
		Where("log_type in ?", logTypes).
		Order("created_at DESC").
		Find(&list).Error

	return
}
