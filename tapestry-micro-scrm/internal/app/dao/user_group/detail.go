package userGroupDao

import (
	"errors"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	userGroupModel "micro-scrm/internal/app/models/user_group"
	"micro-scrm/internal/pkg/db/gormx"

	"gorm.io/gorm"
)

type GroupDetail struct {
	dao       dao.Dao
	tableName string
}

func NewGroupDetail() *GroupDetail {
	return &GroupDetail{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		tableName: userGroupModel.UserGroupDetailModel{}.TableName(),
	}
}

func (d *GroupDetail) GetList(c app.Ctx, enterpriseId, groupId int64) ([]userGroupModel.UserGroupDetailModel, error) {
	var result []userGroupModel.UserGroupDetailModel
	err := d.dao.WithContext(c).Table(d.tableName).Where(map[string]interface{}{
		"enterprise_id": enterpriseId,
		"user_group_id": groupId,
	}).Find(&result).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return result, nil
	}
	return result, err
}

func (d *GroupDetail) UpdateDetailISExtract(c app.Ctx, enterpriseId, groupId int64) error {
	err := d.dao.WithContext(c).Table(d.tableName).Where(map[string]interface{}{
		"user_group_id": groupId,
		"enterprise_id": enterpriseId,
	}).UpdateColumn("is_extract", 1).Error
	return err
}
