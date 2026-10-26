package userGroupDao

import (
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	userGroupModel "micro-scrm/internal/app/models/user_group"
	"micro-scrm/internal/pkg/db/gormx"
)

type GroupImportDataLog struct {
	dao       dao.Dao
	tableName string
}

func NewGroupImportDataLog() *GroupImportDataLog {
	return &GroupImportDataLog{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		tableName: userGroupModel.UserGroupImportDataLogModel{}.TableName(),
	}
}

// CreateInBatches 批量插入数据
func (d *GroupImportDataLog) CreateInBatches(c app.Ctx, enterpriseId int64, data []userGroupModel.UserGroupImportDataLogModel) error {
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where(data).CreateInBatches(&data, len(data)).Error
	return err
}

func (d *GroupImportDataLog) Update(c app.Ctx, enterpriseId, groupId int64, unionId string) error {
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where(map[string]interface{}{
		"user_group_id": groupId,
		"union_id":      unionId,
	}).UpdateColumn("import_res", "成功").Error
	return err
}
