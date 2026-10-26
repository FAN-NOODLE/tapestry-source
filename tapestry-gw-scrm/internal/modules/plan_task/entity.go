package plantask

// PlanTasksAddRequst 新增子任务请求参数
type PlanTasksRequst struct {
	PlanId                  int64                    `json:"plan_id" binding:"required" desc:"子任务对应的计划的ID"`
	TaskId                  int64                    `json:"task_id"  desc:"子任务编号"`
	TaskType                int                      `json:"task_type"  desc:"任务周期 1 单次 2持续"`
	SendType                int                      `json:"send_type"  desc:"任务推送对象 1 BA 2管理bA"`
	Name                    string                   `json:"name" binding:"required" desc:"子任务名称"`
	Dimension               int                      `json:"dimension" binding:"required" desc:"子任务的维度；1社群数据；2用户画像"`
	ConditionType           int                      `json:"condition_type"  desc:"子任务的条件类型；1标签；2群人数；3群内新增任务；4参与活动；5购买行为；6群活跃度；7退群名单"`
	IsOpen                  int                      `json:"is_open" desc:"任务在小程序中是否可见。0不可见，1可见；默认可见"`
	Status                  int                      `json:"status"  desc:"任务状态；0未开始、1进行中、2已完成"`
	Remark                  string                   `json:"remark" binding:"required" desc:"子任务描述"`
	PlanTaskConditions      []map[string]interface{} `json:"plan_task_conditions" binding:"required" desc:"子任务条件与内容"`
	DeletPlanTaskConditions []int64                  `json:"delet_task_conditions"  desc:"子任务条件与内容"`
	PlanTaskSpecifyCrowd    []map[string]interface{} `json:"plan_task_specify_crowd" desc:"任务指定人群"`
	BirthdayPreTime         int                      `json:"birthday_pre_time" desc:"生日消息推送的提前天数"`
	ConsumerFileId          int                      `json:"consumer_file_id" desc:"导入客户模板的文件ID"`
	TaskCategory            int                      `json:"task_category" desc:"子任务的任务类型，沟通任务/朋友圈任务"`
	GroupSendType           int                      `json:"group_send_type" desc:"子任务群发类型，标签任务群发/企业任务群发"`
	IsAddTag                int                      `json:"is_add_tag" desc:"是否生成任务标签"`
	TagExpireTime           int                      `json:"tag_expire_time" desc:"标签任务的过期时间,0为永不过期"`
	TagName                 string                   `json:"tag_name" desc:"子任务的打标签的标签名称"`
	TagId                   string                   `json:"tag_id" desc:"企微的tag_id"`
	TagExpireType           int                      `json:"tag_expire_type" desc:"企微标签的过期类型,1任务结束后，2自定义删除时间，3长期有效"`
	FinishRate              int                      `json:"finish_rate"`
	TaskStartTime           int64                    `json:"task_start_time"`
	TaskEndTime             int64                    `json:"task_end_time"`
	TaskExpireType          int                      `json:"task_expire_type"`
	AllowSelect             int                      `json:"allow_select"`
}

// 任务详情
type PlanTasksInfoRequst struct {
	TaskId int64 `json:"task_id"  binding:"required" desc:"子任务编号"`
}
type PlanTasksResponse struct {
	PlanId             int64                       `json:"plan_id"`
	TaskId             int64                       `json:"task_id"`
	Name               string                      `json:"name"`
	Dimension          int                         `json:"dimension"`
	ConditionType      int                         `json:"condition_type"`
	IsOpen             int                         `json:"is_open"`
	Status             int                         `json:"status"`
	Remark             string                      `json:"remark"`
	CreatedAt          int64                       `json:"created_at"`
	UpdatedAt          int64                       `json:"updated_at"`
	PlanTaskConditions []PlanTaskConditionResponse `json:"plan_task_conditions"`
}
type PlanTaskConditionResponse struct {
	Id            int64                  `json:"id"`
	PlanId        int64                  `json:"plan_id"`
	TaskId        int64                  `json:"task_id"`
	Conditions    map[string]interface{} `json:"conditions"`
	Content       map[string]interface{} `json:"content"`
	ConditionType int                    `json:"condition_type"`
	CreatedAt     int64                  `json:"created_at"`
	UpdatedAt     int64                  `json:"updated_at"`
}

// PlanListRequest 计划列表请求参数
type PlanListRequest struct {
	Id             int64  `json:"id,omitempty"`
	Name           string `json:"name,omitempty"`
	SearchType     int    `json:"search_type,omitempty"`
	Status         int    `json:"status,omitempty"`
	StartTime      int    `json:"start_time,omitempty"`
	EndTime        int    `json:"end_time,omitempty"`
	CreateName     string `json:"create_name,omitempty"`
	EnterpriseId   int64  `json:"enterprise_id,omitempty"`
	TaskLabel      int32  `json:"task_label"`
	PlanTemplateId int64  `json:"plan_template_id"`
	PageParam
}

type PageParam struct {
	Page int `json:"page,omitempty"`
	Size int `json:"size,omitempty"`
}

// TaskListRequest 任务请求参数
type TaskListRequest struct {
	PlanId           int64  `json:"plan_id,omitempty"`
	TaskId           int64  `json:"task_id" comment:"子任务id"`
	Category         int32  `json:"category" comment:"任务类型，0沟通，1朋友圈"`
	TemplateId       int64  `json:"template_id" comment:"任务模板id"`
	TaskName         string `json:"task_name" comment:"子任务名称"`
	GroupType        int    `json:"group_type" comment:"0标签群发，1企业群发，2 1v1沟通任务"`
	PlanName         string `json:"plan_name"`
	CreatedStartTime int64  `json:"created_start_time" comment:"创建时间开始"`
	CreatedEndTime   int64  `json:"created_end_time" comment:"创建时间结束"`
	TaskStarTime     int64  `json:"task_star_time" comment:"开始时间"`
	TaskEndTime      int64  `json:"task_end_time" comment:"开始时间结束"`
	OriginType       int    `json:"origin_type"`
	PageParam
}

// 计划新增
type PlanRequest struct {
	PlanName        string `json:"plan_name"`
	Type            int    `json:"type"`
	StartTime       int64  `json:"start_time"`
	EndTime         int64  `json:"end_time"`
	Description     string `json:"description"`  // 描述
	WhetherOpen     bool   `json:"whether_open"` // 是否开启
	OperationUserId int64  `json:"operation_user_id"`
	Category        int32  `json:"category"`
}

// 计划编辑
type ReceivePlanModel struct {
	ID              int64  `json:"id"`
	CreatedAt       int64  `json:"created_at"`
	PlanName        string `json:"plan_name"`
	TaskNum         int    `json:"task_num"`
	OperationUserId int64  `json:"operation_user_id"`
	Type            int    `json:"type"`
	StartTime       int64  `json:"start_time"`
	EndTime         int64  `json:"end_time"`
	WhetherOpen     bool   `json:"whether_open"`
	Status          int    `json:"status"`
	Description     string `json:"description"`
	TaskLabel       int32  `json:"task_label"`
	Category        int32  `json:"category"`
}

type PlanStatus struct {
	ID     int64 `json:"id"`
	IsOpen bool  `json:"is_open"`
}
type PlanTasksConditionsRequst struct {
	Id                  int64  `json:"id"  binding:"required" desc:"条件编号"`
	PlanStartTime       int64  `json:"plan_start_time" desc:"任务开始时间"`
	PlanEndTime         int64  `json:"plan_end_time" desc:"任务结束时间"`
	DepartmentId        int64  `json:"department_id" desc:"部门编号"`
	IsFinished          int    `json:"is_finished" desc:"完成状态；1是、2否"`
	BaName              string `json:"ba_name" desc:"导购名称"`
	CompletionStartTime int64  `json:"completion_start_time" desc:"任务完成时间开始时间"`
	CompletionEndTime   int64  `json:"completion_end_time" desc:"任务完成时间结束时间"`
	PageParam
}
type PlanTasksConditionsAllotLogRespone struct {
	UserId         string `json:"user_id"`
	Name           string `json:"name"`
	DepartmentId   int64  `json:"department_id"`
	DepartmentName string `json:"department_name"`
	StartTime      string `json:"start_time"`
	EndTime        string `json:"end_time"`
	IsFinished     string `json:"is_finished"`
	CustomerNum    string `json:"customer_num"`
}
type PlanTasksConditionsBaRequst struct {
	Id            int64  `json:"id"  binding:"required" desc:"条件编号"`
	User_id       string `json:"user_id" binding:"required" desc:"导购编号"`
	Customer_name string `json:"customer_name" desc:"顾客名称"`
	PageParam
}
type PlanTasksConditionsBaRespone struct {
	TagName string `json:"tag_name"`
	Name    string `json:"name"`
	Mobile  string `json:"mobile"`
	Sex     string `json:"sex"`
}

type PlanId struct {
	ID int64 `json:"id"`
}

type UserInfoReq struct {
	ActivityUserId int64 `json:"activity_user_id"`
}

type UserListReq struct {
	ActivityId int64 `json:"activity_id"`
	StartTime  int64 `json:"start_time"`
	EndTime    int64 `json:"end_time"`
	Status     int64 `json:"status"`
	Page       int   `json:"page"`
	Size       int   `json:"size"`
}

type ExportUserDataReq struct {
	ActivityId int64 `json:"activity_id"`
	StartTime  int64 `json:"start_time"`
	EndTime    int64 `json:"end_time"`
	Status     int64 `json:"status"`
}

type UserModel struct {
	WxName          string `json:"wx_name"`
	UnionId         string `json:"union_id"`
	MemberCode      string `json:"member_code"`
	ReceiverName    string `json:"receiver_name"`
	ReceiverPhone   string `json:"receiver_phone"`
	ReceiverAddress string `json:"receiver_address"`
	StatusText      string `json:"status_text"`
	CreatedTimeText string `json:"created_time_text"`
}

type CustomerDataReq struct {
	PlanID    int64 `json:"plan_id"`
	TaskID    int64 `json:"task_id"`
	StartTime int64 `json:"start_time"`
	EndTime   int64 `json:"end_time"`
}

type ExportTasksRequst struct {
	Id                  int64  `json:"id"  binding:"required" desc:"条件编号"`
	PlanStartTime       int64  `json:"plan_start_time" desc:"任务开始时间"`
	PlanEndTime         int64  `json:"plan_end_time" desc:"任务结束时间"`
	DepartmentId        int64  `json:"department_id" desc:"部门编号"`
	IsFinished          int    `json:"is_finished" desc:"完成状态；1是、2否"`
	BaName              string `json:"ba_name" desc:"导购名称"`
	CompletionStartTime int64  `json:"completion_start_time" desc:"任务完成时间开始时间"`
	CompletionEndTime   int64  `json:"completion_end_time" desc:"任务完成时间结束时间"`
	EndTime             int64  `json:"end_time"`
	TaskType            int    `json:"task_type"`
}

type TaskListRequst struct {
	ExportTasksRequst
	PageParam
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
	WhetherToResign          int64  `json:"whether_to_resign"`
	TotalReachCustomerNumber int64  `json:"total_reach_customer_number"` // 累计触达客户数
	TotalOrderCustomerNumber int64  `json:"total_order_customer_number"` // 累计下单客户数
	TotalOrderNumber         int64  `json:"total_order_number"`          // 累计订单数
	TotalOrderFee            int64  `json:"total_order_fee"`             // 累计交易金额
	Category                 int    `json:"category"`
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
	GroupChatName            string `json:"group_chat_name"`
	BelongingArea            string `json:"belonging_area"`
	DepartmentName           string `json:"department_name"`
	TotalReachCustomerNumber int64  `json:"total_reach_customer_number"` // 累计触达客户数
	TotalOrderCustomerNumber int64  `json:"total_order_customer_number"` // 累计下单客户数
	TotalOrderNumber         int64  `json:"total_order_number"`          // 累计订单数
	TotalOrderFee            int64  `json:"total_order_fee"`             // 累计交易金额
}

// PlanStatRequest 任务统计接收参数
type PlanStatRequest struct {
	EnterpriseID int64 `json:"enterprise_id"`
	StartTime    int64 `json:"start_time"`
	EndTime      int64 `json:"end_time"`
	PlanId       int64 `json:"plan_id"`
}

// StatDetailRequest 任务统计明细参数
type StatDetailRequest struct {
	StartTime      int64  `json:"start_time"`
	EndTime        int64  `json:"end_time"`
	PlanId         int64  `json:"plan_id"`
	EnterpriseID   int64  `json:"enterprise_id"`
	SubtasksNumber string `json:"subtasks_number"`       // 子任务名称（模糊查询）
	TaskType       int8   `json:"task_type" default:"1"` // 任务类型 1：社群数据 2：用户画像，
	PageParam
}

// CustomerPortraitDetailResult 客户画像返回数据
type CustomerPortraitDetailResult struct {
	ID                       int64  `json:"id"`
	SubtasksNumber           string `json:"subtasks_number"`             // 子任务名称
	TaskType                 int64  `json:"task_type"`                   // 子任务类型
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

type ResultCommunityDataDetail struct {
	Count int                          `json:"count"`
	Data  []*CommunityDataDetailResult `json:"data"`
}

type ResultCustomerPortraitDetail struct {
	Count int                             `json:"count"`
	Data  []*CustomerPortraitDetailResult `json:"data"`
}

type ConsumerList struct {
	ConsumerFileID int64 `json:"consumer_file_id,omitempty"`
	EnterpriseID   int64 `json:"enterprise_id"`
	Total          int64 `json:"total"`
	TaskType       int   `binding:"required" json:"task_type"`
	PageParam
}

type ResultConsumerData struct {
	ID        int64  `json:"id"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
	DeletedAt int64  `json:"deleted_at"`
	FileID    int64  `json:"file_id"`
	TaskID    int64  `json:"task_id"`
	UnionId   string `json:"unionid"`
	Said      string `json:"said"`
	Name      string `json:"name"`
	ImportRes string `json:"import_res"`
}

type ResultConsumerList struct {
	Total int64 `json:"total"`
	List  []ResultConsumerData
}
