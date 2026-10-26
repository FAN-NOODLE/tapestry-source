package activityQRCodeModel

import (
	"micro-scrm/internal/app/models"
)

type ActivityQRCodeCategoryModel struct {
	models.BaseModel
	Name  string `json:"name" gorm:"column:name;type:varchar(255);not null;default:'';comment:名称"`
	TagId int64  `json:"tag_id" gorm:"column:tag_id;type:bigint(10);not null;default:0;index:parent_id_index;comment:渠道标签id"`
}

func (ActivityQRCodeCategoryModel) TableName() string {
	return "activity_qrcode_category"
}
