package storeVisitDao

import (
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	storeVisitModel "micro-scrm/internal/app/models/store_visit"
)

type ConfigDao struct {
	dao.EnterpriseRepo[storeVisitModel.StoreVisitConfigModel]
	model     *storeVisitModel.StoreVisitConfigModel
	tableName string
}

func NewConfigDao() *ConfigDao {
	return &ConfigDao{
		EnterpriseRepo: dao.EnterpriseRepo[storeVisitModel.StoreVisitConfigModel]{
			Dao: dao.Dao{DB: app.DB()},
		},
		model:     new(storeVisitModel.StoreVisitConfigModel),
		tableName: storeVisitModel.StoreVisitConfigModel{}.TableName(),
	}
}
