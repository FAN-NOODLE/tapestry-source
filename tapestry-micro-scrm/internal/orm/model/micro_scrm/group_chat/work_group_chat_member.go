package group_chat

import "micro-scrm/internal/orm/model"

type WorkGroupChatMemberModel struct {
	model.MicroScrmBaseModel

	UserID      string             `json:"userid" gorm:"column:userid;type:varchar(64);not null;default:0;comment:群成员id"`
	UnionID     string             `json:"unionid" gorm:"column:unionid;type:varchar(32);not null;comment:群成员id"`
	GroupChatID int64              `json:"group_chat_id" gorm:"column:group_chat_id;type:bigint(20);not null;index:group_chat_id_index;default:0"`
	Type        int32              `json:"type" gorm:"column:type;type:tinyint(1);not null;default:1;comment:成员类型:1 - 企业成员,2 - 外部联系人"`
	MemberID    int64              `json:"member_id" gorm:"column:member_id;type:bigint(20);not null;;default:0;comment:内部成员ID"`
	Status      int                `json:"status" gorm:"column:status;type:int(10);not null;index:status_index;default:0"`
	JoinTime    int64              `json:"join_time" gorm:"column:join_time;type:int(10);not null;index:leave_time_index;default:0"`
	LeaveTime   int64              `json:"leave_time" gorm:"column:leave_time;type:int(10);not null;index:join_time_index;default:0"`
	JoinScene   int                `json:"join_scene" gorm:"column:join_scene;type:int(10);not null;default:0;comment:入群方式。\n1 - 由成员邀请入群（直接邀请入群）\n2 - 由成员邀请入群（通过邀请链接入群）\n3 - 通过扫描群二维码入群"`
	GroupChat   WorkGroupChatModel `gorm:"foreignKey:GroupChatID"`
}

func (WorkGroupChatMemberModel) TableName() string {
	return "work_group_chat_member"
}
