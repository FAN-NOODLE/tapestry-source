package swCrmCustomerModel

import (
	"micro-scrm/internal/app/models"
)

type CrmCustomerExtensionModel struct {
	models.BaseModel
	CustomerID        string  `json:"customer_id" gorm:"column:customer_id;type:char(64);not null;default:'';comment:customer_id"`
	UnionID           string  `json:"union_id" gorm:"column:union_id;type:char(64);not null;default:'';comment:union_id"`
	PriceTag          string  `json:"price_tag" gorm:"column:price_tag;type:char(64);not null;default:'';comment:price_tag"`
	AgeTag            string  `json:"age_tag" gorm:"column:age_tag;type:char(64);not null;default:'';comment:age_tag"`
	R12PurchaseAmount float64 `json:"r12_purchase_amount" gorm:"column:r12_purchase_amount;type:decimal(10,2);not null;default:0.00;comment:r12_purchase_amont"`
	AtValues          string  `json:"at_values" gorm:"column:at_values;type:varchar(255);not null;default:'';comment:at_values"`
}

func (CrmCustomerExtensionModel) TableName() string {
	return "crm_customer_extension"
}
