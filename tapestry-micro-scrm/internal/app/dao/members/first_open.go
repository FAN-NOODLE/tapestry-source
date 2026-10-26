package membersDao

import (
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	membersModel "micro-scrm/internal/app/models/members"
	"micro-scrm/internal/pkg/db/gormx"
)

type FirstOpen struct {
	dao       dao.Dao
	tableName string
	model     *membersModel.FirstOpenModel
}

func NewFirstOpen() *FirstOpen {
	return &FirstOpen{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		tableName: membersModel.FirstOpenModel{}.TableName(),
		model:     new(membersModel.FirstOpenModel),
	}
}

func (d *FirstOpen) GetOne(c app.Ctx, query interface{}) (membersModel.FirstOpenModel, error) {
	var result membersModel.FirstOpenModel
	err := d.dao.WithContext(c).Model(d.model).Table(d.tableName).Where(query).Find(&result).Error
	return result, err
}

func (d *FirstOpen) Create(c app.Ctx, data membersModel.FirstOpenModel) error {
	return d.dao.WithContext(c).Model(d.model).Create(&data).Error
}
