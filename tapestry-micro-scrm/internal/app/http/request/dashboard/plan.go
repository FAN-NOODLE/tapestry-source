package dashboardReq

type EnterpriseBase struct {
	EnterpriseId int64 `json:"enterprise_id" binding:"required" required:"企业ID必填"`
}

type SearchTime struct {
	StartTime int64 `json:"start_time" desc:"开始时间"`
	EndTime   int64 `json:"end_time" desc:"结束时间"`
}

type PlanStatistics struct {
	EnterpriseBase
	SearchTime
	MemberName   string `json:"member_name"`
	MemberUserid string `json:"member_userid"`
	StoreIds     []int  `json:"store_ids" common:"企微部门id(查询的范围)"`
}
