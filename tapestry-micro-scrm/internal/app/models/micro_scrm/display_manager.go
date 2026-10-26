package microScrmModel

import (
	"micro-scrm/internal/app/models"
)

type DisplayManagerModel struct {
	models.BaseModel
	MemberID    int64  `json:"member_id" gorm:"column:member_id;type:bigint(20);not null;default:0"`
	UserId      string `json:"user_id" gorm:"column:user_id;type:varchar(255);not null;comment:导购的工号，user_id"`
	IsCheck     int    `json:"is_check" gorm:"column:is_check;type:tinyint(1);not null;default:0;comment:是否为VM审核人员,0为否,1为是"`
	IsManager   int    `json:"is_manager" gorm:"column:is_manager;type:tinyint(1);not null;default:0;comment:是否为管理者,0为否,1为是"`
	IsSpecial   int    `json:"is_special" gorm:"column:is_special;type:tinyint(1);not null;default:0;comment:是否为special人员：0否，1是"`
	IsCorporate int    `json:"is_corporate" gorm:"column:is_corporate;type:tinyint(1);not null;default:0;comment:是否为corporate人员：0否，1是"`
}

// TableName returns the table name of the DisplayManager model
func (DisplayManagerModel) TableName() string {
	return "display_manager"
}
