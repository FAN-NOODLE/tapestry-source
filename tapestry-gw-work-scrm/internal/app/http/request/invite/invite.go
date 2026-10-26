package invite

type EnterpriseBase struct {
	EnterpriseId int64 `json:"enterprise_id" binding:"required" comment:"企业ID"`
}

type TemplateConfigListReq struct {
	EnterpriseBase
	Mobile   string `json:"mobile" comment:"手机号"`
	ThirdId  int64  `json:"third_ids"`
	MemberId int64  `json:"member_id"`
}

type TemplateConfigDetailReq struct {
	EnterpriseBase
	Mobile string `json:"mobile" comment:"手机号"`
	ID     int64  `json:"id"`
}

type RecordInfoReq struct {
	EnterpriseBase
	RecordInfo
}

type RecordInfo struct {
	InviteTemplateId int64  `json:"invite_template_id" binding:"required"`
	MemberId         int64  `json:"member_id" `
	MemberUserid     string `json:"member_userid" `
	UserUnionid      string `json:"user_unionid" `
	UserOpenid       string `json:"user_openid" binding:"required"`
	Mobile           string `json:"mobile"`
	StyleType        int    `json:"style_type"  binding:"required"`
	BackgroundUrl    string `json:"background_url"  binding:"required"`
	SendContent      string `json:"send_content"  binding:"required"`
}

type GetRecordInfoReq struct {
	EnterpriseBase
	InviteRecordId int64  `json:"invite_record_id" binding:"required"`
	UnionId        string `json:"union_id"`
}
