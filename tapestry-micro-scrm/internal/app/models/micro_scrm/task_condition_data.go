package microScrmModel

import "micro-scrm/internal/app/models"

type TaskConditionDataModel struct {
	models.BaseModel
	TagName            string `gorm:"column:tag_name;NOT NULL" json:"tag_name"`                             // CampaignTool标签名称
	WwTagId            string `gorm:"column:ww_tag_id;NOT NULL" json:"ww_tag_id"`                           // 企微标签id
	UnionId            string `gorm:"column:union_id;NOT NULL" json:"union_id"`                             // 顾客的UnionID
	WorkExternalUserId string `gorm:"column:work_external_user_id;NOT NULL" json:"work_external_user_id"`   // 外部联系人user_id,可能空
	IsUnion            int    `gorm:"column:is_union;default:0;NOT NULL" json:"is_union"`                   // 关系状态；0没有好友关系、1有好友关系
	Mobile             string `gorm:"column:mobile" json:"mobile"`                                          // 手机
	MarsId             string `gorm:"column:mars_id;NOT NULL" json:"mars_id"`                               // BA的MarsCode
	UserId             string `gorm:"column:user_id;NOT NULL" json:"user_id"`                               // 员工的企微user_id
	ConditionId        int64  `gorm:"column:condition_id;default:0;NOT NULL" json:"condition_id"`           // 条件对应的子任务的ID
	TaskId             int64  `gorm:"column:task_id;default:0;NOT NULL" json:"task_id"`                     // 条件对应的子任务的ID
	PlanId             int64  `gorm:"column:plan_id;default:0;NOT NULL" json:"plan_id"`                     // 条件对应的子任务的ID
	TaskAllotLogId     int64  `gorm:"column:task_allot_log_id;default:0;NOT NULL" json:"task_allot_log_id"` // 记录表的ID
	CustomerID         string `json:"customer_id" gorm:"column:customer_id;type:varchar(255);not null;comment:用户的customerid信息"`
	IsFinished         int32  `json:"is_finished" gorm:"column:is_finished;type:tinyint(2);not null;default:0;comment:判断任务是否标记完成"`
	FinishedAt         int    `json:"finished_at" gorm:"column:finished_at;type:int(10);not null;default:0;comment:标记完成的时间"`
	Remark             string `json:"remark" gorm:"column:remark;type:varchar(255);not null;comment:联系说明"`
	IsContact          int    `json:"is_contact" gorm:"column:is_contact;type:tinyint(2);not null;default:0;comment:是否联系好友"`
	IsCall             int    `json:"is_call" gorm:"column:is_call;type:tinyint(2);not null;default:0;comment:是否拨打电话"`
	IsAdd              int    `json:"is_add" gorm:"column:is_add;type:tinyint(2);not null;default:0;comment:是否点击添加好友"`
	IsMatch            int    `json:"is_match" gorm:"column:is_match;type:tinyint(2);not null;default:1;comment:是否匹配成功"`
	UserName           string `json:"user_name"  gorm:"column:user_name;type:varchar(255);not null;comment:用户的customer昵称信息"`
	IsTouch            int    `json:"is_touch"  gorm:"column:is_touch;type:tinyint(2);not null;default:0;comment:是否触达,0为否，1为是"`
	GroupChatId        int64  `json:"group_chat_id"  gorm:"column:group_chat_id;type:bigint(20);not null;default:0;comment:社群id"`
	CustomerStatus     int    `json:"customer_status" gorm:"column:customer_status;type:tinyint(2);not null;default:0;comment:顾客沟通类型，0为群发沟通，1为单独沟通"`
	Forward            int    `json:"forward" gorm:"column:forward;type:tinyint(2);not null;default:0;comment:是否转发给顾客，0为否，1为是"`
}

const (
	TouchStatusNo    = 0 // 未触达
	TouchStatusYes   = 1 // 已触达
	TouchStatusFail  = 2 // 好友已删除
	TouchStatusLimit = 3 // 群发消息已达上线
)

const (
	FinishedNo  = iota // 未完成
	FinishedYes        // 已完成
)

const (
	TaskReach    = iota // 顾客列表的筛选，任务触达
	TaskExist           // 顾客列表的筛选，有任务
	InteractTime        // 顾客列表的筛选，SW活跃沟通时间
)

const (
	SingleForward = iota
	AllForward
)

func (TaskConditionDataModel) TableName() string {
	return "task_condition_data"
}
