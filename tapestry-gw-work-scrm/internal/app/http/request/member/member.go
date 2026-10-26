package member

import "gw-work-scrm/internal/app/http/request"

type GuidePositionReq struct {
	request.EnterpriseBaseReq
	AdminId    int64  `json:"admin_id"`
	Mobile     string `json:"mobile"`
	IsNew      int    `json:"is_new"`
	MemberType int    `json:"member_type"  comment:"当前导购类型"`
}

type GuidePositionRes struct {
	Role int `json:"role"`
}

type FirstOpenReq struct {
	request.EnterpriseBaseReq
	Type   int    `json:"type"`
	Mobile string `json:"mobile"`
}

type GuideListReq struct {
	StoreOrDepartmentsReq
	Domain string `json:"domain"`
}

type StoreOrDepartmentsReq struct {
	request.EnterpriseBaseReq
	AdminId int64  `json:"admin_id"`
	Mobile  string `json:"mobile"`
	DpIds   []int  `json:"dp_ids"`
}
