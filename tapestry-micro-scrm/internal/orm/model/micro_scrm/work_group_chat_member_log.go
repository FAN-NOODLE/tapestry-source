package micro_scrm

import "micro-scrm/internal/orm/model"

const (
	WGCML_MEMBER_TYPE        = 2
	WGCML_OPERATION_STATUS   = 1
	EXTERNAL_CONTACTS_STRING = "wm%"
)

type WorkGroupChatMemberLog struct {
	model.MicroScrmBaseModel
	GroupChatId     int64  `gorm:"column:group_chat_id;NOT NULL" json:"group_chat_id"`                 // 群编号,work_group_chat.id
	UnionId         string `gorm:"column:union_id;NOT NULL" json:"union_id"`                           // 用户unionid
	ExternalUserid  string `gorm:"column:external_userid;NOT NULL" json:"external_userid"`             // 外部联系人userid
	OperationStatus int    `gorm:"column:operation_status;default:0;NOT NULL" json:"operation_status"` // 操作类型：1加群、2离群
	OperationTime   int    `gorm:"column:operation_time;default:0;NOT NULL" json:"operation_time"`     // 操作类型：1加群、2离群
	MemberType      int    `gorm:"column:member_type;default:0;NOT NULL" json:"member_type"`           // 群成员类型：1企业成员、2外部联系人
	JoinScene       int    `gorm:"column:join_scene;default:0;NOT NULL" json:"join_scene"`             // 加群场景：1 - 由群成员邀请入群（直接邀请入群）2 - 由群成员邀请入群（通过邀请链接入群）3 - 通过扫描群二维码入群 0 表示非加群操作
	GroupNickname   string `gorm:"column:group_nickname;NOT NULL" json:"group_nickname"`               // 群昵称
	Name            string `gorm:"column:name;NOT NULL" json:"name"`                                   // 用户名称(名字。仅当 need_name = 1 时返回如果是微信用户，则返回其在微信中设置的名字如果是企业微信联系人，则返回其设置对外展示的别名或实名)
	InvitorUserid   string `gorm:"column:invitor_userid;NOT NULL" json:"invitor_userid"`               // 邀请者的userid

}

func (w WorkGroupChatMemberLog) TableName() string {
	return "work_group_chat_member_log"
}
