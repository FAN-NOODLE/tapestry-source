package microScrmModel

import (
	"micro-scrm/internal/app/models"
)

type NotificationCrowdDataModel struct {
	models.BaseModel
	CrowdId  int64 `json:"crowd_id" gorm:"column:crowd_id;type:bigint(20);not null;default:0;comment:crowd表主键id"`
	MemberId int64 `json:"member_id" gorm:"column:member_id;type:bigint(20);not null;default:0;comment:member_id"`
}

// TableName returns the table name of the NotificationCorwd model
func (NotificationCrowdDataModel) TableName() string {
	return "notification_crowd_data"
}
