package micro_scrm

import "micro-scrm/internal/orm/model"

type UserEventModel struct {
	model.MicroScrmBaseModel
	UserId    int64  `gorm:"column:user_id;not null"`
	Event     string `gorm:"column:event;type:varchar(60)"`
	EventDesc string `gorm:"column:event_desc;type:varchar(255);"`
	Extra     string `gorm:"column:extra;default:'';type:varchar(255);"`
	Object    string `gorm:"column:object;type:varchar(255)"`
}

// TableName returns the table name of the struct_name model
func (UserEventModel) TableName() string {
	return "user_events_v1"
}
