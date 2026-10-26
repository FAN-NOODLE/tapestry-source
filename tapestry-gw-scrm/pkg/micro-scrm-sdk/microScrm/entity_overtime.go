package microScrm

// 接收消息方为员工
const ReceiveTypeMember = 1

//  chat_overtime_records 表实体结构
type Records struct {
	ChatMsgID    int64 `json:"chat_msg_id" db:"chat_msg_id" desc:"chat_msg:消息id"`
	EnterpriseID int64 `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	MsgAt        int64 `json:"msg_at"  desc:"触发时间"`
}

// ListRequest 获取聊天超时列表请求参数
type ListRequest struct {
	CustomerName string `json:"customer_name" desc:"客户姓名"`
	KeyWords     string `json:"keywords" desc:"超时内容"`
	SearchType   int32  `json:"search_type" desc:"搜索类型"`
	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id" validate:"required" required:"企业ID不能为空" desc:"企业ID"`
	MemberID     int    `json:"member_id"`
	DepartmentID int    `json:"department_id"`
	TimeStart    int    `json:"time_start"`
	TimeEnd      int    `json:"time_end"`
	Page         int    `json:"page"`
	Size         int    `json:"size"`
}

// ResponseList 返回聊天超时记录列表
type PageListResponse struct {
	List  []RecordsListResponse `json:"list"`
	Total int64                 `json:"total"`
	Page  int                   `json:"page"`
	Size  int                   `json:"size"`
}

// RecordsListResponse 聊天超时记录列表返回参数
type RecordsListResponse struct {
	Records
	ChatMsgMembers
	ChatMsg
}

// ChatMsg 聊天消息表结构体
type ChatMsg struct {
	ID           int64                  `json:"id" db:"id" desc:"ID"`
	EnterpriseID int64                  `json:"enterprise_id" db:"enterprise_id" desc:"企业id"`
	FromUser     string                 `json:"from_user" db:"from_user" desc:"发送用户"`
	ToUser       string                 `json:"to_user" db:"to_user" desc:"消息接收用户"`
	UserChannel  int32                  `json:"user_channel" db:"user_channel" desc:"接收用户类型"`
	Content      string                 `json:"content" db:"content" desc:"消息内容"`
	ContentMap   map[string]interface{} `json:"content_map" desc:"内容decode"`
	ContentType  string                 `json:"type" db:"type" desc:"消息类型"`
	TimeStart    int64                  `json:"time_start" db:"msg_at" desc:"发送时间"`
	TimeEnd      int64                  `json:"time_end" desc:"发送截止时间"`
}

// ChatMsgMembers 聊天记录对应员工和客户
type ChatMsgMembers struct {
	EnterpriseID       int64             `json:"enterprise_id" db:"enterprise_id" desc:"企业id"`
	MemberUserID       string            `json:"member_user_id" db:"member_user_id" desc:"员工userid"`
	MemberID           int64             `json:"member_id" db:"member_id" desc:"员工id"`
	MemberName         string            `json:"member_name" db:"member_name" desc:"员工姓名"`
	MemberAvatar       string            `json:"member_avatar" db:"member_avatar" desc:"员工头像"`
	EXternalUserID     string            `json:"external_user_id" db:"external_user_id" desc:"客户userid"`
	EXternalUserName   string            `json:"external_user_name" db:"external_user_name" desc:"客户姓名"`
	EXternalUserAvatar string            `json:"external_user_avatar" db:"external_user_avatar" desc:"客户头像"`
	EXternalUserRemark string            `json:"external_user_remark" db:"external_user_remark" desc:"客户备注"`
	MemberDepartsList  MemberDepartsList `json:"member_departs" usage:"部门信息"`
}

// 添加超时记录结构体
type AddOverRecord struct {
	ID           int64 `json:"id" db:"id" desc:"ID"`
	ChatMsgID    int64 `json:"chat_msg_id" db:"chat_msg_id" desc:"chat_msg:消息id"`
	EnterpriseID int64 `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	CreatedAt    int64 `json:"created_at" db:"created_at" desc:"创建时间"`
}
