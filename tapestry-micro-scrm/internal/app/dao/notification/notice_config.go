package notificationDao

import (
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/pkg/db/gormx"
)

type NoticeConfigDao struct {
	dao       dao.Dao
	model     *microScrmModel.NotificationNoticeConfigModel
	tableName string
}

func NewNoticeConfigDao() *NoticeConfigDao {
	return &NoticeConfigDao{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		model:     new(microScrmModel.NotificationNoticeConfigModel),
		tableName: microScrmModel.NotificationNoticeConfigModel{}.TableName(),
	}
}

func (d NoticeConfigDao) GetOne(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) (microScrmModel.NotificationNoticeConfigModel, error) {
	var result microScrmModel.NotificationNoticeConfigModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).Where(query, args...).Find(&result).Error
	return result, err
}

func (d NoticeConfigDao) GetAll(c app.Ctx, enterpriseId int64) ([]microScrmModel.NotificationNoticeConfigModel, error) {
	var result []microScrmModel.NotificationNoticeConfigModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).Find(&result).Error
	return result, err
}

func (d NoticeConfigDao) UpdateData(c app.Ctx, enterpriseId int64, query interface{}, data map[string]interface{}) error {
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(d.model).Where(query).UpdateColumns(data).Error
	return err
}
