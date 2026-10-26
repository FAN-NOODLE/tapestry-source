package microScrmModel

import (
	"micro-scrm/internal/app/models"
)

type NotificationAllotModel struct {
	models.BaseModel
	NoticeID   int64  `json:"notice_id" gorm:"column:notice_id;type:bigint(20);not null;default:0"`
	CrowdID    int64  `json:"crowd_id" gorm:"column:crowd_id;type:bigint(20);not null;default:0"`
	UserID     string `json:"user_id" gorm:"column:user_id;type:varchar(255);not null;comment:导购工号"`
	MemberID   int64  `json:"member_id" gorm:"column:member_id;type:bigint(20);not null;default:0;comment:导购的member_id"`
	AgentMsgID string `json:"agent_msg_id" gorm:"column:agent_msg_id;type:varchar(255);not null;comment:应用msg_id"`
	LookNum    int    `json:"look_num" gorm:"column:look_num;type:int(10);not null;default:0;comment:查看次数"`
	IsLook     int32  `json:"is_look" gorm:"column:is_look;type:tinyint(1);not null;default:0;comment:是否查看"`
}

// TableName returns the table name of the NotificationAllot model
func (NotificationAllotModel) TableName() string {
	return "notification_allot"
}
