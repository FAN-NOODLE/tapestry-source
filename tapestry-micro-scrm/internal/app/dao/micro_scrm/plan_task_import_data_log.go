package microScrmDao

import (
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	"micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/pkg/db/gormx"
)

type PlanTaskImportDataLogDao struct {
	dao       dao.Dao
	tableName string
}

func NewPlanTaskImportDataLogDao() *PlanTaskImportDataLogDao {
	return &PlanTaskImportDataLogDao{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		tableName: microScrmModel.PlanTaskImportDataLogModel{}.TableName(),
	}
}

// InsertImportDataLog 导入的客户和导购的匹配记录
func (d *PlanTaskImportDataLogDao) InsertImportDataLog(c app.Ctx, data microScrmModel.PlanTaskImportDataLogModel, enterpriseID int64) error {
	return d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseID)).Create(&data).Error
}

// GetImportDataLogCount 查询导入的客户和导购的匹配记录
func (d *PlanTaskImportDataLogDao) GetImportDataLogCount(c app.Ctx, enterpriseID, fileID int64) (total int64, err error) {
	err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseID)).
		Where("file_id", fileID).Count(&total).Error
	return
}

func (d *PlanTaskImportDataLogDao) GetImportDataLogList(c app.Ctx, enterpriseID, fileID int64, page int, size int) (data []microScrmModel.PlanTaskImportDataLogModel, err error) {
	offset := (page - 1) * size
	err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseID)).
		Where("file_id", fileID).Offset(offset).Limit(size).Find(&data).Error
	return
}

// InsertImportDataLogBatch 导入的客户和导购的匹配记录
func (d *PlanTaskImportDataLogDao) InsertImportDataLogBatch(c app.Ctx, data []microScrmModel.PlanTaskImportDataLogModel, enterpriseID int64) error {
	return d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseID)).Create(&data).Error
}
