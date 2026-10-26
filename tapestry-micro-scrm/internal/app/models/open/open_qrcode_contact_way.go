package openModel

import "micro-scrm/internal/app/models"

type OpenQrcodeContactWayModel struct {
	models.BaseModel
	State    string `json:"state" gorm:"column:state;type:varchar(255);not null;default:'';comment:联系我二维码state"`
	ThirdId  int    `json:"third_id" gorm:"column:third_id;type:bigint(20);not null;default:0;comment:门店id"`
	UserId   string `json:"user_id" gorm:"column:user_id;type:varchar(255);not null;default:'';comment:企微用户id"`
	ConfigId string `json:"config_id" gorm:"column:config_id;type:varchar(255);not null;default:'';comment:联系我配置id"`
	QrCode   string `json:"qr_code" gorm:"column:qr_code;type:varchar(255);not null;default:'';comment:联系我二维码"`
}

// TableName returns the table name of the OpenUser model
func (OpenQrcodeContactWayModel) TableName() string {
	return "open_qrcode_contact_way"
}
