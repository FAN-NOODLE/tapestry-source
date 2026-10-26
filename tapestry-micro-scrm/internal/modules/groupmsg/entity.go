package groupmsg

import "micro-scrm/internal/pkg/common"

const (
	ChatTypeGroup          = "group"  // 企微群发任务的类型，group表示发送给客户群
	ChatTypeUser           = "single" // single 发送给用户
	MsgTaskProcessLockTime = 60       // 群发任务处理锁时间
)

// 任务触达状态
const (
	Sending        = 6 // 发送中 未下发给执行员工
	AtOnceNotStart = 1 // 定时未开始
	NotSend        = 2 // 未发送
	PartNotSend    = 3 // 部分未发送
	AlreadySend    = 4 // 已发送
	SendFail       = 5 // 发送失败 发送给执行员工失败
)

// 群发类型
const (
	AtOnce = 1 // 立即
	Timing = 2 // 定时
)

// 任务成员-调用企微接口下发任务状态
const (
	MsgTaskMemberTaskStatusUndo   = 1 // 未发送
	MsgTaskMemberTaskStatusDone   = 2 // 已发送
	MsgTaskMemberTaskStatusFailed = 3 // 发送失败
)

// 任务成员-群发发送消息状态
const (
	MsgTaskMemberStatusNotSend = 1 // 未发送
	MsgTaskMemberStatusSend    = 2 // 已发送
)

// 群发消息调用企业微信接口下发任务状态
const (
	MsgTaskTaskStatusUndo       = 1 // 定时任务未发送
	MsgTaskTaskStatusDoing      = 2 // 任务发送中
	MsgTaskTaskStatusDone       = 3 // 任务已发送
	MsgTaskTaskStatusFailed     = 4 // 任务发送失败
	MsgTaskTaskStatusFailedDone = 5 // 任务发送失败结束
)

// 发送对象
const (
	MsgTaskTaskSendTypeUser  = 1 // 人员
	MsgTaskTaskSendTypeGroup = 2 // 群组
)

// ChatGroupMsgTaskMember 群发消息员工表结构体
type ChatGroupMsgTaskMember struct {
	ID           int64  `db:"id" json:"id"`
	EnterpriseId int64  `db:"enterprise_id" json:"enterprise_id"`
	MsgTaskId    int64  `db:"msg_task_id" json:"msg_task_id"`
	MemberId     int64  `db:"member_id" json:"member_id"`
	Status       int    `db:"status" json:"status"`
	CreatedAt    int64  `db:"created_at" json:"created_at"`
	UpdatedAt    int64  `db:"updated_at" json:"updated_at"`
	DeletedAt    int64  `db:"deleted_at" json:"deleted_at"`
	TaskStatus   int    `db:"task_status" json:"task_status"`
	FailMsg      string `db:"fail_msg" json:"fail_msg"`
	SuccessCount int    `db:"success_count" json:"success_count"`
	NotSendCount int    `db:"not_send_count" json:"not_send_count"`
	FailedCount  int    `db:"failed_count" json:"failed_count"`
	MsgID        string `db:"msgid" json:"msgid"`
	SendTime     int64  `db:"send_time" json:"send_time"`
}

// ChatGroupMsgTask 群发消息任务表结构体
type ChatGroupMsgTask struct {
	ID              int64  `db:"id" json:"id"`
	EnterpriseId    int64  `db:"enterprise_id" json:"enterprise_id"`
	Name            string `db:"name" json:"name"`
	Type            int    `db:"type" json:"type"`
	SendType        int    `db:"send_type" json:"send_type"`
	Status          int    `db:"status" json:"status"`
	TaskStatus      int    `db:"task_status" json:"task_status"`
	FailMsg         string `db:"fail_msg" json:"fail_msg"`
	SendTime        int64  `db:"send_time" json:"send_time"`
	LastSendMsg     int    `db:"last_send_msg" json:"last_send_msg"`
	LatestCheckTime int64  `db:"latest_check_time" json:"latest_check_time"`
	Content         string `db:"content" json:"content"`
	CreatedAt       int64  `db:"created_at" json:"created_at"`
	UpdatedAt       int64  `db:"updated_at" json:"updated_at"`
	DeletedAt       int64  `db:"deleted_at" json:"deleted_at"`
}

// ChatGroupMsgSendRecord 群发结果记录表结构体
type ChatGroupMsgSendRecord struct {
	ID              int64  `db:"id" json:"id"`
	EnterpriseId    int64  `db:"enterprise_id" json:"enterprise_id"`
	MsgTaskId       int64  `db:"msg_task_id" json:"msg_task_id"`
	MsgTaskMemberId int64  `db:"msg_task_member_id" json:"msg_task_member_id"`
	ChatId          string `db:"chat_id" json:"chat_id"`
	ExternalUserID  string `db:"external_userid" json:"external_user_id"`
	Status          int    `db:"status" json:"status"`
	SendTime        int64  `db:"send_time" json:"send_time"`
	CreatedAt       int64  `db:"created_at" json:"created_at"`
	UpdatedAt       int64  `db:"updated_at" json:"updated_at"`
	DeletedAt       int64  `db:"deleted_at" json:"deleted_at"`
}

// 任务队列内容结构体
type TaskQueueData struct {
	EnterpriseID int64 `json:"enterprise_id"`
	TaskID       int64 `json:"task_id"`
}

// GroupMsgAddRequest 新增群发消息任务请求结构体
type GroupMsgAddRequest struct {
	ID           int64                  `db:"id" json:"id"`
	EnterpriseId int64                  `db:"enterprise_id" json:"enterprise_id" validate:"required" required:"企业id不能为空"`
	Name         string                 `db:"name" json:"name" validate:"required" required:"任务名不能为空"`
	Type         int                    `db:"type" json:"type" validate:"required" required:"任务类型不能为空"`
	SendType     int                    `db:"send_type" json:"send_type" validate:"required" required:"发送类型不能为空"`
	Status       int                    `db:"status" json:"status"`
	TaskStatus   int                    `db:"task_status" json:"task_status"`
	SendTime     int64                  `db:"send_time" json:"send_time"`
	LastSendMsg  int                    `db:"last_send_msg" json:"last_send_msg"`
	CreatedAt    int64                  `db:"created_at" json:"created_at"`
	UpdatedAt    int64                  `db:"updated_at" json:"updated_at"`
	DeletedAt    int64                  `db:"deleted_at" json:"deleted_at"`
	Content      map[string]interface{} `json:"content" validate:"required" required:"内容不能为空"`
	MemberIds    []int64                `json:"member_ids" validate:"required" required:"成员不能为空"`
}

// GroupMsgAddRequest 新增群发消息任务请求结构体
type GroupMsgEditRequest struct {
	ID           int64                  `db:"id" json:"id" validate:"required" required:"任务id不能为空"`
	EnterpriseId int64                  `db:"enterprise_id" json:"enterprise_id" validate:"required" required:"企业id不能为空"`
	Name         string                 `db:"name" json:"name" validate:"required" required:"任务名不能为空"`
	Type         int                    `db:"type" json:"type" validate:"required" required:"任务类型不能为空"`
	SendType     int                    `db:"send_type" json:"send_type" validate:"required" required:"发送类型不能为空"`
	Status       int                    `db:"status" json:"status"`
	TaskStatus   int                    `db:"task_status" json:"task_status"`
	SendTime     int64                  `db:"send_time" json:"send_time"`
	LastSendMsg  int                    `db:"last_send_msg" json:"last_send_msg"`
	CreatedAt    int64                  `db:"created_at" json:"created_at"`
	UpdatedAt    int64                  `db:"updated_at" json:"updated_at"`
	DeletedAt    int64                  `db:"deleted_at" json:"deleted_at"`
	Content      map[string]interface{} `json:"content" validate:"required" required:"内容不能为空"`
	MemberIds    []int64                `json:"member_ids" validate:"required" required:"成员不能为空"`
	ContentStr   string                 `db:"content_str" json:"content_str"`
}

// MemberListRequest任务明细-执行群主下拉列表 请求结构体
type MemberListRequest struct {
	EnterpriseID int64 `db:"enterprise_id" json:"enterprise_id"`
	ID           int64 `json:"id" db:"id"`
}

// Member成员
type Member struct {
	ID        int64  `json:"id" db:"id"`
	UserID    string `json:"user_id" db:"user_id"`
	Name      string `json:"name" db:"name"`
	MemberId  int64  `json:"member_id" db:"member_id"`
	SendTime  int    `db:"send_time" json:"send_time"`
	CreatedAt int64  `db:"created_at" json:"created_at"`
	UpdatedAt int64  `db:"updated_at" json:"updated_at"`
	DeletedAt int64  `db:"deleted_at" json:"deleted_at"`
}

// TaskIfoRequest任务详情 请求结构体
type TaskIfoRequest struct {
	EnterpriseID int64 `db:"enterprise_id" json:"enterprise_id"`
	ID           int64 `json:"id" db:"id"`
}

// 群发消息列表请求结构体
type MsgTaskRequest struct {
	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id" validate:"required" required:"企业Id必传"`
	Name         string `json:"name" db:"name"`
	Status       int64  `json:"status" db:"status"`
	SendType     int    `json:"send_type" db:"send_type"`
	MemberID     int64  `json:"member_id" db:"member_id"`
	StartAt      int64  `json:"start_at" db:"start_at"`
	EndAt        int64  `json:"end_at" db:"end_at"`
	Page         int    `json:"page" db:"page"`
	Size         int    `json:"size" db:"size"`
}

// 群发消息列表请求结构体回参数
type MsgTaskListResponse struct {
	List  []ChatGroupMsgTaskList `json:"list"`
	Total int64                  `json:"total"`
	Page  int                    `json:"page" db:"page" desc:"页码"`
	Size  int                    `json:"size"  db:"size" desc:"分页"`
}

// ChatGroupMsgTask 群发消息任务表结构体
type ChatGroupMsgTaskList struct {
	ID           int64                  `db:"id" json:"id"`
	EnterpriseId int64                  `db:"enterprise_id" json:"enterprise_id"`
	Name         string                 `db:"name" json:"name"`
	Members      string                 `db:"members" json:"members"`
	Type         int                    `db:"type" json:"type"`
	SendType     int                    `db:"send_type" json:"send_type"`
	Status       int                    `db:"status" json:"status"`
	TaskStatus   int                    `db:"task_status" json:"task_status"`
	FailMsg      string                 `db:"fail_msg" json:"fail_msg"`
	SendTime     int64                  `db:"send_time" json:"send_time"`
	LastSendMsg  int                    `db:"last_send_msg" json:"last_send_msg"`
	Content      string                 `db:"content" json:"content"`
	ContentObj   map[string]interface{} `db:"content_obj" json:"content_obj"`
	CreatedAt    int64                  `db:"created_at" json:"created_at"`
	UpdatedAt    int64                  `db:"updated_at" json:"updated_at"`
	DeletedAt    int64                  `db:"deleted_at" json:"deleted_at"`
}

// TaskInfoResponse任务详情 请求结构体
type TaskInfoResponse struct {
	ID           int64                  `db:"id" json:"id"`
	EnterpriseId int64                  `db:"enterprise_id" json:"enterprise_id"`
	Name         string                 `db:"name" json:"name"`
	Type         int                    `db:"type" json:"type"`
	SendType     int                    `db:"send_type" json:"send_type"`
	Status       int                    `db:"status" json:"status"`
	TaskStatus   int                    `db:"task_status" json:"task_status"`
	FailMsg      string                 `db:"fail_msg" json:"fail_msg"`
	SendTime     int64                  `db:"send_time" json:"send_time"`
	LastSendMsg  int                    `db:"last_send_msg" json:"last_send_msg"`
	Content      map[string]interface{} `db:"content" json:"content"`
	CreatedAt    int64                  `db:"created_at" json:"created_at"`
	UpdatedAt    int64                  `db:"updated_at" json:"updated_at"`
	DeletedAt    int64                  `db:"deleted_at" json:"deleted_at"`
	Members      []*Member              `json:"members"`
}

// TaskDetailRequest任务明细 请求结构体
type TaskDetailRequest struct {
	EnterpriseID int64  `db:"enterprise_id" json:"enterprise_id"`
	ID           int64  `json:"id" db:"id"`
	ChatName     string `db:"chat_name" json:"chat_name"`
	StartTime    int64  `db:"start_time" json:"start_time"`
	EndTime      int64  `db:"end_time" json:"end_time"`
	MemberId     int64  `db:"member_id" json:"member_id"`
	Page         int    `db:"page" json:"page"`
	Size         int    `db:"size" json:"size"`
}

// TaskDetailResponse任务明细 响应结构体
type TaskDetailResponse struct {
	Page  int               `db:"page" json:"page"`
	Size  int               `db:"size" json:"size"`
	List  []*TaskDetailItem `json:"list"`
	Total int64             `db:"total" json:"total"`
}

type TaskDetailItem struct {
	ChatName string `db:"chat_name" json:"chat_name"`
	Name     string `db:"name" json:"name"`
	SendTime int64  `db:"send_time" json:"send_time"`
}

type TaskDetailMemberRequest struct {
	EnterpriseID int64 `db:"enterprise_id" json:"enterprise_id"`
	ID           int64 `db:"id" json:"id"`
	MemberId     int64 `db:"member_id" json:"member_id"`
	StartTime    int64 `db:"start_time" json:"start_time"`
	TaskStatus   int64 `db:"task_status" json:"task_status"`
	Status       int64 `db:"status" json:"status"`
	EndTime      int64 `db:"end_time" json:"end_time"`
	Page         int   `db:"page" json:"page"`
	Size         int   `db:"size" json:"size"`
}

type TaskDetailMemberResponse struct {
	Page  int                     `db:"page" json:"page"`
	Size  int                     `db:"size" json:"size"`
	List  []*TaskDetailMemberItem `json:"list"`
	Total int64                   `db:"total" json:"total"`
}

type TaskDetailMemberItem struct {
	Name           string `db:"name" json:"name"`
	MemberSendTime int64  `db:"member_send_time" json:"member_send_time"`
	MemberId       int64  `db:"member_id" json:"member_id"`
	Msgid          string `db:"msgid" json:"msgid"`
	TaskStatus     int64  `db:"task_status" json:"task_status"`
	Status         int64  `db:"status" json:"status"`
	SuccessCount   int64  `db:"success_count" json:"success_count"`
	NotSendCount   int64  `db:"not_send_count" json:"not_send_count"`
	FailedCount    int64  `db:"failed_count" json:"failed_count"`
}

// MsgContent 群发消息content结构
type MsgContent struct {
	Text        *common.Text        `json:"text,omitempty" desc:"消息文本内容"`
	Image       *common.Image       `json:"image,omitempty" desc:"图片"`
	Link        *common.Link        `json:"link,omitempty" desc:"图片"`
	MiniProgram *common.MiniProgram `json:"miniprogram,omitempty" desc:"miniprogram"`
}

type TaskMemberSendListReq struct {
	EnterpriseID int64  `json:"enterprise_id" binding:"required" required:"企业Id必填"`
	MemberId     int64  `json:"member_id" binding:"required" required:"员工member_id必填"`
	MsgId        string `json:"msg_id"`
	Export       int    `json:"export"`
	Page         int    `db:"page" json:"page"`
	Size         int    `db:"size" json:"size"`
}

type TaskMemberSendListResp struct {
	Page  int                  `json:"page"`
	Size  int                  `json:"size"`
	List  []TaskMemberSendList `json:"list"`
	Total int64                `json:"total"`
}

type TaskMemberSendList struct {
	ChatId         string `json:"chat_id"`
	ChatName       string `json:"chat_name"`
	UserId         string `json:"user_id"`
	UserName       string `json:"user_name"`
	ExternalUserID string `json:"external_userid"`
	Status         int    `json:"status"`
	SendTime       int64  `json:"send_time"`
}
