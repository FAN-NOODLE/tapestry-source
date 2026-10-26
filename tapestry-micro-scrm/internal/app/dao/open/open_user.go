package openDao

import (
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	openModel "micro-scrm/internal/app/models/open"
	"micro-scrm/internal/pkg/db/gormx"
)

type OpenUser struct {
	dao       dao.Dao
	tableName string
	model     *openModel.OpenUserModel
}

func NewOpenUser() *OpenUser {
	return &OpenUser{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		tableName: openModel.OpenUserModel{}.TableName(),
		model:     new(openModel.OpenUserModel),
	}
}

func (d *OpenUser) GetInfo(c app.Ctx, data interface{}) (openModel.OpenUserModel, error) {
	var result openModel.OpenUserModel
	err := d.dao.WithContext(c).Model(&d.model).Where(data).Find(&result).Error
	return result, err
}
