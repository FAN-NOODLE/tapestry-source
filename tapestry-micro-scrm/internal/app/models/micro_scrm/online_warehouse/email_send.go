package onlineWearhouseModel

import "micro-scrm/internal/app/models"

type OnlineWarehouseEmailSendModel struct {
	Email string `json:"email" gorm:"column:email;type:varchar(255);not null;comment:邮箱地址"`
	Type  int32  `json:"type" gorm:"column:type;type:tinyint(1);not null;default:1;comment:1 网仓"`
	models.BaseModel
}

// TableName returns the table name of the OnlineWarehouseOrder model
func (OnlineWarehouseEmailSendModel) TableName() string {
	return "online_warehouse_email_send"
}
