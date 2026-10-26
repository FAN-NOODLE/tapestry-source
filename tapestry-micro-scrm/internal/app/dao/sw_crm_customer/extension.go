package swCrmCustomerDao

import (
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	"micro-scrm/internal/app/models/sw_crm_customer"
)

type CrmCustomerExtensionDao struct {
	dao.EnterpriseRepo[swCrmCustomerModel.CrmCustomerExtensionModel]
	model     *swCrmCustomerModel.CrmCustomerExtensionModel
	tableName string
}

func NewExtensionDao() *CrmCustomerExtensionDao {
	return &CrmCustomerExtensionDao{
		EnterpriseRepo: dao.EnterpriseRepo[swCrmCustomerModel.CrmCustomerExtensionModel]{
			Dao: dao.Dao{DB: app.DB()},
		},
		model:     new(swCrmCustomerModel.CrmCustomerExtensionModel),
		tableName: swCrmCustomerModel.CrmCustomerExtensionModel{}.TableName(),
	}
}
