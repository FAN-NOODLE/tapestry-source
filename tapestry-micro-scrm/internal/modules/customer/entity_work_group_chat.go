package customer

import "micro-scrm/internal/pkg/sdk/work"

// WorkChatGroupMemberTypeQw 内部成员
const WorkChatGroupMemberTypeQw = 1

// WorkChatGroupMemberTypeWx 外部联系人（企微/微信）
const WorkChatGroupMemberTypeWx = 2

// WorkGroupChatMemberJoinStatus 群成员入群
const WorkGroupChatMemberJoinStatus = 1

// WorkGroupChatMemberLeaveStatus 群成员退群
const WorkGroupChatMemberLeaveStatus = 2

// WorkGroupChatInnerGroup 内部群
const WorkGroupChatInnerGroup = 1

// WorkGroupTypeExists 群存在
const WorkGroupTypeExists = 1

// WorkGroupTypeDissolve 群解散
const WorkGroupTypeDissolve = 2

// GroupChatNoticeLimit 群公告长度限制
const GroupChatNoticeLimit = 20000

// WorkGroupChat 客户群信息
type WorkGroupChat struct {
	ID               int64  `json:"id" db:"id" desc:"主键"`
	EnterpriseID     int64  `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	Type             int    `json:"type" db:"type" desc:"群类型"`
	ChatID           string `json:"chat_id" db:"chat_id" desc:"群ID"`
	Name             string `json:"name" db:"name" desc:"群名称"`
	OwnerMemberID    int64  `json:"owner_member_id" db:"owner_member_id" desc:"群主ID"`
	Notice           string `json:"notice" db:"notice" desc:"群公告"`
	CreateTime       int64  `json:"create_time" db:"create_time" desc:"群创建时间"`
	Status           int    `json:"status" db:"status" desc:"群状态 1=群存在 2=群已解散"`
	DissolveTime     int64  `json:"dissolve_time" db:"dissolve_time" desc:"群解散时间"`
	CreatedAt        int64  `json:"created_at" db:"created_at"`
	UpdatedAt        int64  `json:"updated_at" db:"updated_at"`
	DeletedAt        int64  `json:"deleted_at" db:"deleted_at"`
	Owner            string `json:"owner"`
	MemberNum        int    `json:"member_num" db:"member_num"`
	AdminList        string `db:"admin_list" desc:"管理员列表"`
	ChatTarget       string `json:"chat_target" db:"chat_target" desc:"群目标"`
	ChannelID        int64  `json:"channel_id" db:"channel_id" desc:"渠道ID"`
	ChatTypeID       int64  `json:"chat_type_id" db:"chat_type_id" desc:"群类型ID"`
	OperatorMemberID string `json:"operator_member_id" db:"operator_member_id" desc:"运营人员userid"`
}

// IsOpenChat 判断是不是外部群
func (c *WorkGroupChat) IsOpenChat() bool {
	return c.Type == 1
}

// FormatNotice 格式化群公告
func (c *WorkGroupChat) FormatNotice() {
	t := []rune(c.Notice)
	if len(t) > GroupChatNoticeLimit {
		c.Notice = string(t[:GroupChatNoticeLimit])
	}
}

// WorkGroupChatMember 群成员
type WorkGroupChatMember struct {
	ID           int64  `json:"id" db:"id" desc:"主键"`
	GroupChatID  int64  `json:"group_chat_id" db:"group_chat_id" desc:"work_group_chat 主表ID"`
	Type         int    `json:"type" db:"type" desc:"成员类型:1 - 企业成员,2 - 外部联系人"`
	UserID       string `json:"userid" db:"userid" desc:"群成员id"`
	JoinTime     int64  `json:"join_time" db:"join_time" desc:"加群时间"`
	JoinScene    int    `json:"join_scene" db:"join_scene" desc:"入群方式"`
	LeaveTime    int64  `json:"leave_time" db:"leave_time" desc:"退群时间"`
	Status       int    `json:"status" db:"status" desc:"群状态1=入群 2=退群"`
	MemberID     int64  `json:"member_id" db:"member_id" desc:"成员ID"`
	UnionID      string `json:"unionid" db:"unionid" desc:"微信unionid"`
	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id"`
	NickName     string `json:"-" db:"nickname" desc:"群昵称"`
	Name         string `json:"-" db:"name" desc:"名称"`
	InviteUserID string `json:"-" db:"invite_userid" desc:"邀请人ID"`
	CreatedAt    int64  `json:"created_at" db:"created_at"`
	UpdatedAt    int64  `json:"updated_at" db:"updated_at"`
	DeletedAt    int64  `json:"deleted_at" db:"deleted_at"`
}

// BatchWorkGroupMemberRequest 批量查询客户群成员
type BatchWorkGroupMemberRequest struct {
	EnterpriseID     int64   `json:"enterprise_id" db:"enterprise_id"`
	GroupChatID      int64   `json:"group_chat_id" db:"group_chat_id"`
	GroupChatIds     []int64 `json:"group_chat_ids"`
	GroupChatIdCount int64
	DeletedAt        int64 `json:"deleted_at" db:"deleted_at"`
	Page             int   `json:"page"`
	Size             int   `json:"size"`
}

// TransferWorkGroupChat 结构转换
func TransferWorkGroupChat(response *work.ExternalGroupResponse) *WorkGroupChat {
	return &WorkGroupChat{
		ChatID:     response.ChatID,
		Name:       response.Name,
		Notice:     response.Notice,
		CreateTime: response.CreateTime,
		Owner:      response.Owner,
	}
}

// TransferWorkGroupChatMember 群成员转换
func TransferWorkGroupChatMember(r *work.ExternalGroupMemberList) *WorkGroupChatMember {
	return &WorkGroupChatMember{
		Type:         r.Type,
		JoinTime:     r.JoinTime,
		JoinScene:    r.JoinScene,
		UnionID:      r.UnionID,
		UserID:       r.UserID,
		NickName:     r.GroupNickname,
		Name:         r.Name,
		InviteUserID: r.Invitor.UserID,
	}
}
