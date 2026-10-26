package onlineWarehouseDao

import (
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	onlineWarehouseModel "micro-scrm/internal/app/models/micro_scrm/online_warehouse"
)

type StoreDao struct {
	dao.EnterpriseRepo[onlineWarehouseModel.OnlineWarehouseStoreModel]
	model     *onlineWarehouseModel.OnlineWarehouseStoreModel
	tableName string
}

func NewStoreDao() *StoreDao {
	return &StoreDao{
		EnterpriseRepo: dao.EnterpriseRepo[onlineWarehouseModel.OnlineWarehouseStoreModel]{
			Dao: dao.Dao{DB: app.DB()},
		},
		model:     new(onlineWarehouseModel.OnlineWarehouseStoreModel),
		tableName: onlineWarehouseModel.OnlineWarehouseStoreModel{}.TableName(),
	}
}
