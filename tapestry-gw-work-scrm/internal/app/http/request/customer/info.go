package customerReq

type SaveInfoReq struct {
	EnterpriseID     int64  `json:"enterprise_id" comment:"企业ID"`
	EnterpriseUserID int64  `json:"enterprise_user_id" comment:"客户ID"`
	Address          string `json:"address" comment:"客户地址"`
	ShoeSize         string `json:"shoe_size" comment:"鞋码"`
	Month            int    `json:"month" comment:"客户生日（月）"`
	Day              int    `json:"day" comment:"客户生日（日）"`
	OperateUserId    string `json:"operate_user_id" comment:"操作人Id(企业微信userid)"`
	OperateName      string `json:"operate_name" comment:"操作人姓名"`
	ExternalUserID   string `json:"external_userid" db:"external_userid"`
}

type SimpleInfoRequest struct {
	EnterpriseID   int64  `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	ExternalUserID string `json:"external_userid" db:"external_userid"`
	MemberId       int64  `json:"member_id"`
}
