package swCrmCustomerModel

import (
	"micro-scrm/internal/app/models"
)

type CrmCustomerPurchaseModel struct {
	models.BaseModel
	CustomerID   string `json:"customer_id" gorm:"column:customer_id;type:char(64);not null;default:'';comment:customer_id"`
	UnionID      string `json:"union_id" gorm:"column:union_id;type:char(64);not null;default:'';comment:union_id"`
	SaCode       string `json:"sa_code" gorm:"column:sa_code;type:char(64);not null;default:'';comment:sa_code"`
	CategoryCode string `json:"category_code" gorm:"column:category_code;type:varchar(255);not null;default:'';comment:category_code"`
	PurchaseCode string `json:"purchase_code" gorm:"column:purchase_code;type:char(64);not null;default:'';comment:purchase_code"`
	PurchaseTime string `json:"purchase_time" gorm:"column:purchase_time;type:char(64);not null;default:'';comment:purchase_time"`
}

func (CrmCustomerPurchaseModel) TableName() string {
	return "crm_customer_purchase"
}
