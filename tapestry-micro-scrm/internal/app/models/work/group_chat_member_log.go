package workModels

import "micro-scrm/internal/app/models"

type WorkGroupChatMemberLogModel struct {
	models.MicroBaseModel
	GroupChatID     int64  `json:"group_chat_id" gorm:"column:group_chat_id;type:bigint(20);not null;index:chat_id_index,chat_union_index;comment:群编号,work_group_chat.id"`
	UnionID         string `json:"union_id" gorm:"column:union_id;type:varchar(255);not null;index:chat_union_index,unionid_index;comment:用户unionid"`
	ExternalUserid  string `json:"external_userid" gorm:"column:external_userid;type:varchar(255);not null;comment:外部联系人userid"`
	OperationStatus int32  `json:"operation_status" gorm:"column:operation_status;type:tinyint(1);not null;index:enterprise_id_index;default:0;comment:操作类型：1加群、2离群"`
	OperationTime   int    `json:"operation_time" gorm:"column:operation_time;type:int(10);not null;default:0;comment:操作类型：1加群、2离群"`
	MemberType      int32  `json:"member_type" gorm:"column:member_type;type:tinyint(1);not null;index:enterprise_id_index;default:0;comment:群成员类型：1企业成员、2外部联系人"`
	JoinScene       int32  `json:"join_scene" gorm:"column:join_scene;type:tinyint(1);not null;default:0;comment:加群场景：1 - 由群成员邀请入群（直接邀请入群）2 - 由群成员邀请入群（通过邀请链接入群）3 - 通过扫描群二维码入群 0 表示非加群操作"`
	GroupNickname   string `json:"group_nickname" gorm:"column:group_nickname;type:varchar(255);not null;comment:在群里的昵称"`
	Name            string `json:"name" gorm:"column:name;type:varchar(255);not null;comment:用户名称(名字。仅当 need_name = 1 时返回如果是微信用户，则返回其在微信中设置的名字如果是企业微信联系人，则返回其设置对外展示的别名或实名)"`
	InvitorUserid   string `json:"invitor_userid" gorm:"column:invitor_userid;type:varchar(255);not null;comment:邀请者的userid"`
}

// TableName returns the table name of the WorkGroupChatMemberLog model
func (WorkGroupChatMemberLogModel) TableName() string {
	return "work_group_chat_member_log"
}

const (
	JOINGROUPSTATUS = iota + 1 //操作类型，1加群，2离群
	LEAVEGROUPSTATUS
)

const (
	MEMBERTYPEWORK = iota + 1 //群成员类型：1企业成员、2外部联系人
	MEMBERTYPEWECHAT
)
