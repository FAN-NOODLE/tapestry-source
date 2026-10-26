package microBossDao

import (
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	microBossModel "micro-scrm/internal/app/models/micro_boss"
	"micro-scrm/internal/pkg/db/gormx"
)

type Enterprise struct {
	dao       dao.Dao
	tableName string
}

func NewEnterprise() *Enterprise {
	return &Enterprise{
		dao:       dao.Dao{DB: gormx.GetBossInstance()},
		tableName: microBossModel.EnterpriseModel{}.TableName(),
	}
}

func (d *Enterprise) GetOne(c app.Ctx, query interface{}) (microBossModel.EnterpriseModel, error) {
	var result microBossModel.EnterpriseModel
	err := d.dao.WithContext(c).Table(d.tableName).Where(query).Find(&result).Error
	return result, err
}
