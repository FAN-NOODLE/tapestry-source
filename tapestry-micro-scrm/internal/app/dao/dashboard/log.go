package dashboardDao

import (
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	dashboardModel "micro-scrm/internal/app/models/dashboard"
	"micro-scrm/internal/pkg/db/gormx"
)

type Log struct {
	dao       dao.Dao
	model     *dashboardModel.DashboardLogModel
	tableName string
}

func NewLog() *Log {
	return &Log{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		model:     new(dashboardModel.DashboardLogModel),
		tableName: dashboardModel.DashboardLogModel{}.TableName(),
	}
}

func (d *Log) CreateData(c app.Ctx, data dashboardModel.DashboardLogModel) error {
	err := d.dao.WithContext(c).Model(d.model).Create(&data).Error
	return err
}
