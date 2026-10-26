package microScrmModel

import "micro-scrm/internal/app/models"

type InventorySkuUpcMappingModel struct {
	models.BaseModel
	BatchNo        string `json:"batch_no" gorm:"column:batch_no;type:varchar(64);not null;default:'';comment:import batch number"`
	SkuItem        string `json:"sku_item" gorm:"column:sku_item;type:varchar(64);not null;default:'';comment:SKU item"`
	UpcCode        string `json:"upc_code" gorm:"column:upc_code;type:varchar(32);not null;default:'';comment:UPC code"`
	EffectiveMonth string `json:"effective_month" gorm:"column:effective_month;type:char(7);not null;default:'';comment:YYYY-MM"`
}

func (InventorySkuUpcMappingModel) TableName() string {
	return "inventory_sku_upc_mapping"
}
