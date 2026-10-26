package activityQRCodeModel

import (
	"micro-scrm/internal/app/models"
)

type ActivityQRCodeWelcomeModel struct {
	models.BaseModel
	RelType int    `json:"rel_type" gorm:"column:rel_type;type:tinyint(1);not null;default:0;comment:关联id类型 1 活动id 2渠道id"`
	RelID   int64  `json:"rel_id" gorm:"column:rel_id;type:bigint(20);not null;default:0;index:rel_id_index;comment:相关id"`
	Content string `json:"content" gorm:"column:content;type:text;not null;comment:内容"`
	Type    int    `json:"type" gorm:"column:type;type:tinyint(1);not null;default:0;comment:0 通用 1 会员 2非会员"`
}

func (ActivityQRCodeWelcomeModel) TableName() string {
	return "activity_qrcode_welcomes"
}

const WelcomeTypeActivity = 1 //活动欢迎语
const WelcomeTypeCategory = 2 //投放渠道欢迎语

const WelcomeTypeHasCardLevel = 1 //会员欢迎语
const WelcomeTypeNoCardLevel = 2  //非会员欢迎语
