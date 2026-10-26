package source

type GroupRelationListRequest struct {
	EnterpriseID  int64 `json:"enterprise_id"`
	SourceGroupId int64 `json:"source_group_id"`
	Page          int   `json:"page"`
	Size          int   `json:"size"`
}

type GroupRelationSaveRequest struct {
	EnterpriseID   int64   `json:"enterprise_id"`
	SourceGroupIds []int64 `json:"source_group_ids"`
	DepartmentIds  []int64 `json:"department_ids"`
	ThirdIds       []int   `json:"third_ids"`
}

type GroupRelationFindRequest struct {
	EnterpriseID  int64  `json:"enterprise_id"`
	SourceGroupId int64  `json:"source_group_id"`
	MemberID      int64  `json:"member_id" comment:"员工Id"`
	MemberUserId  string `json:"member_user_id" comment:"企业员工Id"`
}

// RequestSelectList 各素材列表请求结构
type RequestSelectList struct {
	EnterpriseID  int64   `json:"enterprise_id" db:"enterprise_id" validate:"required" required:"企业ID必传"`
	SourceGroupID int64   `json:"source_group_id" db:"source_group_id"`
	Keyword       string  `json:"keyword" db:"keyword"`
	CreatedStart  int64   `json:"create_start" db:"create_start"`
	CreatedEnd    int64   `json:"create_end" db:"create_end"`
	Size          int     `json:"size" db:"size"`
	Page          int     `json:"page" db:"page"`
	CategoryIds   []int64 `json:"category_ids"`
	DepartmentId  int64   `json:"department_id"`
	Mobile        string  `json:"mobile" db:"mobile"`
	MemberID      int64   `json:"member_id" comment:"员工Id"`
	MemberUserId  string  `json:"member_user_id" comment:"企业员工Id"`
	Type          int     `json:"type" comment:"类型(1 文本 2 图片 3 链接 4 小程序 5 视频 6 文件 7 朋友圈)"`
	CollectSearch bool    `json:"collect_search" comment:"是否收藏搜索"`
}

type RequestNewSource struct {
	EnterpriseID int64  `json:"enterprise_id" validate:"required" required:"企业ID必传"`
	MemberID     int64  `json:"member_id" comment:"员工Id"`
	MemberUserId string `json:"member_user_id" comment:"企业员工Id"`
	Mobile       string `json:"mobile" `
}
