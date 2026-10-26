package openReq

// ConvertOpenUserIdRequest 将open_userid转为userid请求
type ConvertOpenUserIdRequest struct {
	EnterpriseID   int64    `json:"enterprise_id"`                       // 企业ID
	OpenUseridList []string `json:"open_userid_list" binding:"required"` // 代开发应用的open_userid
	SourceAgentId  int      `json:"source_agentid" binding:"required"`   //代开发应用的agentid
}

// ConvertOpenExternalUserIdRequest 将open_external_userid转为external_userid请求
type ConvertOpenExternalUserIdRequest struct {
	EnterpriseID   int64  `json:"enterprise_id"`                      // 企业ID
	ExternalUserid string `json:"external_userid" binding:"required"` // 代开发应用的外部联系人id
	SourceAgentId  int    `json:"source_agentid" binding:"required"`  // 代开发应用的agentid
}
