package activityQRCodeModel

import (
	"micro-scrm/internal/app/models"
)

type ActivityQRCodePointModel struct {
	models.BaseModel
	ActivityID int64  `json:"activity_id" gorm:"column:activity_id;type:bigint(20);not null;default:0;index:activity_id_index;comment:活动id"`
	Name       string `json:"name" gorm:"column:name;type:varchar(255);not null;default:'';comment:点位名称"`
	IsPoint    int    `json:"is_point" gorm:"column:is_point;type:tinyint(1);not null;default:0;comment:0 活动码 1 活动点位码"`
	CategoryID int64  `json:"category_id" gorm:"column:category_id;type:bigint(20);not null;default:0;comment:投放渠道id"`
	CodeURL    string `json:"code_url" gorm:"column:code_url;type:varchar(255);not null;default:'';comment:二维码图片url"`
	LinkURL    string `json:"link_url" gorm:"column:link_url;type:varchar(255);not null;default:'';comment:链接url"`
}

func (ActivityQRCodePointModel) TableName() string {
	return "activity_qrcode_point"
}
