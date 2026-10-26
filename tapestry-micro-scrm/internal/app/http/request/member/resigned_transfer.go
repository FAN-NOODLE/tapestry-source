package memberReq

type ResignedTransferCustomerReq struct {
	ExternalUserid []string `json:"external_userid" binding:"required" comment:"企业ID"` // 外部联系人userid(openid)
	EnterpriseID   int64    `json:"enterprise_id" binding:"required" comment:"企业ID"`
	HandoverUserid string   `json:"handover_userid" binding:"required" comment:"原跟进人"` // 原跟进成员的userid
	TakeoverUserid string   `json:"takeover_userid" binding:"required" comment:"接替人"`  // 接替人userid
}

type ResignedTransferCustomerGroupReq struct {
	EnterpriseID   int64    `json:"enterprise_id" binding:"required" comment:"企业ID"`
	ChatIdList     []string `json:"chat_id_list" binding:"gt=0" comment:"客户群ID"`       // 需要转群主的客户群ID列表。取值范围： 1 ~ 100
	HandoverUserid string   `json:"handover_userid" binding:"required" comment:"原跟进人"` // 原跟进成员的userid
	TakeoverUserid string   `json:"takeover_userid" binding:"required" comment:"接替人"`  // 接替人userid
}

type EUidList struct {
	Euid int64 `json:"euid" db:"euid"`
}
