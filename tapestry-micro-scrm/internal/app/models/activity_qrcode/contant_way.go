package activityQRCodeModel

import (
	"micro-scrm/internal/app/models"
)

type ActivityQRCodeContactWayModel struct {
	models.BaseModel
	ActivityID int64  `json:"activity_id" gorm:"column:activity_id;type:bigint(20);not null;default:0;index:activity_id_index;comment:活动id"`
	State      string `json:"state" gorm:"column:state;type:varchar(255);not null;default:'';comment:点位id(联系我二维码state)"`
	ThirdId    int    `json:"third_id" gorm:"column:third_id;type:bigint(20);not null;default:0;comment:门店id"`
	StoreCode  string `json:"store_code" gorm:"column:store_code;type:varchar(255);not null;default:'';comment:门店编码"`
	StoreName  string `json:"store_name" gorm:"column:store_name;type:varchar(255);not null;default:'';comment:门店名称"`
	ConfigId   string `json:"config_id" gorm:"column:config_id;type:varchar(255);not null;default:'';comment:联系我配置id"`
	QrCode     string `json:"qr_code" gorm:"column:qr_code;type:varchar(255);not null;default:'';comment:联系我二维码"`
}

func (ActivityQRCodeContactWayModel) TableName() string {
	return "activity_qrcode_contact_way"
}
