package micro_scrm_enterprise

type JobReceive struct {
	TagName            string `json:"tag_name"`
	TagId              string `json:"tag_id"`
	UnionId            string `json:"union_id"`
	WorkExternalUserId string `json:"work_external_user_id"`
	UserId             string `json:"user_id"`
	Mobile             string `json:"mobile"`
	RelationId         int64  `json:"relation_id"`
	TagTableId         int64  `json:"tag_table_id"`
}
type ResponCondition struct {
	PlanId         int64  `json:"plan_id"`
	PlanName       string `json:"plan_name"`
	TaskId         int64  `json:"task_id"`
	TaskName       string `json:"task_name"`
	PlanStatus     int    `json:"plan_status"`
	TaskStatus     int    `json:"task_status"`
	ConditionId    int64  `json:"condition_id"`
	IsSpecifyCrowd int    `json:"is_specify_crowd"` // 任务是否指定人群 1 否 2 是
	// IsPush        int    `json:"is_push"`
	Conditions     string `json:"conditions"`
	ConditionType  int    `json:"conditions"`
	StartTime      int64  `json:"start_time"`
	EndTime        int64  `json:"end_time"`
	ContentType    int    `json:"content_type"`
	Content        string `json:"content"`
	TaskType       int    `json:"task_type"`
	SendType       int    `json:"send_type"`
	TaskIsPush     int    `json:"task_is_push"`
	PlanSendTime   int64  `json:"plan_send_time"`
	OriginTaskId   int64  `json:"origin_task_id"`
	TaskStartTime  int64  `json:"task_start_time"`
	TaskEndTime    int64  `json:"task_end_time"`
	TaskExpireType int    `json:"task_expire_type"`
}
type CloserResult struct {
	ExternalTagId      int64  `json:"external_tag_id"`
	RelationId         int64  `json:"relation_id"`
	UserId             string `json:"user_id"`
	WorkExternalUserId string `json:"work_external_user_id"`
	TagId              string `json:"tag_id"`
}

type GetGroupChatName struct {
	ID     int64  `json:"id"`
	ChatID string `json:"chat_id"`
	Name   string `json:"name"`
}
type GetLeaveInfo struct {
	GroupChatID int64 `json:"group_chat_id"`
	LeaveTime   int64 `json:"leave_time"`
}

type CountInfo struct {
	UnionID int64 `json:"unionid"`
}

type DetailInfo struct {
	Name               string `json:"name"`
	Gender             int    `json:"gender"`
	WorkExternalUserId string `json:"work_external_user_id"`
	Avatar             string `json:"avatar"`
}
type DetailInfo2 struct {
	Name               string `json:"name"`
	Gender             int    `json:"gender"`
	WorkExternalUserId string `json:"work_external_user_id"`
	Avatar             string `json:"avatar"`
	GroupChatID        int64  `json:"group_chat_id"`
}

type NewResponseFrontPlanTaskList struct {
	Id             int64           `json:"id"`
	Name           string          `json:"name"`
	Remark         string          `json:"remark"`
	ConditionType  int             `json:"condition_type"`
	Dimension      int             `json:"dimension"` // dimension
	Status         int             `json:"status"`
	CustomerNum    int             `json:"customer_num"`
	ConditionInfo  ConditionDetail `json:"condition_info"`
	ContentInfo    Content         `json:"content_info"`
	ConditionID    int64           `json:"condition_id"`
	GroupNumber    int64           `json:"group_number"`
	LeaveNumber    int64           `json:"leave_number"`
	GroupChatID    int64           `json:"group_chat_id"`
	StatusCode     int             `json:"status_code"`
	TaskType       int             `json:"task_type"`
	TaskCategory   int             `json:"task_category"`
	NoFinishedNum  int64           `json:"no_finished_num"`
	CloseTips      int             `json:"close_tips"`
	PlanId         int64           `json:"plan_id"`
	PlanName       string          `json:"plan_name"`
	TaskId         int64           `json:"task_id"`
	GroupSendType  int             `json:"group_send_type"`
	FinishRate     int             `json:"finish_rate"`
	TaskStartTime  int64           `json:"task_start_time"`
	TaskEndTime    int64           `json:"task_end_time"`
	TaskExpireType int             `json:"task_expire_type"`
}

type GroupChatInfo struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	ChatID string `json:"chat_id"`
}

type ConditionInfo struct {
	ID         int64  `json:"id"`
	Conditions string `json:"conditions"`
	Content    string `json:"content"`
}

type NewResponseFrontPlanTaskLists struct {
	Id            int64  `json:"id"`
	Name          string `json:"name"`
	Remark        string `json:"remark"`
	ConditionType int    `json:"condition_type"`
	Dimension     int    `json:"dimension"` // dimension
	Status        int    `json:"status"`
	CustomerNum   int    `json:"customer_num"`
	ConditionID   int64  `json:"condition_id"`
	GroupChatID   int64  `json:"group_chat_id"`
}

type Receives struct {
	ID             int64  `json:"id"`
	PlanID         int64  `json:"plan_id"`
	CreatedAt      int64  `json:"created_at"`
	UserId         string `json:"user_id"`
	TaskID         int64  `json:"task_id"`
	ConditionID    int64  `json:"condition_id"`
	GroupChatId    int64  `json:"group_chat_id"`
	IsFinished     int    `json:"is_finished"`
	StartTime      int64  `json:"start_time"`
	EndTime        int64  `json:"end_time"`
	Conditions     string `json:"conditions"`
	Content        string `json:"content"`
	Remark         string `json:"remark"`
	Dimension      int    `json:"dimension"`
	ConditionType  int    `json:"condition_type"`
	ContentType    int    `json:"content_type"`
	FinishRate     int    `json:"finish_rate"`
	TaskStartTime  int64  `json:"task_start_time"`
	TaskEndTime    int64  `json:"task_end_time"`
	TaskExpireType int    `json:"task_expire_type"`
}
type Receives2 struct {
	ID          int64  `json:"id"`
	UserId      string `json:"user_id"`
	ConditionID int64  `json:"condition_id"`
	StartTime   int64  `json:"start_time"`
	EndTime     int64  `json:"end_time"`
	IsFinished  int    `json:"is_finished"`
	Status      int    `json:"status"`
	IsForward   int    `json:"is_forward"`
	TaskId      int64  `json:"task_id"`
	PlanName    string `json:"plan_name"`
	PlanId      int64  `json:"plan_id"`
}
type AddGroupInfo struct {
	ID          int64  `json:"id"`
	GroupChatID int64  `json:"group_chat_id"`
	IsFinished  int    `json:"is_finished"`
	GroupName   string `json:"group_name"`
}
type Receives3 struct {
	ID          int64 `json:"id"`
	ConditionID int64 `json:"condition_id"`
	IsFinished  int   `json:"is_finished"`
}

type PlanInquireInfo struct {
	ID            int64  `json:"id"`
	TaskName      string `json:"task_name"` // 计划名称
	PlanStartTime int64  `json:"plan_start_time"`
	PlanEndTime   int64  `json:"plan_end_time"`
	Dimension     int8   `json:"dimension"` // 1社群数据；2用户画像
	Status        int8   `json:"status"`
}

// MatchGroupIDList 获取匹配的社群
type MatchGroupIDList struct {
	GroupChatID int64 `json:"group_chat_id"`
}

// MatchUnionID 获取匹配客户的union_id
type MatchUnionID struct {
	UnionID string `json:"union_id"`
}

type TaskName struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type TaskNumberCount struct {
	TaskID int64 `json:"task_id"`
	Number int64 `json:"number"`
}

type TaskGroupNumberCount struct {
	TaskID      int64 `json:"task_id"`
	TaskNumber  int64 `json:"task_number"`
	GroupNumber int64 `json:"group_number"`
	LordNumber  int64 `json:"lord_number"`
}
type TaskCustomerCount struct {
	TaskID            int64 `json:"task_id"`
	CustomerNumber    int64 `json:"customer_number"`
	DisCustomerNumber int64 `json:"dis_customer_number"`
}

// CustomerPortraitDetailResult 客户画像返回数据
type CustomerPortraitDetailResult struct {
	ID                       int64  `json:"id"`
	SubtasksNumber           string `json:"subtasks_number"`             // 子任务名称
	TaskType                 int64  `json:"task_type"`                   // 1单次,2持续
	ReachGuideNumber         int64  `json:"reach_guide_number"`          // 触达导购人数
	MatchCustomersNumber     int64  `json:"match_customers_number"`      // 匹配客户数
	MatchCustomersNumber2    int64  `json:"match_customers_number2"`     // 匹配客户数（去重）
	TaskNumber               int64  `json:"task_number"`                 // 任务总件数
	TaskFinishNumber         int64  `json:"task_finish_number"`          // 已完成任务数
	TaskFinishPercentage     string `json:"task_finish_percentage"`      // 任务完成率
	TotalReachCustomerNumber int64  `json:"total_reach_customer_number"` // 累计触达客户数
	TotalOrderCustomerNumber int64  `json:"total_order_customer_number"` // 累计下单客户数
	TotalOrderNumber         int64  `json:"total_order_number"`          // 累计订单数
	TotalOrderFee            int64  `json:"total_order_fee"`             // 累计交易金额
	Dimension                int    `json:"dimension"`                   // 子任务类型
}

// CommunityDataDetailResult 社群数据返回明细
type CommunityDataDetailResult struct {
	ID                       int64  `json:"id"`
	SubtasksNumber           string `json:"subtasks_number"`             // 子任务名称
	TaskType                 int64  `json:"task_type"`                   // 子任务类型
	ReachLordNumber          int64  `json:"reach_lord_number"`           // 触达群主数量
	MatchGroupNumber         int64  `json:"match_group_number"`          // 匹配社群数
	TaskNumber               int64  `json:"task_number"`                 // 任务总件数
	TaskFinishNumber         int64  `json:"task_finish_number"`          // 已完成任务数
	TaskFinishPercentage     string `json:"task_finish_percentage"`      // 任务完成率
	TotalReachCustomerNumber int64  `json:"total_reach_customer_number"` // 累计触达客户数
	TotalOrderCustomerNumber int64  `json:"total_order_customer_number"` // 累计下单客户数
	TotalOrderNumber         int64  `json:"total_order_number"`          // 累计订单数
	TotalOrderFee            int64  `json:"total_order_fee"`             // 累计交易金额
}

type SyncMemberData struct {
	UnionID    string `json:"union_id" `
	IsMember   uint8  `json:"is_member"`
	GradeCode  string `json:"grade_code"`
	MemberCode string `json:"member_code" `
	AdminBa    string `json:"admin_ba" `
}

type WxUserCorpExternalEuid struct {
	UnionID    string `json:"union_id" `
	GradeCode  string `json:"grade_code"`
	MemberCode string `json:"member_code" `
	AdminBa    string `json:"admin_ba" `
	Euid       int64  `json:"euid"`
}

type ResponseBirthday struct {
	PlanId         int64  `json:"plan_id"`
	PlanName       string `json:"plan_name"`
	TaskId         int64  `json:"task_id"`
	TaskName       string `json:"task_name"`
	ConditionId    int64  `json:"condition_id"`
	IsSpecifyCrowd int    `json:"is_specify_crowd"` // 任务是否指定人群 1 否 2 是
	// IsPush        int    `json:"is_push"`
	Conditions      string `json:"conditions"`
	StartTime       int64  `json:"start_time"`
	EndTime         int64  `json:"end_time"`
	ContentType     int    `json:"content_type"`
	Content         string `json:"content"`
	BirthdayPreTime int    `json:"birthday_pre_time"`
}

type ResponseTemplateCondition struct {
	PlanId         int64  `json:"plan_id"`
	PlanName       string `json:"plan_name"`
	TaskId         int64  `json:"task_id"`
	TaskName       string `json:"task_name"`
	PlanStatus     int    `json:"plan_status"`
	TaskStatus     int    `json:"task_status"`
	ConditionId    int64  `json:"condition_id"`
	IsSpecifyCrowd int    `json:"is_specify_crowd"` // 任务是否指定人群 1 否 2 是
	// IsPush        int    `json:"is_push"`
	Conditions     string `json:"conditions"`
	ConditionType  int    `json:"conditions"`
	StartTime      int64  `json:"start_time"`
	EndTime        int64  `json:"end_time"`
	ContentType    int    `json:"content_type"`
	Content        string `json:"content"`
	TaskType       int    `json:"task_type"`
	SendType       int    `json:"send_type"`
	TaskIsPush     int    `json:"task_is_push"`
	ConsumerFileID int    `json:"consumer_file_id"`
	IsExtract      int    `json:"is_extract"`
	TaskStartTime  int64  `json:"task_start_time"`
	TaskEndTime    int64  `json:"task_end_time"`
	TaskExpireType int    `json:"task_expire_type"`
}
