package drainageContactModel

import "micro-scrm/internal/app/models"

type DrainageContactsWelcomesModel struct {
	models.BaseModel
	ContactGroupId int64  `json:"contact_group_id" gorm:"column:contact_group_id;type:bigint(20);not null;default:0;index:rel_id_index;comment:联系我活动id"`
	Content        string `json:"content" gorm:"column:content;type:text;not null;comment:内容"`
	Type           int    `json:"type" gorm:"column:type;type:tinyint(1);not null;default:0;comment:0 通用 1 会员 2非会员"`
}

func (DrainageContactsWelcomesModel) TableName() string {
	return "drainage_contact_group_welcomes"
}

const WelcomeTypeHasCardLevel = 1 //会员欢迎语
const WelcomeTypeNoCardLevel = 2  //非会员欢迎语
