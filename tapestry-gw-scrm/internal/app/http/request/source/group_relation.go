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
	EnterpriseID  int64 `json:"enterprise_id"`
	SourceGroupId int64 `json:"source_group_id"`
	DepartmentId  int64 `json:"department_id"`
	ThirdId       int   `json:"third_id"`
}
