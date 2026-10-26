package onlineWearhouseModel

import "micro-scrm/internal/app/models"

type OnlineWarehouseStoreModel struct {
	StoreName string `json:"store_name" gorm:"column:store_name;type:varchar(255);not null;comment:门店名称"`
	StoreCode string `json:"store_code" gorm:"column:store_code;type:varchar(255);not null;comment:门店code"`
	ThirdId   int    `json:"third_id" gorm:"column:third_id;type:int(10);not null;default:0;comment:企微部门id"`
	models.BaseModel
}

// TableName returns the table name of the OnlineWarehouseOrder model
func (OnlineWarehouseStoreModel) TableName() string {
	return "online_warehouse_store"
}
