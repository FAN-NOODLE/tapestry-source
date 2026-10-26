package plan_task

type ResponseFrontPlanList struct {
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	Remark    string `json:"remark"`
	Status    int    `json:"status"`
	StatusTex string `json:"status_tex"`
	IsOpen    int    `json:"is_open"`
	StartTime int64  `json:"start_time"`
	EndTime   int64  `json:"end_time"`
	CreatedAt int64  `json:"created_at"`
}

type ResponseFrontPlanTaskList struct {
	Id             int64  `json:"id"`
	Name           string `json:"name"`
	Remark         string `json:"remark"`
	ConditionType  int    `json:"condition_type"`
	Dimension      int    `json:"dimension"` // dimension
	Status         int    `json:"status"`
	CustomerNum    int    `json:"customer_num"`
	ConditionID    int64  `json:"condition_id"`
	GroupChatID    int64  `json:"group_chat_id"`
	StartTime      int64  `json:"start_time"`
	EndTime        int64  `json:"end_time"`
	TaskType       int    `json:"task_type"`
	TaskCategory   int    `json:"task_category"`
	GroupSendType  int    `json:"group_send_type"`
	CloseTips      int    `json:"close_tips"`
	ChatNum        int    `json:"chat_num"`
	PlanId         int64  `json:"plan_id"`
	PlanName       string `json:"plan_name"`
	TaskId         int64  `json:"task_id"`
	FinishRate     int    `json:"finish_rate"`
	TaskStartTime  int64  `json:"task_start_time"`
	TaskEndTime    int64  `json:"task_end_time"`
	TaskExpireType int    `json:"task_expire_type"`
}

type ResponseFrontPlanTaskConList struct {
	Id            int64  `json:"id"`
	TaskId        int64  `json:"task_id"`
	Conditions    string `json:"conditions"`
	Content       string `json:"content"`
	ContentType   int    `json:"content_type"`
	ConditionType int    `json:"condition_type"`
	ConditionID   int64  `json:"condition_id"`
	IsFinished    int    `json:"is_finished"`
	Dimension     int    `json:"dimension"`
	GroupChatID   int64  `json:"group_chat_id"`
	StartTime     int64  `json:"start_time"`
	EndTime       int64  `json:"end_time"`
	Remark        string `json:"remark"`
}

type ResponseFrontTagUserList struct {
	Unionid            string   `json:"unionid"`
	Avatar             string   `json:"avatar"`
	Name               string   `json:"name"`
	TagName            []string `json:"tag_name"`
	Remark             string   `json:"remark"`
	IsUnion            int      `json:"is_union"`
	Mobile             string   `json:"mobile"`
	WorkExternalUserId string   `json:"work_external_user_id"`
}

type ResponseFrontTagsList struct {
	UnionId string `json:"union_id"`
	Avatar  string `json:"avatar"`
	Name    string `json:"name"`
	TagName string `json:"tag_name"`
	Remark  string `json:"remark"`
	IsUnion int    `json:"is_union"`
}

type ResponseFinishDays struct {
	UserId      string `json:"user_id"`
	CreatedAt   string `json:"created_at"`
	LogCount    string `json:"log_count"`
	IsFinishNum string `json:"is_finish_num"`
}
