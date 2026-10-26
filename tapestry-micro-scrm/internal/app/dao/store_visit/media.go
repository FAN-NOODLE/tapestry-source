package storeVisitDao

import (
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	storeVisitModel "micro-scrm/internal/app/models/store_visit"
)

type MediaDao struct {
	dao.EnterpriseRepo[storeVisitModel.StoreVisitMediaModel]
	model     *storeVisitModel.StoreVisitMediaModel
	tableName string
}

func NewMediaDao() *MediaDao {
	return &MediaDao{
		EnterpriseRepo: dao.EnterpriseRepo[storeVisitModel.StoreVisitMediaModel]{
			Dao: dao.Dao{DB: app.DB()},
		},
		model:     new(storeVisitModel.StoreVisitMediaModel),
		tableName: storeVisitModel.StoreVisitMediaModel{}.TableName(),
	}
}
