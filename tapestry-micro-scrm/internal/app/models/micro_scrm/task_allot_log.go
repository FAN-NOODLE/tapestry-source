package microScrmModel

import "micro-scrm/internal/app/models"

type TaskAllotLogModel struct {
	models.BaseModel
	ID              int64  `json:"id" gorm:"primaryKey;column:id;type:int(11);autoIncrement"`
	PlanId          int64  `gorm:"column:plan_id;default:0;NOT NULL" json:"plan_id"`           // 计划ID
	ConditionId     int64  `gorm:"column:condition_id;default:0;NOT NULL" json:"condition_id"` // 条件对应的子条件的ID
	TaskId          int64  `gorm:"column:task_id;default:0;NOT NULL" json:"task_id"`           // 条件对应的子任务的ID
	UserId          string `gorm:"column:user_id;NOT NULL" json:"user_id"`                     // 员工的企微user_id
	StartTime       int    `gorm:"column:start_time;default:0;NOT NULL" json:"start_time"`     // 任务开始时间
	EndTime         int    `gorm:"column:end_time;default:0;NOT NULL" json:"end_time"`         // 任务完成时间
	IsFinished      int    `gorm:"column:is_finished;default:0;NOT NULL" json:"is_finished"`   // 完成状态；0否、1是
	CustomerNum     int    `gorm:"column:customer_num;default:0;NOT NULL" json:"customer_num"` // 符合条件的客户数
	GroupChatId     int64  `gorm:"column:group_chat_id" json:"chat_id"`
	TagId           int64  `gorm:"column:tag_id;default:0" json:"tag_id"`                                // 标签编号
	TagName         string `gorm:"column:tag_name" json:"tag_name"`                                      // 标签名称                       // 群id
	IsForward       int    `gorm:"column:is_forward;default:0;NOT NULL" json:"is_forward"`               // 导购是否点击转发给顾客按钮，0否，1是
	GroupMsgId      string `gorm:"column:group_msg_id;default:'';NOT NULL" json:"group_msg_id"`          // 企业群发的msg_id
	IsGroupSend     int    `gorm:"column:is_group_send;default:0;NOT NULL" json:"is_group_send"`         // 企业群发是否发送，1为已发送
	AgentMsgId      string `gorm:"column:agent_msg_id;default:'';NOT NULL" json:"agent_msg_id"`          // 自建应用的消息msg_id
	AgentMsgTime    int    `gorm:"column:agent_msg_time;default:0;NOT NULL" json:"agent_msg_time"`       // 自建应用消息发送时间
	ChatNum         int    `gorm:"column:chat_num;default:0;NOT NULL" json:"chat_num"`                   // 满足条件的社群数
	CloseTips       int    `gorm:"column:close_tips;default:0;NOT NULL" json:"close_tips"`               // 是否关闭未触达的提示
	ContactedCount  int    `gorm:"column:contacted_count;default:0;NOT NULL" json:"contacted_count"`     // 已联系的顾客人数
	FinishRate      int    `gorm:"column:finish_rate;default:0;NOT NULL" json:"finish_rate"`             // 导购的任务完成率，默认为0%
	MomentLookCount int    `gorm:"column:moment_look_count;default:0;NOT NULL" json:"moment_look_count"` // 导购的朋友圈联系人数
}

func (TaskAllotLogModel) TableName() string {
	return "task_allot_log"
}

const (
	PLAN_FISISHING = iota // 待完成
	PLAN_FINISHED         // 已完成
	PLAN_NOFINISH         // 已结束，未完成
)

const (
	QyGroupMsgSend       = iota // 未发送
	QyGroupMsgSendFinish        // 已发送
	QyGroupMsgSendCancel        // 已取消
)

const (
	IS_UNION_TRUE = iota + 1 // 筛选是否是好友，1为好友，2为非好友
	IS_UNION_FALSE
)

type TaskIdGroupNum struct {
	PlanId        int64 `json:"plan_id"`
	TaskId        int64 `json:"task_id"`
	GroupMsgCount int64 `json:"group_msg_count"`
	CreatedAt     int64 `json:"created_at"`
}

type TaskFinish struct {
	TaskId int64 `json:"task_id"`
	Num    int   `json:"num"`
}
