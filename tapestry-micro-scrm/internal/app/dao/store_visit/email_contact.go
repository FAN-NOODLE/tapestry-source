package storeVisitDao

import (
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	storeVisitModel "micro-scrm/internal/app/models/store_visit"
)

type EmailContactDao struct {
	dao.EnterpriseRepo[storeVisitModel.StoreVisitEmailContactModel]
	model     *storeVisitModel.StoreVisitEmailContactModel
	tableName string
}

func NewEmailContactDao() *EmailContactDao {
	return &EmailContactDao{
		EnterpriseRepo: dao.EnterpriseRepo[storeVisitModel.StoreVisitEmailContactModel]{
			Dao: dao.Dao{DB: app.DB()},
		},
		model:     new(storeVisitModel.StoreVisitEmailContactModel),
		tableName: storeVisitModel.StoreVisitEmailContactModel{}.TableName(),
	}
}
