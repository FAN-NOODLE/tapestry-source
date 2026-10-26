package sensitive

// TypeCart 敏感词类型 名片
const TypeCart = 1

// TypeText 敏感词类型 文本
const TypeText = 2

// UserChannelMember 消息接收人类型 1 员工
const UserChannelMember = 1

// ListRequest 获取敏感词列表请求参数
type ListRequest struct {
	SensitiveText string `json:"sensitive_text" db:"sensitive_text" desc:"敏感词搜索"`
	EnterpriseID  int64  `json:"enterprise_id" db:"enterprise_id" validate:"required" required:"企业ID不能为空" desc:"企业ID"`
	Page          int    `json:"page" desc:"页码"`
	Size          int    `json:"size" desc:"分页"`
}

// ResponseList 返回敏感词列表
type ResponseList struct {
	List  []Detail `json:"list"`
	Total int64    `json:"total"`
	Page  int      `json:"page"`
	Size  int      `json:"size"`
}

// Detail 敏感词字段
type Detail struct {
	ID              int64  `json:"id" db:"id" desc:"主键id"`
	SensitiveText   string `json:"sensitive_text" db:"sensitive_text" desc:"敏感词"`
	SensitiveType   int8   `json:"sensitive_type" db:"sensitive_type" desc:"敏感词类型 1 名片 2 文本"`
	EnterpriseID    int64  `json:"enterprise_id" db:"enterprise_id" desc:"企业id"`
	TriggerNum      int32  `json:"trigger_num" db:"trigger_num" desc:"总触发次数"`
	TodayTriggerNum int32  `json:"today_trigger_num" desc:"今日触发次数"`
	CreatedAt       int64  `json:"created_at" db:"created_at" desc:"创建时间"`
}

// CountBySensitive 获取敏感词每天对应的聊天数
type CountBySensitive struct {
	SensitiveID     int64 `json:"sensitive_id" db:"sensitive_id" desc:"敏感词id"`
	TodayTriggerNum int32 `json:"today_trigger_num" db:"today_trigger_num" desc:"今日触发次数"`
}

// SaveRequest 敏感词保存请求参数
type SaveRequest struct {
	ID            int64  `json:"id" db:"id" desc:"ID"`
	SensitiveType int8   `json:"sensitive_type" db:"sensitive_type" validate:"required" required:"敏感词类型不能为空" desc:"敏感词类型 1 名片 2 文本"`
	SensitiveText string `json:"sensitive_text" db:"sensitive_text" desc:"敏感词"`
	EnterpriseID  int64  `json:"enterprise_id" db:"enterprise_id" validate:"required" required:"企业ID不能为空" desc:"企业ID"`
	TriggerNum    int32  `json:"trigger_num" db:"trigger_num" desc:"总触发次数"`
}

// DeleteRequest 敏感词删除请求参数
type DeleteRequest struct {
	ID int64 `json:"id" db:"id"  validate:"required" required:"敏感词ID不能为空" desc:"ID"`
}

// DetailListRequest 敏感词明细列表请求参数
type DetailListRequest struct {
	ID   int64 `json:"id" db:"id"  validate:"required" required:"敏感词ID不能为空" desc:"ID"`
	Page int   `json:"page" desc:"页码"`
	Size int   `json:"size" desc:"分页"`
}

// Records 敏感词关联聊天记录字段
type Records struct {
	ID           int64 `json:"id" db:"id" desc:"主键id"`
	EnterpriseID int64 `json:"enterprise_id" db:"enterprise_id" validate:"required" required:"企业ID不能为空" desc:"企业ID"`
	ChatMsgID    int64 `json:"chat_msg_id" db:"chat_msg_id" desc:"聊天记录id"`
	SensitiveID  int64 `json:"sensitive_id" db:"sensitive_id" desc:"敏感词id"`
	CreatedAt    int64 `json:"created_at" db:"created_at" desc:"创建时间"`
}

// DetailListPageResponse 敏感词明细列表分页返回参数
type DetailListPageResponse struct {
	List  []DetailListResponse `json:"list"`
	Total int64                `json:"total"`
	Page  int                  `json:"page"`
	Size  int                  `json:"size"`
}

// DetailListResponse 敏感词明细列表返回参数
type DetailListResponse struct {
	ID          int64 `json:"id" db:"id" desc:"主键id"`
	ChatMsgID   int64 `json:"chat_msg_id" db:"chat_msg_id" desc:"聊天记录id"`
	SensitiveID int64 `json:"sensitive_id" db:"sensitive_id" desc:"敏感词id"`
	CreatedAt   int64 `json:"created_at" db:"created_at" desc:"创建时间"`
	ChatMsgMembers
}

// ChatMsgMembers 聊天记录对应员工和客户
type ChatMsgMembers struct {
	EnterpriseID       int64  `json:"enterprise_id" db:"enterprise_id" desc:"企业id"`
	MemberID           string `json:"member_id" db:"member_id" desc:"员工id"`
	MemberUserID       string `json:"member_user_id" db:"member_user_id" desc:"员工userid"`
	MemberName         string `json:"member_name" db:"member_name" desc:"员工姓名"`
	MemberAvatar       string `json:"member_avatar" db:"member_avatar" desc:"员工头像"`
	EXternalUserID     string `json:"external_user_id" db:"external_user_id" desc:"客户userid"`
	EXternalUserName   string `json:"external_user_name" db:"external_user_name" desc:"客户姓名"`
	EXternalUserAvatar string `json:"external_user_avatar" db:"external_user_avatar" desc:"客户头像"`
	EXternalUserRemark string `json:"external_user_remark" db:"external_user_remark" desc:"客户备注"`
}

// ChatMsg 聊天消息表结构体
type ChatMsg struct {
	ID           int64  `json:"id" db:"id" desc:"ID"`
	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id" desc:"企业id"`
	FromUser     string `json:"from_user" db:"from_user" desc:"发送用户"`
	ToUser       string `json:"to_user" db:"to_user" desc:"消息接收用户"`
	UserChannel  int64  `json:"user_channel" db:"user_channel" desc:"消息接收用户类型"`
}

// SendMessageRequest 发送请求
type SendMessageRequest struct {
	EnterpriseID  int64  `json:"enterprise_id" db:"enterprise_id" desc:"企业id"`
	ChatMessageID int64  `json:"chat_message_id" desc:"聊天记录id"`
	Type          string `json:"type" desc:"聊天记录类型"`
	ChatRecord    string `json:"chat_record" desc:"聊天记录"`
	FromUser      string `json:"from_user" desc:"发送方userid"`
	ToUser        string `json:"to_user" desc:"接收方userid"`
}

// UserInfo 聊天用户信息
type UserInfo struct {
	ObjectID     int64  `json:"object_id" db:"object_id" desc:"员工ID,外部联系人ID"`
	UserID       string `json:"user_id" db:"userid" desc:"企微userid"`
	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	Name         string `json:"name" db:"name" desc:"姓名"`
	Avatar       string `json:"avatar" db:"avatar" desc:"头像"`
	Remark       string `json:"remark" desc:"备注"`
	Status       int32  `json:"status" db:"status" desc:"备注"`
}
