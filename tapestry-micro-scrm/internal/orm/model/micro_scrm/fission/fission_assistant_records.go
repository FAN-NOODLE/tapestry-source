package fission

import "micro-scrm/internal/orm/model"

const (
	FAR_IsVip_No = iota
	FAR_IsVip_Yes
)

type FissionAssistantRecordsModel struct {
	model.MicroScrmBaseModel

	ActivityID          int64  `gorm:"column:activity_id;type:bigint(20);not null;index:deleted_enterprise_activity_index;default:0;comment:活动id"`
	FissionUserID       int64  `gorm:"column:fission_user_id;type:bigint(20);not null;index:fission_user_id_index;default:0;comment:裂变用户关系表id"`
	AssistantUserID     int64  `gorm:"column:assistant_user_id;type:bigint(20);not null;index:assistant_user_member_idx;default:0;comment:助力用户id,work_external_user表user_id"`
	ExternalUserID      string `gorm:"column:external_user_id;type:varchar(255);not null;comment:外部联系人userid"`
	AddMemberStatus     int32  `gorm:"column:add_member_status;type:tinyint(1);not null;default:0;comment:添加好友记录：0-未添加，1-已添加，2-已删除"`
	AddMemberAt         uint   `gorm:"column:add_member_at;type:int(10) unsigned;not null;default:0;comment:好友添加时间"`
	JoinChatGroupStatus int32  `gorm:"column:join_chat_group_status;type:tinyint(1);not null;default:0;comment:入群状态：0-未入群，1-已入群，2-已退群"`
	JoinChatGroupAt     uint   `gorm:"column:join_chat_group_at;type:int(10) unsigned;not null;default:0;comment:入群时间"`
	MemberID            uint64 `gorm:"column:member_id;type:bigint(20) unsigned;not null;index:assistant_user_member_idx;default:0;comment:添加的员工ID"`
	Score               int    `gorm:"column:score;type:int(10);not null;default:0;comment:计分"`
	IsVip               int    `gorm:"column:is_vip;type:tinyint(1) unsigned;not null;default:0;comment:vip状态：0-非vip，1-vip"`
}

func (FissionAssistantRecordsModel) TableName() string {
	return "fission_assistant_records"
}
