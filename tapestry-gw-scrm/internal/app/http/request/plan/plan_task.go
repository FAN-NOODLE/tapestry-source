package plan

type BaseReq struct {
	EnterpriseId int64 `json:"enterprise_id"`
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
	EndTime             int64  `json:"end_time"`
	TaskType            int    `json:"task_type"`
}

type TaskAllotDetailExportBase struct {
	PlanId          int64  `json:"plan_id"`
	PlanName        string `json:"plan_name" desc:"任务计划名称"`
	TaskName        string `json:"task_name" desc:"子任务名称"`
	Category        int    `json:"category" desc:"任务的创建类型"`
	MemberLists     []TaskAllotMemberLists
	TaskAllotLogIds []int64
	TaskForwards    map[int64]int
	TaskFinish      map[int64]TaskAllotFinish
	TaskCategory    int `json:"task_category" desc:"子任务类型，沟通任务，朋友圈任务"`
	GroupSendType   int `json:"group_send_type" desc:"发送类型,标签任务群发/企业任务群发"`
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

type TaskDetailExportLists struct {
	Total int64
	List  []TaskDetailExportRes
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

// NewStatReq 任务统计接收参数
type NewStatReq struct {
	BaseReq
	StartTime int64 `json:"start_time"`
	EndTime   int64 `json:"end_time"`
	PlanId    int64 `json:"plan_id"`
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
	Count int                    `json:"count"`
	Data  []*NewStatDetailResult `json:"data"`
}

type NewStatDetailResult struct {
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

type TaskDetailExportListNewReq struct {
	BaseReq
	TaskAllotDetailExportBase
	TaskAllotLogId []int64 `json:"task_allot_log_id"`
	Page           int     `json:"page" binding:"min=1" comment:"页码"`
	Size           int     `json:"size" binding:"min=1,max=99999999" comment:"分页数"`
}
type TaskPlanNameListReq struct {
	BaseReq
	PlanId int64 `json:"plan_id"`
}
