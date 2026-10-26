package drainageContactModel

import (
	"micro-scrm/internal/app/models"
)

type DrainageContactsMemberQrcodeModel struct {
	models.MicroBaseModel
	ContactId int64  `gorm:"column:contact_id;type:bigint(20);not null;default:0;comment:联系我表主键"`
	MemberId  int64  `gorm:"column:member_id;type:bigint(20);not null;default:0;comment:导购id"`
	Userid    string `gorm:"column:userid;type:varchar(64);not null;comment:员工企微账号"`
	QrCodeURL string `gorm:"column:qr_code_url;type:varchar(1000);not null;comment:二维码地址"`
	ConfigID  string `gorm:"column:config_id;type:varchar(255);not null;comment:微信返回配置id"`
}

// TableName returns the table name of the DisplayTask model
func (DrainageContactsMemberQrcodeModel) TableName() string {
	return "drainage_contact_member_qrcode"
}
