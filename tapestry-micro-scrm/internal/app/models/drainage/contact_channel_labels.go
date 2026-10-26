package drainageContactModel

import (
	"micro-scrm/internal/app/models"
	"micro-scrm/internal/orm/model/micro_scrm"
)

type DrainageContactChannelLabelsModel struct {
	models.MicroBaseModel                                 // 假设包含 id, enterprise_id, created_at, updated_at, deleted_at
	ChannelID             int64                           `gorm:"column:channel_id;type:bigint(20);not null;index:channel_id;comment:渠道id" json:"channel_id"`
	LabelGroupID          int64                           `gorm:"column:label_group_id;type:bigint(20);not null;index:label_group_id;comment:标签分组id" json:"label_group_id"`
	LabelID               int64                           `gorm:"column:label_id;type:bigint(20);not null;comment:标签id" json:"label_id"`
	WxTagGroupID          string                          `gorm:"column:wx_tag_group_id;type:varchar(100);not null;comment:企微标签分组ID" json:"wx_tag_group_id"`
	WxTagID               string                          `gorm:"column:wx_tag_id;type:varchar(100);not null;comment:企微标签ID" json:"wx_tag_id"`
	Tags                  micro_scrm.WorkExternalTagModel `json:"tags" gorm:"foreignKey:ID;references:LabelID"`
}

// TableName returns the table name of the DisplayTask model
func (DrainageContactChannelLabelsModel) TableName() string {
	return "drainage_contact_channel_labels"
}
