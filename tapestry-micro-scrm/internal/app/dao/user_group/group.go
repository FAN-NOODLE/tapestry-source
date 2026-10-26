package userGroupDao

import (
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	userGroupModel "micro-scrm/internal/app/models/user_group"
	"micro-scrm/internal/pkg/db/gormx"
)

type Group struct {
	dao       dao.Dao
	tableName string
}

func NewGroup() *Group {
	return &Group{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		tableName: userGroupModel.UserGroupModel{}.TableName(),
	}
}

func (d *Group) GetList(c app.Ctx, data userGroupModel.UserGroupModel) ([]userGroupModel.UserGroupModel, error) {
	var result []userGroupModel.UserGroupModel
	err := d.dao.WithContext(c).Table(d.tableName).Where(data).Find(&result).Error
	return result, err
}
