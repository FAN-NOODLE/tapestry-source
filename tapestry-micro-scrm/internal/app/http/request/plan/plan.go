package planReq

import (
	"micro-scrm/internal/app/http/request"
	customerSearchReq "micro-scrm/internal/app/http/request/customer_search"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
)

type BaseReq struct {
	EnterpriseId int64 `json:"enterprise_id" binding:"required" comment:"企业ID"`
}

type MonthListsReq struct {
	BaseReq
	StartTime int64  `json:"start_time" binding:"required" comment:"开始时间"`
	EndTime   int64  `json:"end_time" binding:"required" comment:"结束时间"`
	Mobile    string `json:"mobile"`
	Type      int    `json:"type" description:"类型:1任务计划,2播种计划,3邀约计划"`
	PlanId    int64  `json:"plan_id" description:"类型:针对于计划任务的plan_id"`
}

type MonthListsStatus struct {
	Finished bool `json:"finished"`
	Afoot    bool `json:"afoot"`
}

type MonthListsRes struct {
	List []MonthListsData
}
type MonthListsData struct {
	Time int64 `json:"time"`
	MonthListsStatus
}

type PlanAndTaskAllotFinishList struct {
	Id         int64 `json:"id"`
	IsFinished int   `json:"is_finished"`
	CreatedAt  int64 `json:"created_at"`
}

type ListsReq struct {
	BaseReq
	StartTime int64  `json:"start_time" binding:"required" comment:"开始时间"`
	EndTime   int64  `json:"end_time" binding:"required" comment:"结束时间"`
	Name      string `json:"name"`
	Mobile    string `json:"mobile"`
	Status    int    `json:"status"` // 任务状态，0为待完成,1为已完成,2为未完成,该状态以task_allot_log表is_finished状态为准,-1查全部、
}

type PlansListInfo struct {
	Id         int64  `json:"id"`
	Name       string `json:"name"`
	StartTime  int64  `json:"start_time"`
	EndTime    int64  `json:"end_time"`
	Remark     string `json:"remark"`
	IsFinished int    `json:"is_finished"`
	Status     int    `json:"status"` // 任务状态，0为待完成,1为已完成,2为未完成,该状态以task_allot_log表is_finished状态为准,-1查全部
	AllNum     int    `json:"all_num" comment:"子任务总数"`
	FinishNum  int    `json:"finish_num" comment:"任务的完成数"`
	TaskType   int    `json:"task_type" comment:"持续任务，单次任务的展示；1单次任务，2持续任务"`
	Category   int    `json:"category" comment:"任务类型，2是MA任务(模板创建)，1是手动创建任务"`
	CloseTips  int    `json:"close_tips" comment:"已完成的企业群发类型的沟通任务有未联系的顾客需要提示"`
}

type PlanInfo struct {
	BaseReq
	Id        int64  `json:"id"`
	Mobile    string `json:"mobile"`
	StartTime int64  `json:"start_time"`
	EndTime   int64  `json:"end_time"`
}

type PlanInfoRes struct {
	Id              int64  `json:"id"`
	Name            string `json:"name"`
	StartTime       int64  `json:"start_time"`
	EndTime         int64  `json:"end_time"`
	Remark          string `json:"remark"`
	TaskNum         string `json:"task_num"`
	TaskType        int    `json:"task_type" comment:"主任务下是否有持续性任务，0为否，1为是"`
	Status          int    `json:"status"  comment:"主任务的status,1为任务进行中,2为任务已结束"`
	Category        int    `json:"category" comment:"任务的模板类型"`
	PlanTemplateID  int64  `json:"plan_template_id" comment:"任务模板id"`
	StatusCode      int    `json:"status_code" comment:"任务的完成状态，0待完成，1已完成，2已过期"`
	FinishNum       int64  `json:"finish_num" comment:"任务的完成数"`
	OperationUserID int    `json:"operation_user_id"`
}

type PlanListRes struct {
	Lists []PlansListInfo
}

type TipsReq struct {
	BaseReq
	StartTime int64  `json:"start_time" binding:"required" comment:"开始时间"`
	EndTime   int64  `json:"end_time" binding:"required" comment:"结束时间"`
	Mobile    string `json:"mobile"`
}
type TipsRes struct {
	PlanNums    int64 `json:"plan_nums"`
	SowTaskNums int64 `json:"sow_task_nums"`
	ReserveNums int64 `json:"reserve_nums"`
}

type LeaveGroupChatReq struct {
	BaseReq
	Id     int64 `json:"id"`
	Status int   `json:"status" comment:"0为全部顾客,1是好友关系的顾客,2非好友关系的顾客"`
	Page   int   `json:"page" binding:"min=1" comment:"页码"`
	Size   int   `json:"size" binding:"min=10,max=99999999" comment:"分页数"`
}

type TaskAllotGroupInfo struct {
	PlanId    int64  `json:"plan_id"`
	StartTime int64  `json:"start_time"`
	EndTime   int64  `json:"end_time"`
	UserId    string `json:"user_id"`
}

type AllRes struct {
	UserId   string `json:"user_id"`
	AllCount int64  `json:"all_count"`
}

type RatioRes struct {
	UserId string  `json:"user_id"`
	Ratio  float64 `json:"ratio"`
}

type CustomerListsReq struct {
	BaseReq
	Id           int64  `json:"id"`
	Page         int    `json:"page" binding:"min=1" comment:"页码"`
	Size         int    `json:"size" binding:"min=10,max=99999999" comment:"分页数"`
	FilterType   int    `json:"filterType"`
	Keywords     string `json:"keywords"`
	IsUnion      int    `json:"is_union"`
	IsTag        int    `json:"is_tag" comment:"0为不展示tag,1展示tag"`
	TaskCategory int    `json:"task_category" comment:"子任务类型"`
}

type CustomerListsRes struct {
	microScrmModel.TaskConditionDataModel
	CustomerInfoRes
	TagName []string `json:"tag_name,omitempty"`
}

type CustomerInfoRes struct {
	Name                 string `json:"name"`
	Avatar               string `json:"avatar"`
	RelRemark            string `json:"rel_remark"`
	VipType              int    `json:"vip_type"`
	ValidReserveAt       int    `json:"valid_reserve_at"`
	SuccessReserveAt     int    `json:"success_reserve_at"`
	TaskFinishAt         int    `json:"task_finish_at"`
	CouponStatus         int    `json:"coupon_status"`
	CouponBy30daysStatus int    `json:"coupon_by_30_days_status"`
	CustomerCard         string `json:"customer_card"`
}

type ConfirmReq struct {
	BaseReq
	TaskAllotLogId int64  `json:"task_allot_log_id"`
	MemberId       int64  `json:"member_id"`
	Mobile         string `json:"mobile"`
	UserId         string `json:"user_id"`
	UnionId        string `json:"union_id"`
	CustomerId     string `json:"customer_id"`
	Status         int    `json:"status"  binding:"required" comment:"参数有误"`
}

type SetRemark struct {
	BaseReq
	TaskAllotLogId int64  `json:"task_allot_log_id"`
	MemberId       int64  `json:"member_id"`
	UserId         string `json:"user_id"`
	UnionId        string `json:"union_id"`
	CustomerId     string `json:"customer_id"`
	Remark         string `json:"remark"`
}

type ConditionDataTagsList struct {
	Id      int    `json:"id"`
	TagName string `json:"tag_name"`
	UnionId string `json:"union_id"`
	WwTagId string `json:"ww_tag_id"`
}

type SetForward struct {
	BaseReq
	Mobile string `json:"mobile"`
	Id     int64  `json:"id"`
}

type TaskDetailExportReq struct {
	BaseReq
	Id                  int64  `json:"id"`
	PlanStartTime       int64  `json:"plan_start_time" desc:"任务开始时间"`
	PlanEndTime         int64  `json:"plan_end_time" desc:"任务结束时间"`
	DepartmentId        int64  `json:"department_id" desc:"部门编号"`
	IsFinished          int    `json:"is_finished" desc:"完成状态；1是、2否,-1全部"`
	BaName              string `json:"ba_name" desc:"导购名称"`
	CompletionStartTime int64  `json:"completion_start_time" desc:"任务完成时间开始时间"`
	CompletionEndTime   int64  `json:"completion_end_time" desc:"任务完成时间结束时间"`
	EndTime             int64  `json:"end_time" desc:"task_allot_log表的完成时间"`
	TaskType            int    `json:"task_type"`
}

type TaskAllotDetailExportBase struct {
	PlanId          int64  `json:"plan_id"`
	PlanName        string `json:"plan_name" desc:"任务计划名称"`
	TaskName        string `json:"task_name" desc:"子任务名称"`
	Category        int    `json:"category" desc:"任务标签"`
	MemberLists     []TaskAllotMemberLists
	TaskAllotLogIds []int64
	TaskForwards    map[int64]int
	TaskFinish      map[int64]TaskAllotFinish
	TaskCategory    int `json:"task_category" desc:"子任务类型 沟通任务，朋友圈任务"`
	GroupSendType   int `json:"group_send_type" desc:"发送类型 标签任务群发/企业任务群发"`
	TaskFinishRate  map[int64]int
}
type TaskAllotFinish struct {
	IsFinished int
	EndTime    int
}
type TaskAllotMemberLists struct {
	SaName    string `json:"sa_name" desc:"触达导购"`
	SaCode    string `json:"sa_code" desc:"导购ID"`
	StoreName string `json:"store_name" desc:"所属门店"`
	StoreType string `json:"store_type" desc:"门店属性"`
}

type TaskDetailExportRes struct {
	Id             int64  `json:"id"`
	UserId         string `json:"user_id"`
	IsFinished     int    `json:"is_finished" desc:"子任务是否完成"`
	FinishedAt     int    `json:"finished_at" desc:"完成时间"`
	UnionID        string `json:"union_id" desc:"union_id"`
	CustomerID     string `json:"customer_id" desc:"customer_id"`
	IsMatch        int    `json:"is_match" desc:"是否匹配"`
	IsForward      int    `json:"is_forward" desc:"是否点击转发给顾客"`
	IsContact      int    `json:"is_contact" desc:"是否点击立即沟通"`
	IsCall         int    `json:"is_call" desc:"是否点击拨打电话"`
	Remark         string `json:"remark" desc:"联系说明"`
	IsAdd          int    `json:"is_add" desc:"是否点击添加好友"`
	CustomerStatus int    `json:"customer_status" desc:"顾客状态"`
	IsTouch        int    `json:"is_touch" desc:"触达状态，0-未发送 1-已发送 2-因客户不是好友导致发送失败 3-因客户已经收到其他群发消息导致发送失败"`
}

type TaskDetailExportListReq struct {
	BaseReq
	PlanId         int64   `json:"plan_id"`
	TaskAllotLogId []int64 `json:"task_allot_log_id"`
	Page           int     `json:"page" binding:"min=1" comment:"页码"`
	Size           int     `json:"size" binding:"min=50,max=99999999" comment:"分页数"`
	Total          int64   `json:"total"`
}

type TaskRateRes struct {
	UserId string  `json:"user_id"`
	Ratio  float64 `json:"ratio"`
}

type ConditionContent struct {
	Text struct {
		Content string `json:"content"`
	} `json:"text"`
	Image struct {
		ImgList []struct {
			PicUrl  string `json:"pic_url"`
			MediaId string `json:"media_id"`
		} `json:"imgList"`
	} `json:"image"`
	Voice struct {
		MediaId string `json:"media_id"`
		Format  string `json:"format"`
	} `json:"voice"`
	Video struct {
		MediaId      string `json:"media_id"`
		ThumbMediaId string `json:"thumb_media_id"`
		VideoUrl     string `json:"video_url"`
	} `json:"video"`
	Location struct {
		Location_X string `json:"location_x"`
		Location_Y string `json:"location_y"`
		Scale      string `json:"scale"`
		Label      string `json:"label"`
	} `json:"location"`
	Link struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Url         string `json:"url"` // 以link_url为主
		LinkUrl     string `json:"link_url"`
		PicUrl      string `json:"pic_url"`
		CoverUrl    string `json:"cover_url"` // PicUrl的地址
	} `json:"link"`
	Applets struct {
		Id            int64  `json:"id"`
		Type          int    `json:"type"`
		AdminId       int64  `json:"admin_id"`
		Title         string `json:"title"`
		MiniProgramId int64  `json:"mini_program_id"`
		AppId         string `json:"app_id"`
		AppName       string `json:"app_name"`
		Path          string `json:"path"`
		CoverUrl      string `json:"cover_url"`
	} `json:"applets"`
	Card struct {
		Id                 int      `json:"id"`
		TemplateName       string   `json:"template_name"`
		TemplatePictureURL []string `json:"template_picture_url"`
		TemplateContent    string   `json:"template_content"`
	} `json:"card"`
}

type NewStatReq struct {
	BaseReq
	StartTime int64 `json:"start_time"`
	EndTime   int64 `json:"end_time"`
	PlanId    int64 `json:"plan_id"`
}

type NewStatRes struct {
	PlanName                 string `json:"plan_name"`                  // 任务名称
	PeriodSubtasksNumber     int64  `json:"period_subtasks_number"`     // 子任务数
	SubtasksFinishNumber     int64  `json:"subtasks_finish_number"`     // 子任务完成数
	SubtasksFinishPercentage string `json:"subtasks_finish_percentage"` // 子任务完成率
	TaskNumber               int64  `json:"task_number"`                // 任务总数 (件)
	TaskFinishNumber         int64  `json:"task_finish_number"`         // 已完成任务数 (件)
	TaskFinishPercentage     string `json:"task_finish_percentage"`     // 任务完成率
	ReachGuideNumber         int64  `json:"reach_guide_number"`         // 触达导购总人数
	MatchCustomersNumber     int64  `json:"match_customers_number"`     // 匹配客户总人数 (不去重)
	MatchCustomersNumber2    int64  `json:"match_customers_number2"`    // 匹配客户总人数 (去重)
	ReachLordNumber          int64  `json:"reach_lord_number"`          // 触达群主总人数 (去重)
	MatchGroupNumber         int64  `json:"match_group_number"`         // 匹配社群总数
	PlanStartTime            int64  `json:"plan_start_time"`            // 任务开始时间
	PlanEndTime              int64  `json:"plan_end_time"`              // 任务结束时间
}

type NewStatDetailReq struct {
	BaseReq
	PlanId        int64  `json:"plan_id"`
	TaskName      string `json:"task_name"`
	Dimension     int    `json:"dimension"`       // 任务的维度，1 社群，2 用户画像 3 导入顾客名单 4 导入顾客&&saCode名单
	TaskCategory  int    `json:"task_category"`   // 任务类型 0 沟通任务 1 朋友圈任务
	GroupSendType int    `json:"group_send_type"` // 任务群发类型 0 标签任务群发 1 企业任务群发
	Page          int    `json:"page"`
	Size          int    `json:"size"`
}

type NewStatDetailRes struct {
	Count int64                  `json:"count"`
	Data  []*NewStatDetailResult `json:"data"`
}

type NewStatDetailResult struct {
	Id                    int64  `json:"id"`
	TaskName              string `json:"task_name"`               // 子任务名称
	TaskCategory          int    `json:"task_category"`           // 子任务类型
	GroupSendType         int    `json:"group_send_type"`         // 任务群发类型
	Dimension             int    `json:"dimension"`               // 条件维度
	ReachGuideNumber      int64  `json:"reach_guide_number"`      // 触达导购总人数
	TaskNumber            int64  `json:"task_number"`             // 任务总数 (件)
	TaskFinishNumber      int64  `json:"task_finish_number"`      // 已完成任务数 (件)
	TaskFinishPercentage  string `json:"task_finish_percentage"`  // 任务完成率
	MatchCustomersNumber  int64  `json:"match_customers_number"`  // 匹配客户总人数 (不去重)
	MatchCustomersNumber2 int64  `json:"match_customers_number2"` // 匹配客户总人数 (去重)
	ReachLordNumber       int64  `json:"reach_lord_number"`       // 触达群主总人数 (去重)
	MatchGroupNumber      int64  `json:"match_group_number"`      // 匹配社群总数
}

type TaskNoFinishedCount struct {
	Id  int64 `json:"id" desc:"task_allot_log_id"`
	Num int64 `json:"num"`
}

type MomentDelTagQueueData struct {
	EnterpriseId int64 `json:"enterprise_id"`
	TaskId       int64 `json:"task_id"`
}

type TaskNoTouchInfo struct {
	Id     int64  `json:"id" desc:"task_allot_log_id"`
	UserId string `json:"user_id"`
}

type TaskDetailExportListNewReq struct {
	BaseReq
	TaskAllotDetailExportBase
	PlanId         int64   `json:"plan_id"`
	TaskAllotLogId []int64 `json:"task_allot_log_id"`
	Page           int     `json:"page" binding:"min=1" comment:"页码"`
	Size           int     `json:"size" binding:"min=1,max=99999999" comment:"分页数"`
}

type TaskDetailExportNewRes struct {
	PlanName                string `json:"plan_name"`
	TaskName                string `json:"task_name"`
	CategoryName            string `json:"category_name"`
	TaskCategoryName        string `json:"task_category_name"`
	SaName                  string `json:"sa_name"`
	SaCode                  string `json:"sa_code"`
	StoreName               string `json:"store_name"`
	StoreType               string `json:"store_type"`
	IsTaskFinished          string `json:"is_task_finished" desc:"子任务是否完成"`
	FinishedTimeStr         string `json:"finished_time_str" desc:"完成时间"`
	UnionID                 string `json:"union_id" desc:"union_id"`
	CustomerID              string `json:"customer_id" desc:"customer_id"`
	IsMatch                 string `json:"is_match" desc:"是否匹配"`
	IsForward               string `json:"is_forward" desc:"是否点击转发给顾客"`
	IsContact               string `json:"is_contact" desc:"是否点击立即沟通"`
	IsCall                  string `json:"is_call" desc:"是否点击拨打电话"`
	Remark                  string `json:"remark" desc:"联系说明"`
	IsAdd                   string `json:"is_add" desc:"是否点击添加好友"`
	CustomerStatus          string `json:"customer_status" desc:"顾客状态"`
	IsTouch                 int    `json:"is_touch" desc:"触达状态，0-未发送 1-已发送 2-因客户不是好友导致发送失败 3-因客户已经收到其他群发消息导致发送失败"`
	TouchName               string `json:"touch_name" desc:"触达状态"`
	FinishRate              string `json:"finish_rate" desc:"触达率"`
	TouchType               string `json:"touch_type" desc:"触达类型，群发触达，单独沟通"`
	CustomerFinishedTimeStr string `json:"customer_finished_time_str" desc:"顾客联系时间"`
}

type MaPushFinishQueueData struct {
	EnterpriseId    int64    `json:"enterprise_id"`
	TaskId          int64    `json:"task_id"`
	PlanId          int64    `json:"plan_id"`
	ExternalUserIds []string `json:"external_user_ids"`
	UserId          string   `json:"user_id"`
	TaskAllotLogId  int64    `json:"task_allot_log_id"`
}

type ExternalSearchListsReq struct {
	BaseReq
	MemberId       int64                        `json:"member_id"`
	Id             int64                        `json:"id"`
	FilterType     int                          `json:"filterType"`
	Keywords       string                       `json:"keywords"`
	CustomerStatus int                          `json:"customer_status"`
	IsTag          int                          `json:"is_tag" comment:"0为不展示tag,1展示tag"`
	Search         customerSearchReq.SearchList `json:"search"`
	request.PageReq
}

type CustomerSearchListsReq struct {
	BaseReq
	Id             int64              `json:"id"`
	FilterType     int                `json:"filterType"`
	Keywords       string             `json:"keywords"`
	CustomerStatus int                `json:"customer_status"`
	IsTag          int                `json:"is_tag" comment:"0为不展示tag,1展示tag"`
	Search         request.SearchList `json:"search"`
	request.PageReq
}

type CustomerMoveReq struct {
	BaseReq
	Id             int64    `json:"id"`
	UnionIds       []string `json:"union_ids"`
	CustomerStatus int      `json:"customer_status" comment:"0为移除，1为移入"`
	OpenIds        []string `json:"open_ids"`
}

type SendFinishReq struct {
	BaseReq
	TaskAllotLogId int64    `json:"task_allot_log_id"`
	MemberId       int64    `json:"member_id"`
	Mobile         string   `json:"mobile"`
	UnionIds       []string `json:"union_ids"`
	OpenIds        []string `json:"open_ids"`
}

type TaskMatchOrder struct {
	TaskId          string  `json:"task_id"`
	PrimaryTaskId   int64   `json:"primary_task_id"`
	SubTaskId       int64   `json:"sub_task_id"`
	SaCode          string  `json:"sa_code"`
	UnionId         string  `json:"union_id"`
	CounterCode     string  `json:"counter_code"`
	PurchaseCode    string  `json:"purchase_code"`
	OriPurchaseCode string  `json:"ori_purchase_code"`
	PurchaseAmount  float64 `json:"purchase_amount"`
	PurchaseDate    string  `json:"purchase_date"`
	CreateDate      string  `json:"create_date"`
	MatchStatus     int64   `json:"match_status"`
}

type TaskMatchOrderProduct struct {
	PurchaseCode string  `json:"purchase_code"`
	SkuName      string  `json:"sku_name"`
	Style        string  `json:"style"`
	Color        string  `json:"color"`
	Amount       float64 `json:"amount"`
	PurchaseDate string  `json:"purchase_date"`
	CreateDate   string  `json:"create_date"`
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
	OriginType       int    `json:"origin_type" comment:"0:默认手动创建，1:convertlab创建，2:startDt创建"`
	request.PageReq
}

type CustomerTaskReq struct {
	request.EnterpriseBaseReq
	request.PageReq
	UserId       string `json:"user_id"`
	UnionId      string `json:"union_id"`
	Status       int    `json:"status"`
	TaskCategory int    `json:"task_category"`
}
