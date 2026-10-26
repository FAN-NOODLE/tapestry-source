package dashboardReq

import "micro-scrm/internal/app/http/request"

type AreaSave struct {
	request.EnterpriseBaseReq
	Id        int    `json:"id"`
	Name      string `json:"name"`
	ParentId  int    `json:"parent_id"`
	ThirdId   int    `json:"third_id"`
	BrotherId int    `json:"brother_id"`
	IsAbove   bool   `json:"is_above"`
}

type AreaTree struct {
	request.EnterpriseBaseReq
	Language         string `json:"language"`
	ParentId         int    `json:"parent_id"`
	HasChildren      bool   `json:"has_children"`
	HasVirtual       bool   `json:"has_virtual"`
	IgnoreCloseStore bool   `json:"ignore_close_store"`
}

type StoreList struct {
	request.EnterpriseBaseReq
	SearchStr string `json:"search_str"`
	Page      int    `json:"page"`
	Size      int    `json:"size"`
}

type AreaRemove struct {
	request.EnterpriseBaseReq
	Id int `json:"id"`
}

type ManagerRemove struct {
	request.EnterpriseBaseReq
	Id int `json:"id"`
}

type ManagerAuthSave struct {
	request.EnterpriseBaseReq
	ID            int64   `json:"id"`
	MemberId      int64   `json:"member_id" binding:"required" comment:"导购member_id"`
	Position      string  `json:"position"  binding:"required" comment:"导购职位"`
	PositionLevel int64   `json:"position_level"`
	AreaIds       []int64 `json:"area_ids"  binding:"required" comment:"导购管理范围"`
}

type ManagerAuthList struct {
	request.EnterpriseBaseReq
	UserId   string `json:"user_id"`
	UserName string `json:"user_name"`
	Page     int    `json:"page"`
	Size     int    `json:"size"`
}

type ManagerAuthDetail struct {
	request.EnterpriseBaseReq
	ID int64 `json:"id" binding:"required" comment:"职位id"`
}

type ManagerAuthAreaIds struct {
	request.EnterpriseBaseReq
	MemberId     int64   `json:"member_id"`
	AreaIds      []int64 `json:"area_ids"`
	TargetAreaId int64   `json:"target_area_id"`
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
	MemberId      int64   `json:"member_id" binding:"required" comment:"导购member_id"`
	Position      string  `json:"position" binding:"required" comment:"导购职位"`
	PositionLevel int64   `json:"position_level"`
	AreaIds       []int64 `json:"area_ids" binding:"required" comment:"导购管理范围"`
}

type HkManagerAuthList struct {
	request.EnterpriseBaseReq
	UserId   string `json:"user_id"`
	UserName string `json:"user_name"`
	Page     int    `json:"page"`
	Size     int    `json:"size"`
}

type HkManagerAuthDetail struct {
	request.EnterpriseBaseReq
	ID int64 `json:"id" binding:"required" comment:"职位id"`
}

type HkManagerRemove struct {
	request.EnterpriseBaseReq
	Id int `json:"id"`
}

type HkManagerAuthAreaIds struct {
	request.EnterpriseBaseReq
	MemberId int64   `json:"member_id"`
	AreaIds  []int64 `json:"area_ids"`
}

type StoreCloseOrOpen struct {
	request.EnterpriseBaseReq
	ThirdId  int  `json:"third_id" binding:"required" comment:"企微门店id"`
	IsClosed bool `json:"is_closed"`
}

type SearchStoreReq struct {
	request.EnterpriseBaseReq
	Language         string `json:"language"`
	MemberId         int64  `json:"member_id"`
	StoreStr         string `json:"store_str"`
	IgnoreCloseStore bool   `json:"ignore_close_store"`
}
