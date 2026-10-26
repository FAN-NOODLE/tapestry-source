package microScrmModel

import (
	"micro-scrm/internal/app/models"
)

type NotificationNoticeConfigModel struct {
	InfoType     int              `json:"info_type" gorm:"column:info_type;type:int(10);not null;default:0;comment:传讯类型"`
	DeadlineWeek int              `json:"deadline_week" gorm:"column:deadline_week;type:int(10);not null;default:0;comment:截止星期几（1-7）"`
	DeadlineTime string           `json:"deadline_time" gorm:"column:deadline_time;type:varchar(255);not null;comment:截止时间（时分秒）"`
	ReleaseWeek  int              `json:"release_week" gorm:"column:release_week;type:int(10);not null;default:0;comment:发布星期几（1-7）"`
	ReleaseTime  string           `json:"release_time" gorm:"column:release_time;type:varchar(255);not null;comment:发布时间（时分秒）"`
	CreatedAt    int64            `gorm:"column:created_at;type:int(10);not null;default:0" json:"created_at"`
	UpdatedAt    int64            `gorm:"column:updated_at;type:int(10);not null;default:0" json:"updated_at"`
	DeletedAt    models.DeletedAt `gorm:"softDelete;column:deleted_at;type:int(10) unsigned;not null;default:0" json:"deleted_at"`
}

// TableName returns the table name of the NotificationInfoMerge model
func (NotificationNoticeConfigModel) TableName() string {
	return "notification_notice_config"
}
