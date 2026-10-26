package microScrmDao

import (
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	"micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/pkg/db/gormx"
)

type PlanTaskImportDataDao struct {
	dao       dao.Dao
	tableName string
}

func NewPlanTaskImportDataDao() *PlanTaskImportDataDao {
	return &PlanTaskImportDataDao{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		tableName: microScrmModel.PlanTaskImportDataModel{}.TableName(),
	}
}

// InsertImportData 插入匹配成功的导购和客户的数据
func (d *PlanTaskImportDataDao) InsertImportData(c app.Ctx, data microScrmModel.PlanTaskImportDataModel, enterpriseID int64) error {
	return d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseID)).Create(&data).Error
}

// GetImportDataCount 获得匹配的导购和客户数据的count值
func (d *PlanTaskImportDataDao) GetImportDataCount(c app.Ctx, enterpriseID, fileID int64) (total int64, err error) {
	err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseID)).
		Where("file_id", fileID).Where("is_match", 1).Count(&total).Error
	return
}

// GetImportDataList 查询匹配的导购和客户数
func (d *PlanTaskImportDataDao) GetImportDataList(c app.Ctx, enterpriseID, fileID int64, page int, size int) (data []microScrmModel.PlanTaskImportDataModel, err error) {
	offset := (page - 1) * size
	err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseID)).
		Where("file_id", fileID).Where("is_match", 1).Offset(offset).Limit(size).Find(&data).Error
	return
}

// InsertImportDataBatch 批量插入匹配成功的导购和客户的数据
func (d *PlanTaskImportDataDao) InsertImportDataBatch(c app.Ctx, data []microScrmModel.PlanTaskImportDataModel, enterpriseID int64) error {
	return d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseID)).Create(&data).Error
}
