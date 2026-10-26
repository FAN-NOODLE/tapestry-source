package microScrm

// FissionActivityListResponse 裂变活动返回结构体
type FissionActivityListResponse struct {
	List       []FissionActivityList `json:"list"`
	Statistics Statistics            `json:"statistics"`
	Page       int64                 `json:"page"`
	Size       int64                 `json:"size"`
}

// Statistics 统计数据
type Statistics struct {
	UserNum       int64 `json:"user_num" desc:"发起人数"`
	FissionNum    int64 `json:"fission_num" desc:"好友裂变数"`
	RunAwayNum    int64 `json:"run_away_num" desc:"流失人数"`
	JoinGroupNum  int64 `json:"join_group_num" desc:"群成员裂变数"`
	LeaveGroupNum int64 `json:"leave_group_num" desc:"群成员退群数"`
}

// FissionActivityList 列表数据
type FissionActivityList struct {
	Id      int64  `json:"id"`
	Name    string `json:"name"`
	Desc    string `json:"desc"`
	Status  int    `json:"status"`
	HtmlUrl string `json:"html_url"`
	ImgUrl  string `json:"img_url"`
}

// FissionActivityListRequest 裂变活动请求结构体
type FissionActivityListRequest struct {
	EnterpriseId   int64  `json:"enterprise_id" db:"enterprise_id" desc:"企业Id"`
	EnterpriseHash string `json:"enterprise_hash" desc:"企业hash"`
	SearchText     string `json:"search_text" desc:"搜索条件"`
	Status         int64  `json:"status" desc:"活动状态：1-上架中，2-已下架"`
	Page           int64  `json:"page"`
	Size           int64  `json:"size"`
}
