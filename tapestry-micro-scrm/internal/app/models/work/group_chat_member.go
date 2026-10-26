package workModels

import "micro-scrm/internal/app/models"

type WorkGroupChatMemberModel struct {
	models.MicroBaseModel
	GroupChatID  int64  `json:"group_chat_id" gorm:"column:group_chat_id;type:bigint(20);not null;index:group_chat_id_index;default:0;comment:群聊表id"`
	Type         int32  `json:"type" gorm:"column:type;type:tinyint(1);not null;index:type_userid_index;default:1;comment:成员类型:1 - 企业成员,2 - 外部联系人"`
	Userid       string `json:"userid" gorm:"column:userid;type:varchar(64);not null;index:type_userid_index,userid;default:0;comment:群成员id"`
	MemberID     int64  `json:"member_id" gorm:"column:member_id;type:bigint(20);not null;index:member_idx;default:0;comment:内部成员ID"`
	Unionid      string `json:"unionid" gorm:"column:unionid;type:varchar(32);not null;index:unionid_index;comment:微信unionid"`
	Status       int32  `json:"status" gorm:"column:status;type:tinyint(1);not null;index:enterprise_id,status_index;default:0;comment:群状态，1=入群 2=退群"`
	JoinTime     int    `json:"join_time" gorm:"column:join_time;type:int(10);not null;index:join_time_index;default:0;comment:入群时间"`
	LeaveTime    int    `json:"leave_time" gorm:"column:leave_time;type:int(10);not null;index:leave_time_index;default:0;comment:退群时间"`
	JoinScene    int    `json:"join_scene" gorm:"column:join_scene;type:int(10);not null;default:1;comment:入群方式。1 - 由成员邀请入群（直接邀请入群）2 - 由成员邀请入群（通过邀请链接入群）3 - 通过扫描群二维码入群"`
	ActNum       int    `json:"act_num" gorm:"column:act_num;type:int(10);not null;default:0;comment:群成员用户参与活动的数量"`
	ActIds       string `json:"act_ids" gorm:"column:act_ids;type:varchar(2048);not null;default:[];comment:群成员用户参与活动的ID集合"`
	Nickname     string `json:"nickname" gorm:"column:nickname;type:varchar(255);not null;comment:群名称"`
	Name         string `json:"name" gorm:"column:name;type:varchar(255);not null;comment:群名称"`
	InviteUserid string `json:"invite_userid" gorm:"column:invite_userid;type:varchar(255);not null;comment:邀请人ID"`
}

// TableName returns the table name of the WorkGroupChatMember model
func (WorkGroupChatMemberModel) TableName() string {
	return "work_group_chat_member"
}
