package http_entity

type PageParam struct {
	Page int `json:"page,omitempty"`
	Size int `json:"size,omitempty"`
}

// 子任务导出任务详情请求
type ExportTasksRequst struct {
	Id                  int64   `json:"id"  binding:"required" desc:"条件编号"`
	PlanStartTime       int64   `json:"plan_start_time" desc:"任务开始时间"`
	PlanEndTime         int64   `json:"plan_end_time" desc:"任务结束时间"`
	DepartmentId        int64   `json:"department_id" desc:"部门编号"`
	DepartmentIds       []int64 `json:"department_ids" desc:"部门集合"`
	IsFinished          int     `json:"is_finished" desc:"完成状态；1是、2否"`
	BaName              string  `json:"ba_name" desc:"导购名称"`
	CompletionStartTime int64   `json:"completion_start_time" desc:"任务完成时间开始时间"`
	CompletionEndTime   int64   `json:"completion_end_time" desc:"任务完成时间结束时间"`
	EndTime             int64   `json:"end_time" desc:"订单统计结束时间"`
}

// 子任务列表数据
type TaskListRequst struct {
	ExportTasksRequst
	PageParam
}
