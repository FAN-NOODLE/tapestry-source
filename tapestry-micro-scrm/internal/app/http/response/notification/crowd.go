package notificationResp

type CrowdDataResp struct {
	MemberId     int64  `json:"member_id"`
	MemberName   string `json:"member_name"`
	DepartmentId int64  `json:"department_id"`
}

type CrowdDetailResp struct {
	Id         int64           `json:"id"`
	Name       string          `json:"name"`
	Type       int             `json:"type"`
	MemberList []CrowdDataResp `json:"member_list"`
}

type ImportDataResp struct {
	StoreCodes []string        `json:"store_codes"`
	MemberList []CrowdDataResp `json:"member_list"`
}
