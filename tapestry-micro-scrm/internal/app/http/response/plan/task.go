package planResp

type CustomerTaskResp struct {
	Id            int64  `json:"id"`
	TaskId        int64  `json:"task_id"`
	Status        int    `json:"status"`
	TaskCategory  int    `json:"task_category"`
	Name          string `json:"name"`
	Remark        string `json:"remark"`
	TaskStartTime int    `json:"task_start_time"`
	TaskEndTime   int    `json:"task_end_time"`
}
