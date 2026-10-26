package fission

import "micro-scrm/internal/orm/model"

type FissionChatGroupModel struct {
	model.MicroScrmBaseModel

	ActivityID int64  `gorm:"column:activity_id;type:bigint(20);not null;index:enterprise_activity_deleted_at_index;default:0;comment:活动id"`
	ChatID     int64  `gorm:"column:chat_id;type:bigint(20);not null;index:chat_id_idx;default:0;comment:群聊id"`
	QrCode     string `gorm:"column:qr_code;type:varchar(500);not null;comment:入群二维码"`
	Sort       int    `gorm:"column:sort;type:int(10);not null;default:0;comment:排序权重，升序"`
}

func (FissionChatGroupModel) TableName() string {
	return "fission_chat_group"
}
