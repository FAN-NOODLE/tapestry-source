package group_chat

import "micro-scrm/internal/orm/model"

const (
	WGC_STATUS_YES = iota + 1
	WGC_STATUS_NO
)

type WorkGroupChatModel struct {
	model.BaseModel

	EnterpriseID int64 `json:"enterprise_id" gorm:"column:enterprise_id;type:bigint(20);not null;index:enterprise_id_chat_id_index;default:0"`

	ChatID        string `json:"chat_id" gorm:"column:chat_id;type:varchar(255);not null;index:enterprise_id_chat_id_index;comment:企业微信群聊id"`
	Name          string `json:"name" gorm:"column:name;type:varchar(255);not null;comment:群聊名称"`
	Type          int32  `json:"type" gorm:"column:type;type:tinyint(1);not null;default:1;comment:群类型 1:外部群 2内部群"`
	OwnerMemberID string `json:"owner_member_id" gorm:"column:owner_member_id;type:varchar(255);not null;index:owner_member_id_index;default:0;comment:群创建人id(对应members表userid)"`
	Notice        string `json:"notice" gorm:"column:notice;type:text;not null;comment:群公告"`
	CreateTime    int    `json:"create_time" gorm:"column:create_time;type:int(10);not null;default:0;comment:创建时间"`
	Status        int32  `json:"status" gorm:"column:status;type:tinyint(1);not null;comment:群状态 1=群存在 2=群已解散"`
	DissolveTime  int    `json:"dissolve_time" gorm:"column:dissolve_time;type:int(10);not null;default:0;comment:群解散时间"`
	MemberNum     uint   `json:"member_num" gorm:"column:member_num;type:int(10) unsigned;not null;default:0;comment:群人数"`
}

func (WorkGroupChatModel) TableName() string {
	return "work_group_chat"
}
