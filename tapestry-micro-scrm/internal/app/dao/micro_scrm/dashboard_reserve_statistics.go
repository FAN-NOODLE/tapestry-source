package microScrmDao

import (
	"fmt"
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/pkg/db/gormx"
)

type ReserveStatistics struct {
	dao   dao.Dao
	model *microScrmModel.DashboardReserveStatisticsModel
}

func NewReserveStatistics() *ReserveStatistics {
	return &ReserveStatistics{
		dao:   dao.Dao{DB: gormx.GetInstance()},
		model: new(microScrmModel.DashboardReserveStatisticsModel),
	}
}

// WithTx 切换db连接,用于事物
func (d *ReserveStatistics) WithTx(tx *gorm.DB) *ReserveStatistics {
	d.dao.DB = tx
	return d
}

func (d *ReserveStatistics) GetOne(c app.Ctx, query interface{}, args ...interface{}) (microScrmModel.DashboardReserveStatisticsModel, error) {
	var result microScrmModel.DashboardReserveStatisticsModel
	err := d.dao.WithContext(c).Model(d.model).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *ReserveStatistics) GetMany(c app.Ctx, query interface{}, args ...interface{}) ([]microScrmModel.DashboardReserveStatisticsModel, error) {
	var result []microScrmModel.DashboardReserveStatisticsModel
	err := d.dao.WithContext(c).Model(d.model).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *ReserveStatistics) CreateData(c app.Ctx, data microScrmModel.DashboardReserveStatisticsModel) error {
	err := d.dao.WithContext(c).Model(d.model).Create(&data).Error
	return err
}

func (d *ReserveStatistics) CreateInBatchesData(c app.Ctx, data []microScrmModel.DashboardReserveStatisticsModel) error {
	err := d.dao.WithContext(c).Model(d.model).CreateInBatches(&data, len(data)).Error
	return err
}

func (d *ReserveStatistics) SaveData(c app.Ctx, data microScrmModel.DashboardReserveStatisticsModel) error {
	err := d.dao.WithContext(c).Model(d.model).Save(&data).Error
	return err
}

func (d *ReserveStatistics) UpdateData(c app.Ctx, query interface{}, data map[string]interface{}) error {
	err := d.dao.WithContext(c).Model(d.model).Where(query).UpdateColumns(data).Error
	return err
}

func (d *ReserveStatistics) DeleteData(c app.Ctx, query interface{}, args ...interface{}) error {
	var model microScrmModel.DashboardReserveStatisticsModel
	err := d.dao.WithContext(c).Model(d.model).Where(query, args...).Delete(&model).Error
	return err
}

func (d *ReserveStatistics) ForceDelete(c app.Ctx, query interface{}, args ...interface{}) error {
	var model microScrmModel.DashboardReserveStatisticsModel
	err := d.dao.WithContext(c).Model(d.model).Where(query, args...).Unscoped().Delete(&model).Error
	return err
}

func (d *ReserveStatistics) Increment(c app.Ctx, query interface{}, column string, num int64) error {
	err := d.dao.WithContext(c).Model(d.model).Where(query).UpdateColumn(column, gorm.Expr(fmt.Sprintf("%s + ?", column), num)).Error
	return err
}

func (d *ReserveStatistics) Decrement(c app.Ctx, query interface{}, column string, num int64) error {
	err := d.dao.WithContext(c).Model(d.model).Where(query).UpdateColumn(column, gorm.Expr(fmt.Sprintf("%s - ?", column), num)).Error
	return err
}
