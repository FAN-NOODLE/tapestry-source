package micro_scrm

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
	CreatedAt    int64  `json:"created_at" db:"created_at"`
	UpdatedAt    int64  `json:"updated_at" db:"updated_at"`
	DeletedAt    int64  `json:"deleted_at" db:"deleted_at"`
}

func (w WorkGroupChatMember) TableName() string {
	return "work_group_chat_member"
}
