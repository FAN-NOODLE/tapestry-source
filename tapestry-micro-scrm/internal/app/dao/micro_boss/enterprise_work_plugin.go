package microBossDao

import (
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	microBossModel "micro-scrm/internal/app/models/micro_boss"
	"micro-scrm/internal/pkg/db/gormx"
)

type EnterpriseWorkPlugin struct {
	dao       dao.Dao
	tableName string
}

func NewEnterpriseWorkPlugin() *EnterpriseWorkPlugin {
	return &EnterpriseWorkPlugin{
		dao:       dao.Dao{DB: gormx.GetBossInstance()},
		tableName: microBossModel.EnterpriseWorkPluginModel{}.TableName(),
	}
}

func (d *EnterpriseWorkPlugin) GetOne(c app.Ctx, query interface{}) (microBossModel.EnterpriseWorkPluginModel, error) {
	var result microBossModel.EnterpriseWorkPluginModel
	err := d.dao.WithContext(c).Table(d.tableName).Where(query).Find(&result).Error
	return result, err
}
