package dashboardReq

import "micro-scrm/internal/app/http/request"

type HkStoreSync struct {
	request.EnterpriseBaseReq
}

type HkStoreList struct {
	request.EnterpriseBaseReq
	SearchStr          string `json:"search_str"`
	Area               string `json:"area"`
	StoreType          int    `json:"store_type"`
	WechatDepartmentID int    `json:"wechat_department_id"`
	ClosedStatus       int    `json:"closed_status" comment:"关店状态：0全部 1营业中 2已关店"`
	Page               int    `json:"page"`
	Size               int    `json:"size"`
}

type HkStoreDetail struct {
	request.EnterpriseBaseReq
	ID int64 `json:"id" binding:"required"`
}

type HkStoreSave struct {
	request.EnterpriseBaseReq
	ID            int64  `json:"id" binding:"required"`
	StoreName     string `json:"store_name"`
	Area          string `json:"area"`
	StoreType     int    `json:"store_type"`
	InventoryName string `json:"inventory_name"`
	Address       string `json:"address"`
	City          string `json:"city"`
	Phone         string `json:"phone"`
	Email         string `json:"email"`
}

type HkStoreClose struct {
	request.EnterpriseBaseReq
	ID int64 `json:"id" binding:"required"`
}

type HkStoreLastSyncTime struct {
	request.EnterpriseBaseReq
}
