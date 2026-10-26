package micro_scrm_enterprise

import (
	"fmt"

	"micro-scrm/internal/orm/model"
)

const (
	FINISHED    = 1
	ISNOTDELETE = 0
)

type TaskAllotLogModel struct {
	model.BranchModel
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
	ChatNum         int    `gorm:"chat_num;default:0;NOT NULL" json:"chat_num"`                          // 满足条件的社群数
	CloseTips       int    `gorm:"chat_num;default:0;NOT NULL" json:"close_tips"`                        // 已完成子任务未触达提示关闭
	ContactedCount  int    `gorm:"column:contacted_count;default:0;NOT NULL" json:"contacted_count"`     // 已联系的顾客人数
	FinishRate      int    `gorm:"column:finish_rate;default:0;NOT NULL" json:"finish_rate"`             // 导购的任务完成率，默认为0%
	MomentLookCount int    `gorm:"column:moment_look_count;default:0;NOT NULL" json:"moment_look_count"` // 导购的朋友圈联系人数
}
type PlanTasksConditionsAllotLog struct {
	UserId                   string `json:"user_id"`
	Name                     string `json:"name"`
	DepartmentId             int64  `json:"department_id"`
	DepartmentName           string `json:"department_name"`
	StartTime                int64  `json:"start_time"`
	EndTime                  int64  `json:"end_time"`
	IsFinished               int    `json:"is_finished"`
	CustomerNum              int    `json:"customer_num"`
	ChatID                   int64  `json:"chat_id"`
	ChatName                 string `json:"chat_name"`
	Count                    rune   `json:"count"`
	TagId                    int64  `json:"tag_id"`
	ID                       int64  `json:"id" db:"id"`
	WhetherToResign          int64  `json:"whether_to_resign"` // 添加是否离职状态   0为在职状态
	TotalReachCustomerNumber int64  `json:"total_reach_customer_number" db:"total_reach_customer_number" comment:"触达人数，历史代码接口定义的字段"`
	TotalOrderCustomerNumber int64  `json:"total_order_customer_number" db:"total_order_customer_number"`
	TotalOrderNumber         int64  `json:"total_order_number" db:"total_order_number"`
	TotalOrderFee            int64  `json:"total_order_fee" db:"total_order_fee"`
	TaskCategory             int    `json:"task_category"`
	GroupSendType            int    `json:"group_send_type"`
	FinishRate               string `json:"finish_rate"`
	ContactedCount           int64  `json:"contacted_count" comment:"沟通任务的触达人数"`
	MomentLookCount          int64  `json:"moment_look_count" comment:"朋友圈任务的可见人数"`
}

type DataBase struct {
	UserId        string `json:"user_id"`
	UserName      string `json:"user_name"`
	PlanName      string `json:"plan_name"`
	TaskName      string `json:"task_name"`
	IsFinished    int    `json:"is_finished"`
	StartTime     int64  `json:"start_time"`
	CompletedTime int64  `json:"completed_time"`
}

type ExportTaskDataResp struct {
	DataBase
	WorkNo                   string `json:"work_no"`
	DepartmentName           string `json:"department_name"`
	BelongingArea            string `json:"belonging_area"`
	CustomerNum              int    `json:"customer_num"`
	ID                       int64  `json:"id" db:"id"`
	WhetherToResign          int64  `json:"whether_to_resign"` // 添加是否离职状态   0为在职状态
	TotalReachCustomerNumber int64  `json:"total_reach_customer_number" db:"total_reach_customer_number"`
	TotalOrderCustomerNumber int64  `json:"total_order_customer_number" db:"total_order_customer_number"`
	TotalOrderNumber         int64  `json:"total_order_number" db:"total_order_number"`
	TotalOrderFee            int64  `json:"total_order_fee" db:"total_order_fee"`
	Category                 int    `json:"category"`
	ContactedCount           int    `json:"contacted_count"`
	MomentLookCount          int    `json:"moment_look_count"`
	FinishRate               string `json:"finish_rate"`
}

type ExportDetailDataResp struct {
	UserID          string `json:"user_id"`
	UserName        string `json:"user_name"`
	CompletedTime   int64  `json:"completed_time"`
	IsFinished      int    `json:"is_finished"`
	PlanName        string `json:"plan_name"`
	TaskName        string `json:"task_name"`
	BaName          string `json:"ba_name"`
	TagName         string `json:"tag_name"`
	UnionID         string `json:"union_id"`
	VipCode         string `json:"vip_code"`
	WhetherToResign int64  `json:"whether_to_resign"`
}

type ExportGroupChatDataResp struct {
	DataBase
	WorkNo                   string `json:"work_no"`
	BelongingArea            string `json:"belonging_area"`
	GroupChatName            string `json:"group_chat_name"`
	DepartmentName           string `json:"department_name"`
	ID                       int64  `json:"id" db:"id"`
	TotalReachCustomerNumber int64  `json:"total_reach_customer_number" db:"total_reach_customer_number"`
	TotalOrderCustomerNumber int64  `json:"total_order_customer_number" db:"total_order_customer_number"`
	TotalOrderNumber         int64  `json:"total_order_number" db:"total_order_number"`
	TotalOrderFee            int64  `json:"total_order_fee" db:"total_order_fee"`
}

type TaskAllotcustomerNum struct {
	UserId      string `json:"user_id"`
	CustomerNum int    `json:"customer_num"`
}

func (TaskAllotLogModel) TableName() string {
	return "task_allot_log"
}

func (t TaskAllotLogModel) AllTableName(enterpriseID int64) string {
	return fmt.Sprintf(" %s%d.%s ", t.DBPrefix(), enterpriseID, t.TableName())
}
