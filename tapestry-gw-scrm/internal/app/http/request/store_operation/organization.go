package store_operation

type OrganizationTreeReq struct {
	EnterpriseID int64 `json:"enterprise_id"`
}

type OrganizationStoreListReq struct {
	EnterpriseID int64  `json:"enterprise_id"`
	NodeID       int64  `json:"node_id"`
	StoreCode    string `json:"store_code"`
	StoreName    string `json:"store_name"`
	Page         int    `json:"page"`
	Size         int    `json:"size"`
}
