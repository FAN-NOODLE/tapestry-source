package work_external_user

import (
	"micro-scrm/internal/orm/model"
)

type MemberExternalSyncTagLogModel struct {
	model.MicroScrmBaseModel

	MemberID    int64 `gorm:"column:member_id;type:bigint(11);not null;"`
	ExternalNum int   `gorm:"column:external_num;type:int(11);not null;"`
}

// TableName returns the table name of the struct_name model
func (w MemberExternalSyncTagLogModel) TableName() string {
	return "member_external_sync_tag_log"
}

func (w MemberExternalSyncTagLogModel) DBPrefix() string {
	return "micro_scrm"
}

