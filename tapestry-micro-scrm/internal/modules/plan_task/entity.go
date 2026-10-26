package plan_task

import (
	"micro-scrm/internal/orm/model/micro_scrm_enterprise"
)

// PlanTasksRequest 实体
// PlanTasksAddRequest 新增子任务请求参数
type PlanTasksRequest struct {
	PlanId                  int64                    `json:"plan_id" binding:"required" desc:"主任务ID"`
	TaskId                  int64                    `json:"task_id"  desc:"子任务编号"`
	TaskType                int                      `json:"task_type"  desc:"任务周期 1 单次 2持续"`
	SendType                int                      `json:"send_type"  desc:"任务推送对象 1 BA 2管理bA"`
	Name                    string                   `json:"name" binding:"required" desc:"子任务名称"`
	Dimension               int                      `json:"dimension" binding:"required" desc:"子任务的维度；1社群数据；2用户画像"`
	ConditionType           int                      `json:"condition_type" binding:"required" desc:"子任务的条件类型；1标签；2群人数；3群内新增任务；4参与活动；5购买行为；6群活跃度；7退群名单"`
	IsOpen                  int                      `json:"is_open" desc:"任务在小程序中是否可见。0不可见，1可见；默认可见"`
	Status                  int                      `json:"status"  desc:"任务状态；0未开始、1进行中、2已完成"`
	Remark                  string                   `json:"remark" binding:"required" desc:"子任务描述"`
	PlanTaskConditions      []map[string]interface{} `json:"plan_task_conditions" binding:"required" desc:"子任务条件与内容"`
	DeletPlanTaskConditions []int64                  `json:"delet_task_conditions"  desc:"子任务条件与内容"`
	PlanTaskSpecifyCrowd    []PlanTaskSpecifyCrowd   `json:"plan_task_specify_crowd" desc:"任务指定人群"`
	BirthdayPreTime         int                      `json:"birthday_pre_time" desc:"生日消息推送的提前天数"`
	ConsumerFileId          int                      `json:"consumer_file_id" desc:"导入客户模板的文件id"`
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
type PlanTaskConditions struct {
	Id              int64                                 `json:"id"`
	Deleted_at      int                                   `json:"deleted_at"`
	ContentType     int                                   `json:"content_type"`
	Content         micro_scrm_enterprise.Content         `json:"content"`
	ConditionDetail micro_scrm_enterprise.ConditionDetail `json:"conditions"`
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
	Mobile         string `json:"mobile"`
	DateTime       int64  `json:"date_time"`
	TaskLabel      int32  `json:"task_label"`
	PlanTemplateId int64  `json:"plan_template_id"`
	PageParam
}

// PlanListDetail 计划返回详情
type PlanListDetail struct {
	Id              int64  `json:"id"`
	Name            string `json:"name"`
	TaskNum         int    `json:"task_num"`
	OperationUserId int64  `json:"operation_user_id"`
	Type            int    `json:"type"`
	StartTime       int    `json:"start_time"`
	EndTime         int    `json:"end_time"`
	IsOpen          bool   `json:"is_open"`
	Status          int    `json:"status"`
	Remark          string `json:"remark"`
	CreatedAt       int64  `json:"created_at"`
	UpdatedAt       int64  `json:"updated_at"`
	DeletedAt       int64  `json:"deleted_at"`
	CreateName      string `json:"create_name"`
	Category        int32  `json:"category"`
	TemplateId      int64  `json:"template_id"`
}

// PlanListResponse 计划响应列表
type PlanListResponse struct {
	PlanList []PlanListDetail `json:"plan_list"`
}

// TaskListRequest 任务请求参数
type TaskListRequest struct {
	PlanId           int64  `json:"plan_id,omitempty"`
	TaskId           int64  `json:"task_id" comment:"子任务id"`
	Category         int    `json:"category" comment:"任务类型，0沟通，1朋友圈"`
	TemplateId       int64  `json:"template_id" comment:"任务模板id"`
	TaskName         string `json:"task_name" comment:"子任务名称"`
	GroupType        int    `json:"group_type" comment:"0标签群发，1企业群发，2 1v1沟通任务"`
	PlanName         string `json:"plan_name"`
	CreatedStartTime int64  `json:"created_start_time" comment:"创建时间开始"`
	CreatedEndTime   int64  `json:"created_end_time" comment:"创建时间结束"`
	TaskStarTime     int64  `json:"task_star_time" comment:"开始时间"`
	TaskEndTime      int64  `json:"task_end_time" comment:"开始时间结束"`
	PageParam
}

type TaskListDetail struct {
	Id             int64  `json:"id"`
	PlanId         int64  `json:"plan_id"`
	Name           string `json:"name"`
	Dimension      int    `json:"dimension"`
	ConditionType  int    `json:"condition_type"`
	IsOpen         bool   `json:"is_open"`
	Status         int    `json:"status"`
	Remark         string `json:"remark"`
	CreatedAt      int64  `json:"created_at"`
	UpdatedAt      int64  `json:"updated_at"`
	DeletedAt      int64  `json:"deleted_at"`
	TaskCategory   int    `json:"task_category"`
	GroupSendType  int    `json:"group_send_type"`
	FinishRate     int    `json:"finish_rate"`
	TaskStartTime  int64  `json:"task_start_time"`
	TaskEndTime    int64  `json:"task_end_time"`
	PlanTemplateId int64  `json:"plan_template_id"`
	OriginType     int    `json:"origin_type"`
}

type TaskListResponse struct {
	TaskList []TaskListDetail `json:"task_list"`
}

type PageParam struct {
	Page int `json:"page,omitempty"`
	Size int `json:"size,omitempty"`
}

// 任务详情
type PlanTasksInfoRequst struct {
	TaskId int64 `json:"task_id"  binding:"required" desc:"子任务编号"`
}

// 任务条件
type PlanTasksConRequest struct {
	Id           int64 `json:"id"`
	TaskId       int64 `json:"task_id"`
	EnterpriseId int64 `json:"enterprise_id"`
	MemberId     int64 `json:"member_id"`
}

type PlanTasksResponse struct {
	micro_scrm_enterprise.PlanTaskModel
	PlanId               int64                       `json:"plan_id"`
	TaskId               int64                       `json:"task_id"`
	TaskType             int                         `json:"task_type"  desc:"任务周期 1 单次 2持续"`
	SendType             int                         `json:"send_type"  desc:"任务推送对象 1 BA 2管理bA"`
	Name                 string                      `json:"name"`
	IsSpecifyCrowd       int                         `json:"is_specify_crowd"` // 任务是否指定人群 1 否 2 是
	Dimension            int                         `json:"dimension"`
	ConditionType        int                         `json:"condition_type"`
	IsOpen               int                         `json:"is_open"`
	Status               int                         `json:"status"`
	Remark               string                      `json:"remark"`
	CreatedAt            int64                       `json:"created_at"`
	UpdatedAt            int64                       `json:"updated_at"`
	PlanTaskConditions   []PlanTaskConditionResponse `json:"plan_task_conditions"`
	PlanTaskSpecifyCrowd []PlanTaskSpecifyCrowd      `json:"plan_task_specify_crowd" desc:"任务指定人群"`
	BirthdayPreTime      int                         `json:"birthday_pre_time" desc:"生日类型任务提前多少天推送消息"`
	ConsumerFileId       int                         `json:"consumer_file_id" desc:"客户导入类型的任务名单文件id"`
	IsExtract            int                         `json:"is_extract" desc:"客户名单是否导入完成0:未完成,1已完成"`
	DownLoadUrl          string                      `json:"down_load_url"  desc:"客户名单导入的匹配结果cos地址"`
	TaskCategory         int                         `json:"task_category"`
	GroupSendType        int                         `json:"group_send_type"`
	Tags                 []string                    `json:"tags"`
	FileName             string                      `json:"file_name"`
	FileUrl              string                      `json:"file_url"`
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

type ReceivePlan struct {
	PlanName        string `json:"plan_name"`
	StartTime       int64  `json:"start_time"`
	EndTime         int64  `json:"end_time"`
	Description     string `json:"description"`  // 描述
	WhetherOpen     bool   `json:"whether_open"` // 是否开启
	OperationUserId int64  `json:"operation_user_id"`
	Category        int32  `json:"category"`
}

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
	Category        int32  `json:"category"`
}

type PlanId struct {
	ID int64 `json:"id"`
}

// 任务挑起
type PlanTasksConditionsRequst struct {
	Id                  int64  `json:"id"  binding:"required" desc:"条件编号"`
	PlanStartTime       int64  `json:"plan_start_time" desc:"任务开始时间"`
	PlanEndTime         int64  `json:"plan_end_time" desc:"任务结束时间"`
	DepartmentId        int64  `json:"department_id" desc:"部门编号"`
	IsFinished          int    `json:"is_finished" desc:"完成状态；1是、2否"`
	BaName              string `json:"ba_name" desc:"导购名称"`
	CompletionStartTime int64  `json:"completion_start_time" desc:"任务完成时间开始时间"`
	CompletionEndTime   int64  `json:"completion_end_time" desc:"任务完成时间结束时间"`
	EndTime             int64  `json:"end_time" desc:"订单统计结束时间"`
	PageParam
}
type PlanTasksConditionsReponse struct {
	UserId string `json:"user_id"`
	Name   string `json:"name"`
}
type PlanTasksConditionsAllotLogRespone struct {
	UserId         string `json:"user_id"`
	Name           string `json:"name"`
	DepartmentId   int64  `json:"department_id"`
	DepartmentName string `json:"department_name"`
	StartTime      string `json:"start_time"`
	EndTime        string `json:"end_time"`
	IsFinished     int    `json:"is_finished"`
	CustomerNum    int    `json:"customer_num"`
}
type PlanTasksConditionsBaRequst struct {
	Id            int64  `json:"id"  binding:"required" desc:"条件编号"`
	User_id       string `json:"user_id" desc:"导购编号"`
	Customer_name string `json:"customer_name" desc:"顾客名称"`
	PageParam
}
type PlanStatus struct {
	ID     int64 `json:"id"`
	IsOpen bool  `json:"is_open"`
}
type PlanTasksConditionsBaRespone struct {
	TagName   string `json:"tag_name"`
	Name      string `json:"name"`
	Mobile    string `json:"mobile"`
	Sex       string `json:"sex"`
	Avatar    string `json:"avatar"`
	LeaveTime int64  `json:"leave_time"`
}

type TaskConditionReq struct {
	ID           int64  `json:"id"  binding:"required"`
	EnterpriseId int64  `json:"enterprise_id"`
	TagName      string `json:"tag_name"`
	Mobile       string `json:"mobile"`
	ConditionId  int64  `json:"condition_id"`
	Status       int    `json:"status"`
	IsUnion      int    `json:"is_union"`
	GroupChatID  int64  `json:"group_chat_id"`
	Dimension    int64  `json:"dimension"`
	PageParam
}

type NewTaskConditionReq struct {
	ID           int64  `json:"id"  binding:"required"`
	EnterpriseId int64  `json:"enterprise_id"`
	TagName      string `json:"tag_name"`
	Mobile       string `json:"mobile"`
	ConditionId  int64  `json:"condition_id"`
	Status       int    `json:"status"`
	IsUnion      int    `json:"is_union"`
	ChatGroupID  int64  `json:"chat_group_id"`
	PageParam
}

type ConditionJobChan struct {
	PlanTaskConditionModel micro_scrm_enterprise.PlanTaskConditionModel `json:"plan_task_condition_model"`
	EnterpriseID           int64                                        `json:"enterprise_id"`
	TagName                string                                       `json:"tag_name"`
	TagId                  int64                                        `json:"tag_id"`
	PlanId                 int64                                        `json:"plan_id_list"`
	SendType               int                                          `json:"send_type"`
	IsSpecifyCrowd         int                                          `json:"is_specify_crowd"` // 任务是否指定人群 1 否 2 是
}

type NewResponseFrontPlanTaskConList struct {
	Id                int64                                 `json:"id"`
	TaskId            int64                                 `json:"task_id"`
	ConditionInfo     micro_scrm_enterprise.ConditionDetail `json:"condition_info"`
	ContentInfo       micro_scrm_enterprise.Content         `json:"content_info"`
	ContentType       int                                   `json:"content_type"`
	IsFinished        int                                   `json:"is_finished"`
	Dimension         int                                   `json:"dimension"`
	GroupChatID       int64                                 `json:"group_chat_id"`
	ConditionType     int                                   `json:"condition_type"`
	Remark            string                                `json:"remark"`
	ConditionID       int64                                 `json:"condition_id"`
	GroupInfo         []*GroupInfo                          `json:"group_info"`
	LeaveDetail       []*LeaveDetail                        `json:"leave_detail"`
	LeaveNumber       int64                                 `json:"leave_number"`
	CustomerChurnRate string                                `json:"customer_churn_rate"`
	StatusCode        int                                   `json:"status_code" desc:"子任务的完成状态"`
	IsForward         int                                   `json:"is_forward"`
	NoFinishedNum     int64                                 `json:"no_finished_num"`
	TaskCategory      int                                   `json:"task_category"`
	GroupSendType     int                                   `json:"group_send_type"`
	IsAddTag          int                                   `json:"is_add_tag"`
	Tags              []string                              `json:"tags"`
	ExpireStatus      int32                                 `json:"expire_status" desc:"主任务的过期状态"`
	PlanId            int64                                 `json:"plan_id"`
	PlanName          string                                `json:"plan_name"`
	FinishRate        int                                   `json:"finish_rate"`
	TaskStartTime     int64                                 `json:"task_start_time"`
	TaskEndTime       int64                                 `json:"task_end_time"`
	TaskExpireType    int                                   `json:"task_expire_type"`
	CustomerNum       int                                   `json:"customer_num"`
}

type GroupInfo struct {
	ID          int64  `json:"id"`
	GroupName   string `json:"group_name"`
	GroupChatID int64  `json:"group_chat_id"`
	IsFinish    int    `json:"is_finish"`
	ChatID      string `json:"chat_id"`
}
type LeaveDetail struct {
	Name               string `json:"name"`
	Gender             string `json:"gender"`
	WorkExternalUserId string `json:"work_external_user_id"`
	Avatar             string `json:"avatar"`
}
type UpdateReceive struct {
	ID           int64 `json:"id"`
	GroupChatID  int64 `json:"group_chat_id"`
	Status       int64 `json:"status"`
	EnterpriseID int64 `json:"enterprise_id"`
}

type ResponseFrontPlanList struct {
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	Remark    string `json:"remark"`
	SendType  int    `json:"send_type"`
	Status    int    `json:"status"`
	StatusTex string `json:"status_tex"`
	IsOpen    int    `json:"is_open"`
	StartTime int64  `json:"start_time"`
	EndTime   int64  `json:"end_time"`
	CreatedAt int64  `json:"created_at"`
}

type GetFileInfo struct {
	TagName        string      `json:"tag_name"`
	EnterpriseID   int64       `json:"enterprise_id"`
	ID             int64       `json:"id"`
	FileName       string      `json:"file_name"`
	File           []FileField `json:"file"`             // csv 文件
	UnionIDEncrypt string      `json:"union_id_encrypt"` // union是否加密 1.加密过 2.未加密
}

type FileField struct {
	Ccode     string `json:"ccode"`
	Unionid   string `json:"unionid"`
	GradeCode string `json:"grade_code"`
	AdminBA   string `json:"admin_ba"`
}

type DownLoadFileInfo struct {
	EnterpriseID int64  `json:"enterprise_id"`
	FileName     string `json:"file_name"`
	TagName      string `json:"tag_name"` // 标签名称
	FileUrl      string `json:"file_url"`
	TagGroupName string `json:"tag_group_name"` // 标签组
	BAType       string `json:"ba_type"`        // BA类型  1.普通BA 2.管理BA
}

type IDInfo struct {
	ID int64 `json:"id"`
}

type CountInfo struct {
	Count int64 `json:"id"`
}

type FileDetection struct {
	EnterpriseID int64  `json:"enterprise_id"`
	FileName     string `json:"file_name"`
}

type UpdateIdInfo struct {
	ID           int64 `json:"id"`
	EnterpriseID int64 `json:"enterprise_id"`
}

type TagInfo struct {
	GroupName string
	TagName   string
	BaType    uint8
}

type PlanTaskSpecifyCrowd struct {
	UserId   string `json:"user_id"`   // 推送的用户编号
	UserName string `json:"user_name"` // 用户/部门名称
	ThirdId  int    `json:"third_id"`  //企微部门id
}

type CustomerDataReq struct {
	PlanID    int64 `json:"plan_id"`
	TaskID    int64 `json:"task_id"`
	StartTime int64 `json:"start_time"`
	EndTime   int64 `json:"end_time"`
}

type CustomerDataResp struct {
	TouchBaCount             int64  `json:"touch_ba_count"`
	TouchMasterCount         int64  `json:"touch_master_count"`
	GroupChatCount           int64  `json:"group_chat_count"`
	CustomerCount            int64  `json:"customer_count"`
	CustomerDistictCount     int64  `json:"customer_distict_count"`
	TaskCount                int64  `json:"task_count"`
	TaskCompleteCount        int64  `json:"task_complete_count"`
	TaskCompleteRate         string `json:"task_complete_rate"`
	NewTaskCount             int64  `json:"new_task_count"`
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

type TaskFinishStatus struct {
	ID         int64 `json:"id"`
	ISFinished int8  `json:"is_finished"`
}

// PlanStatResultData 任务统计返回参数
type PlanStatResultData struct {
	TaskName                 string `json:"task_name"`       // 计划名称
	SubtasksNumber           int64  `json:"subtasks_number"` // 子任务数
	PlanStartTime            int64  `json:"plan_start_time"`
	PlanEndTime              int64  `json:"plan_end_time"`
	PeriodSubtasksNumber     int64  `json:"period_subtasks_number"`      // 时间段内的子任务数
	SubtasksFinishNumber     int64  `json:"subtasks_finish_number"`      // 子任务完成数
	SubtasksFinishPercentage string `json:"subtasks_finish_percentage"`  // 子任务完成率
	TaskNumber               int64  `json:"task_number"`                 // 任务总件数
	TaskFinishNumber         int64  `json:"task_finish_number"`          // 已完成任务数
	TaskFinishPercentage     string `json:"task_finish_percentage"`      // 任务完成率
	ReachGuideNumber         int64  `json:"reach_guide_number"`          // 触达导购人数
	MatchCustomersNumber     int64  `json:"match_customers_number"`      // 匹配客户数
	MatchCustomersNumber2    int64  `json:"match_customers_number2"`     // 匹配客户数（去重）
	ReachLordNumber          int64  `json:"reach_lord_number"`           // 触达群主数量
	MatchGroupNumber         int64  `json:"match_group_number"`          // 匹配社群数
	TotalReachCustomerNumber int64  `json:"total_reach_customer_number"` // 累计触达客户数
	TotalOrderCustomerNumber int64  `json:"total_order_customer_number"` // 累计下单客户数
	TotalOrderNumber         int64  `json:"total_order_number"`          // 累计订单数
	TotalOrderFee            int64  `json:"total_order_fee"`             // 累计交易金额
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
type ResultCommunityDataDetail struct {
	Count int                                                `json:"count"`
	Data  []*micro_scrm_enterprise.CommunityDataDetailResult `json:"data"`
}

type ResultCustomerPortraitDetail struct {
	Count int                                                   `json:"count"`
	Data  []*micro_scrm_enterprise.CustomerPortraitDetailResult `json:"data"`
}
type PlanTaskTag struct {
	Tag []TaskTag `json:"tag"`
}
type TaskTag struct {
	GroupName string `json:"group_name"`
	GroupID   int64  `json:"group_id"`
	QwGroupId string `json:"qw_group_id"`
	TagName   string `json:"tag_name"`
	TagID     int64  `json:"tag_id"`
	QwTagID   string `json:"qw_tag_id"`
	IsDelete  bool   `json:"is_delete"`
}

// UploadUserResponse 上传用户列表响应
type UploadUserResponse struct {
	LineNum int    `json:"line_num"`
	Item    string `json:"item" desc:"员工的userid或者marscode"`
	UserID  string `json:"user_id" desc:"员工的userid"`
	Name    string `json:"user_name" desc:"员工名称"`
	ErrType int    `json:"err_type" desc:"错误类型：0-正常，1-错误，2-重复"`
	Error   string `json:"error"`
}

type TaskEverydayStatistics struct {
	Date                     string `json:"date" db:"date"`
	DataTime                 int64  `json:"data_time" db:"data_time" desc:"数据更新时间"`
	TaskID                   int64  `json:"task_id" db:"task_id"`
	TotalReachCustomerNumber int64  `json:"total_reach_customer_number" db:"total_reach_customer_number"`
	TotalOrderCustomerNumber int64  `json:"total_order_customer_number" db:"total_order_customer_number"`
	TotalOrderNumber         int64  `json:"total_order_number" db:"total_order_number"`
	TotalOrderFee            int64  `json:"total_order_fee" db:"total_order_fee"`
	TaskType                 int64  `json:"task_type" db:"task_type"`
	CreatedAt                int64  `json:"created_at" db:"created_at"`
	UpdatedAt                int64  `json:"updated_at" db:"updated_at"`
}

type MatchUnion struct {
	TaskID  int64  `json:"task_id" db:"task_id"`
	UnionID string `json:"union_id" db:"union_id"`
}

type OrderStatistic struct {
	TotalOrderCustomerNumber int64 `json:"total_order_customer_number" db:"total_order_customer_number"`
	TotalOrderNumber         int64 `json:"total_order_number" db:"total_order_number"`
	TotalOrderFee            int64 `json:"total_order_fee" db:"total_order_fee"`
}

type TaskAllotLog struct {
	ID          int64  `json:"id" db:"id"`
	CreatedAt   int64  `json:"created_at" db:"created_at"`
	UpdatedAt   int64  `json:"updated_at" db:"updated_at"`
	DeletedAt   int64  `json:"deleted_at" db:"deleted_at"`
	PlanId      int64  `json:"plan_id" db:"plan_id"`           // 计划ID
	ConditionId int64  `json:"condition_id" db:"condition_id"` // 条件对应的子条件的ID
	TaskId      int64  `json:"task_id" db:"task_id"`           // 条件对应的子任务的ID
	UserId      string `json:"user_id" db:"user_id"`           // 员工的企微user_id
	StartTime   int    `json:"start_time" db:"start_time"`     // 任务开始时间
	EndTime     int    `json:"end_time" db:"end_time"`         // 任务完成时间
	IsFinished  int    `json:"is_finished" db:"is_finished"`   // 完成状态；0否、1是
	CustomerNum int    `json:"customer_num" db:"customer_num"` // 符合条件的客户数
	GroupChatId int64  `json:"chat_id" db:"chat_id"`
	TagId       int64  `json:"tag_id" db:"tag_id"`     // 标签编号
	TagName     string `json:"tag_name" db:"tag_name"` // 标签名称
}
