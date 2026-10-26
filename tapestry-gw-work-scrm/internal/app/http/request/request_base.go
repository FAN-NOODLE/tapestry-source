package request

type EnterpriseBaseReq struct {
	EnterpriseId int64 `json:"enterprise_id" comment:"企业ID"`
}

type PageReq struct {
	Page int `json:"page" binding:"min=1" comment:"页码"`
	Size int `json:"size" binding:"min=1,max=1000" comment:"分页数"`
}

type MemberInfoReq struct {
	MemberName   string `json:"member_name"`
	MemberUserid string `json:"member_userid"`
	MemberID     int    `json:"member_third_id"`
}
