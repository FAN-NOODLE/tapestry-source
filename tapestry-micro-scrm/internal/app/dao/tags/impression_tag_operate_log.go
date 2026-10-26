package tagsDao

import (
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	tagModels "micro-scrm/internal/app/models/tag"
	"micro-scrm/internal/pkg/db/gormx"
	"time"
)

type ImpressionOperateLog struct {
	dao       dao.Dao
	tableName string
}

func NewTagsOperateLog() *ImpressionOperateLog {
	return &ImpressionOperateLog{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		tableName: tagModels.ImpressionTagOperateLogModel{}.TableName(),
	}
}

func (d *ImpressionOperateLog) BatchCreate(c app.Ctx, data []tagModels.ImpressionTagOperateLogModel) error {
	return d.dao.WithContext(c).Table(d.tableName).CreateInBatches(&data, len(data)).Error
}

func (d *ImpressionOperateLog) Create(c app.Ctx, data tagModels.ImpressionTagOperateLogModel) error {
	return d.dao.WithContext(c).Table(d.tableName).Create(&data).Error
}

func (d *ImpressionOperateLog) GetLogList(c app.Ctx, enterpriseId int64) (data []tagModels.ImpressionTagOperateLogModel, err error) {
	err = d.dao.WithContext(c).Table(d.tableName).
		Select("operate_user_id,third_tag_id,action_type,external_unionid,created_at").
		Where("enterprise_id = ?", enterpriseId).
		Where("deleted_at = 0").Find(&data).Error
	return data, err
}

func (d *ImpressionOperateLog) DeleteList(c app.Ctx, enterpriseId int64, createdAt int64) error {
	return d.dao.WithContext(c).Table(d.tableName).
		Where("enterprise_id = ?", enterpriseId).
		Where("created_at <= ?", createdAt).
		Update("deleted_at", time.Now().Unix()).Error
}

func (d *ImpressionOperateLog) DeleteBeforeRecord(c app.Ctx, enterpriseId int64, data tagModels.ImpressionTagOperateLogModel) error {
	return d.dao.WithContext(c).Table(d.tableName).
		Where("enterprise_id = ?", enterpriseId).
		Where("operate_user_id = ?", data.OperateUserId).
		Where("external_unionid = ?", data.ExternalUnionId).
		Where("third_tag_id = ?", data.ThirdTagID).
		Unscoped().
		Delete(&tagModels.ImpressionTagOperateLogModel{}).Error
}

func (d *ImpressionOperateLog) GetMany(c app.Ctx, query interface{}, args ...interface{}) (data []tagModels.ImpressionTagOperateLogModel, err error) {
	err = d.dao.WithContext(c).Table(d.tableName).
		Where(query, args...).Find(&data).Error
	return data, err
}
