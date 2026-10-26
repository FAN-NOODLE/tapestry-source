package microScrmModel

import (
	"micro-scrm/internal/app/models"
)

type NotificationInfoMemberModel struct {
	models.BaseModel
	InfoId   int64 `json:"info_id" gorm:"column:info_id;type:bigint(20);not null;default:0;comment:传讯id"`
	MemberId int64 `json:"member_id" gorm:"column:member_id;type:bigint(20);not null;default:0;comment:导购id"`
	ViewNum  int   `json:"view_num" gorm:"column:view_num;type:int(10);not null;default:0;comment:查看次数"`
}

// TableName returns the table name of the NotificationInfoMerge model
func (NotificationInfoMemberModel) TableName() string {
	return "notification_info_member"
}
