package customerReq

type PageParam struct {
	Page int `json:"page" comment:"页码"`
	Size int `json:"size" comment:"分页数"`
}

type CustomerBaseReq struct {
	EnterpriseID int64 `json:"enterprise_id" comment:"企业ID"`
}

type CustomerListReq struct {
	CustomerBaseReq
	PageParam
	HandoverUserId       string `json:"handover_userid" comment:"离职成员的userid"`
	HandoverName         string `json:"handover_name" comment:"离职成员姓名"`
	HandoverDepartmentId int64  `json:"handover_department_id" comment:"离职成员所属门店"`
	Status               int    `json:"status" comment:"状态: 0 未分配 1已分配"`
}

type ResignedTransferCustomerReq struct {
	ExternalUserid []string `json:"external_userid" comment:"企业ID"` // 外部联系人userid(openid)
	EnterpriseID   int64    `json:"enterprise_id" comment:"企业ID"`
	HandoverUserid string   `json:"handover_userid" comment:"原跟进人"` // 原跟进成员的userid
	TakeoverUserid string   `json:"takeover_userid" comment:"接替人"`  // 接替人userid
}

type ResignedTransferCustomerGroupReq struct {
	EnterpriseID   int64    `json:"enterprise_id" comment:"企业ID"`
	ChatIdList     []string `json:"chat_id_list" comment:"客户群ID"`   // 需要转群主的客户群ID列表。取值范围： 1 ~ 100
	HandoverUserid string   `json:"handover_userid" comment:"原跟进人"` // 原跟进成员的userid
	TakeoverUserid string   `json:"takeover_userid" comment:"接替人"`  // 接替人userid
}

type RequestCustomerExternalList struct {
	SearchName     string `json:"search_name"`
	HandoverUserid string `json:"handover_userid"`
	EnterpriseID   int64  `json:"enterprise_id"`
	Page           int    `json:"page"`
	Size           int    `json:"size"`
}

type RequestCustomerGroupList struct {
	HandoverUserid string `json:"handover_userid"`
	EnterpriseID   int64  `json:"enterprise_id"`
	Page           int    `json:"page"`
	Size           int    `json:"size"`
}
type RequestTakeOverExternalList struct {
	TakeoverUserId string `json:"takeover_userid"`
	RequestCustomerExternalList
}

type SyncCustomDataRes struct {
	CustomerBaseReq
	ExeSync       bool `json:"exe_sync"`
}
