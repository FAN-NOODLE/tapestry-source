package microWorkDao

import (
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	microWorkModel "micro-scrm/internal/app/models/micro_work"
	"micro-scrm/internal/pkg/db/gormx"
)

type StoreSftpConfig struct {
	dao       dao.Dao
	tableName string
}

func NewStoreSftpConfig() *StoreSftpConfig {
	return &StoreSftpConfig{
		dao:       dao.Dao{DB: gormx.GetWorkInstance()},
		tableName: microWorkModel.StoreSftpConfigModel{}.TableName(),
	}
}

func (d *StoreSftpConfig) GetOne(c app.Ctx, query interface{}) (microWorkModel.StoreSftpConfigModel, error) {
	var result microWorkModel.StoreSftpConfigModel
	err := d.dao.WithContext(c).Table(d.tableName).Where(query).Find(&result).Error
	return result, err
}
