package openReq

type CustomerClearReq struct {
	EnterpriseId int64  `json:"enterprise_id"`
	UnionId      string `json:"unionid" binding:"required" required:"unionid必填"`
}
