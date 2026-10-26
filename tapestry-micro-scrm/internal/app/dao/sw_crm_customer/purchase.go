package swCrmCustomerDao

import (
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	"micro-scrm/internal/app/models/sw_crm_customer"
)

type CrmCustomerPurchaseDao struct {
	dao.EnterpriseRepo[swCrmCustomerModel.CrmCustomerPurchaseModel]
	model     *swCrmCustomerModel.CrmCustomerPurchaseModel
	tableName string
}

func NewPurchaseDao() *CrmCustomerPurchaseDao {
	return &CrmCustomerPurchaseDao{
		EnterpriseRepo: dao.EnterpriseRepo[swCrmCustomerModel.CrmCustomerPurchaseModel]{
			Dao: dao.Dao{DB: app.DB()},
		},
		model:     new(swCrmCustomerModel.CrmCustomerPurchaseModel),
		tableName: swCrmCustomerModel.CrmCustomerPurchaseModel{}.TableName(),
	}
}
