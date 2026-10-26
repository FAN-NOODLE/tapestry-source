package dashboardReq

import "gw-scrm/internal/app/http/request"

//type StoreList struct {
//	SearchStr string `json:"search_str"`
//	request.EnterpriseBaseReq
//}

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
	ID            int64  `json:"id"`
	MemberId      int64  `json:"member_id"`
	Position      string `json:"position"`
	PositionLevel int64  `json:"position_level"`
	AreaIds       []int  `json:"area_ids"`
}

type ManagerAuthDetail struct {
	request.EnterpriseBaseReq
	ID int64 `json:"id"`
}

type ManagerAuthList struct {
	request.EnterpriseBaseReq
	UserId   string `json:"user_id"`
	UserName string `json:"user_name"`
	Position string `json:"position"`
	Page     int    `json:"page"`
	Size     int    `json:"size"`
}
