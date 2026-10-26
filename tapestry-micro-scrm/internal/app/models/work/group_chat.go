package workModels

import "micro-scrm/internal/app/models"

type GroupChatModel struct {
	models.MicroBaseModel
	ChatID           string `json:"chat_id" gorm:"column:chat_id;type:varchar(255);not null;index:enterprise_id_chat_id_index;comment:企业微信群聊id"`
	Name             string `json:"name" gorm:"column:name;type:varchar(255);not null;comment:群聊名称"`
	Type             int32  `json:"type" gorm:"column:type;type:tinyint(1);not null;default:1;comment:群类型 1:外部群 2内部群"`
	OwnerMemberID    int64  `json:"owner_member_id" gorm:"column:owner_member_id;type:bigint(20);not null;index:owner_member_id_index;default:0;comment:群创建人id"`
	Notice           string `json:"notice" gorm:"column:notice;type:text;not null;comment:群公告"`
	CreateTime       int    `json:"create_time" gorm:"column:create_time;type:int(10);not null;default:0;comment:创建时间"`
	Status           int32  `json:"status" gorm:"column:status;type:tinyint(1);not null;comment:群状态 1=群存在 2=群已解散"`
	DissolveTime     int    `json:"dissolve_time" gorm:"column:dissolve_time;type:int(10);not null;default:0;comment:群解散时间"`
	MemberNum        uint   `json:"member_num" gorm:"column:member_num;type:int(10) unsigned;not null;default:0;comment:群人数"`
	ChannelID        int64  `json:"channel_id" gorm:"column:channel_id;type:bigint(20);not null;index:idx_channel_id;default:0;comment:所属渠道ID"`
	ChatTypeID       int    `json:"chat_type_id" gorm:"column:chat_type_id;type:int(10);not null;index:idx_chat_type;default:0;comment:群类型ID"`
	ChatTarget       string `json:"chat_target" gorm:"column:chat_target;type:varchar(512);not null;comment:群目标"`
	OperatorMemberID string `json:"operator_member_id" gorm:"column:operator_member_id;type:varchar(64);not null;comment:群运营人员ID(对应members表的userid)"`
	AdminList        string `json:"admin_list" gorm:"column:admin_list;type:varchar(255);not null;default:[];comment:管理员userID列表"`
}

// TableName returns the table name of the WorkExternalUserTag model
func (GroupChatModel) TableName() string {
	return "work_group_chat"
}

type WorkGroupChatNumList struct {
	OwnerMemberID int64  `json:"owner_member_id"`
	Userid        string `json:"userid"`
	Num           int    `json:"num"`
}
