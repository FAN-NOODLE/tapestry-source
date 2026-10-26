package ws_trace

const (
	SourceTypeOneToOne      = 1 + iota //1v1聊天
	SourceTypeGroupChat                //群聊
	SourceTypeTask                     //任务计划
	SourceTypeGroupChatSend            //社群群发
)

type WsTrace struct {
	WsTraceID     string `json:"ws_trace_id" db:"ws_trace_id" desc:"ws_trace表的id"`
	UkVoucher     string `json:"uk_voucher" db:"uk_voucher" desc:"不同类型的聚合唯一索引"`
	MemberID      int64  `json:"member_id" db:"member_id" desc:"员工id"`
	WxExternalID  int64  `json:"wx_external_id" db:"wx_external_id" desc:"外部联系人id"`
	Unionid       string `json:"unionid" db:"unionid" desc:"员工unionid"`
	GroupChatID   int64  `json:"group_chat_id" db:"group_chat_id" desc:"ws系统外部群表id"`
	MainTaskID    int64  `json:"main_task_id" db:"main_task_id" desc:"主任务id"`
	ChildTaskID   int64  `json:"child_task_id" db:"child_task_id" desc:"子任务id"`
	SubTaskID     int64  `json:"sub_task_id" db:"sub_task_id" desc:"子任务下面的每个员工的task_allot_log表id"`
	QwGroupSendID int64  `json:"qw_group_send_id" db:"qw_group_send_id" desc:"企微群发任务id"`
	SourceType    int64  `json:"source_type" db:"source_type" desc:"来源：0-未知 1-1v1聊天 2-群聊 3-任务计划 4-社群群发"`
	CreatedAt     int64  `json:"created_at" db:"created_at" desc:"创建时间"`
	UpdatedAt     int64  `json:"updated_at" db:"updated_at" desc:"更新时间"`
}
