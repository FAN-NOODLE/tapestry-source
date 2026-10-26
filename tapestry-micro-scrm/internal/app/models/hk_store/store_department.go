package hkStoreModel

import "micro-scrm/internal/app/models"

type StoreDepartmentModel struct {
	models.MicroBaseModel

	StoreCode          string `json:"store_code" gorm:"column:store_code;type:varchar(32);not null;uniqueIndex:unique_store_code;comment:store code"`
	WechatDepartmentID uint   `json:"wechat_department_id" gorm:"column:wechat_department_id;type:int(10) unsigned;not null;default:0;comment:wechat department id"`
	StoreName          string `json:"store_name" gorm:"column:store_name;type:varchar(255);not null;default:'';comment:store name"`
	Area               string `json:"area" gorm:"column:area;type:varchar(10);not null;default:'';comment:store area hk/mo"`
	StoreType          int32  `json:"store_type" gorm:"column:store_type;type:tinyint(1);not null;default:0;comment:store type"`
	InventoryName      string `json:"inventory_name" gorm:"column:inventory_name;type:varchar(255);not null;default:'';comment:inventory name"`
	Address            string `json:"address" gorm:"column:address;type:varchar(255);not null;default:'';comment:address"`
	City               string `json:"city" gorm:"column:city;type:varchar(50);not null;default:'';comment:city"`
	Phone              string `json:"phone" gorm:"column:phone;type:varchar(50);not null;default:'';comment:phone"`
	Email              string `json:"email" gorm:"column:email;type:varchar(255);not null;default:'';comment:email"`
	ClosedAt           int64  `json:"closed_at" gorm:"column:closed_at;type:int(10);not null;default:0;comment:closed time"`
}

func (StoreDepartmentModel) TableName() string {
	return "hkmo_store_department"
}
