package microScrmModel

import (
	"micro-scrm/internal/app/models"
)

type NotificationManagerModel struct {
	models.BaseModel
	MemberID int64 `json:"member_id" gorm:"column:member_id;type:bigint(20);not null;default:0"`
	AdminID  int64 `json:"admin_id" gorm:"column:admin_id;type:bigint(20);not null;default:0"`
	IsFt     int32 `json:"is_ft" gorm:"column:is_ft;type:tinyint(1);not null;default:0;comment:是否function_team，0否,1是"`
	IsRo     int32 `json:"is_ro" gorm:"column:is_ro;type:tinyint(1);not null;default:0;comment:是否ro人员，0否，1是"`
}

// TableName returns the table name of the NotificationManager model
func (NotificationManagerModel) TableName() string {
	return "notification_manager"
}

const NotificationIsFt = 1
const NotificationIsRo = 1
