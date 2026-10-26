package dashboardReq

import "gw-scrm/internal/app/http/request"

type HkStoreSync struct {
	request.EnterpriseBaseReq
}

type HkStoreList struct {
	request.EnterpriseBaseReq
	SearchStr          string `json:"search_str"`
	Area               string `json:"area"`
	StoreType          int    `json:"store_type"`
	WechatDepartmentID int    `json:"wechat_department_id"`
	ClosedStatus       int    `json:"closed_status"`
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

type HkManagerStoreList struct {
	request.EnterpriseBaseReq
	StoreName          string `json:"store_name"`
	StoreCode          string `json:"store_code"`
	Area               string `json:"area"`
	WechatDepartmentID int    `json:"wechat_department_id"`
	ExcludeClosed      bool   `json:"exclude_closed"`
}

type HkManagerAuthSave struct {
	request.EnterpriseBaseReq
	ID            int64   `json:"id"`
	MemberId      int64   `json:"member_id" binding:"required"`
	Position      string  `json:"position" binding:"required"`
	PositionLevel int64   `json:"position_level"`
	AreaIds       []int64 `json:"area_ids" binding:"required"`
}

type HkManagerAuthList struct {
	request.EnterpriseBaseReq
	UserId   string `json:"user_id"`
	UserName string `json:"user_name"`
	Position string `json:"position"`
	Page     int    `json:"page"`
	Size     int    `json:"size"`
}

type HkManagerAuthDetail struct {
	request.EnterpriseBaseReq
	ID int64 `json:"id" binding:"required"`
}

type HkManagerRemove struct {
	request.EnterpriseBaseReq
	Id int `json:"id"`
}
