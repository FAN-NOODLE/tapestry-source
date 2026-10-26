package userGroupDao

import (
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	userGroupModel "micro-scrm/internal/app/models/user_group"
	"micro-scrm/internal/pkg/db/gormx"
)

type GroupImportData struct {
	dao       dao.Dao
	tableName string
}

func NewGroupImportData() *GroupImportData {
	return &GroupImportData{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		tableName: userGroupModel.UserGroupImportDataModel{}.TableName(),
	}
}

func (d *GroupImportData) GetImportDataCount(c app.Ctx, enterpriseId, fileId int64) (total int64, err error) {
	err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).
		Where("file_id = ?", fileId).Count(&total).Error
	return
}

func (d *GroupImportData) GetImportData(c app.Ctx, enterpriseId, fileId int64, offset, size int) (data []userGroupModel.UserGroupImportDataModel, err error) {
	err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).
		Where("file_id = ?", fileId).Offset(offset).Limit(size).Find(&data).Error
	return
}
