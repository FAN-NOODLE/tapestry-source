package dashboardResp

import (
	membersModel "micro-scrm/internal/app/models/members"
)

type ManagerAuthDetail struct {
	ID         int64                    `json:"id"`
	Position   string                   `json:"position"`
	AreaIds    []int64                  `json:"area_ids"`
	MemberInfo membersModel.MemberModel `json:"member_info"`
	StoreNum   int                      `json:"store_num"`
}

type ManagerList struct {
	ID          int      `json:"id"`
	Position    string   `json:"position"`
	Name        string   `json:"name"`
	Userid      string   `json:"userid"`
	PathNameArr []string `json:"path_name_arr"`
	StoreNum    int      `json:"store_num"`
	AreaIds     string   `json:"area_ids"`
}

type HkStoreInfo struct {
	ID                 int64  `json:"id"`
	StoreCode          string `json:"store_code"`
	StoreName          string `json:"store_name"`
	Area               string `json:"area"`
	StoreType          int32  `json:"store_type"`
	InventoryName      string `json:"inventory_name"`
	WechatDepartmentID uint   `json:"wechat_department_id"`
	ClosedAt           int64  `json:"closed_at"`
}

type HkManagerMemberInfo struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Userid string `json:"userid"`
	Mobile string `json:"mobile"`
}

type HkManagerAuthDetail struct {
	ID         int64               `json:"id"`
	Position   string              `json:"position"`
	AreaIds    []int64             `json:"area_ids"`
	MemberInfo HkManagerMemberInfo `json:"member_info"`
	StoreNum   int                 `json:"store_num"`
	StoreList  []HkStoreInfo       `json:"store_list"`
}

type HkManagerList struct {
	ID           int      `json:"id"`
	Position     string   `json:"position"`
	Name         string   `json:"name"`
	Userid       string   `json:"userid"`
	PathNameArr  []string `json:"path_name_arr"`
	StoreNameArr []string `json:"store_name_arr"`
	StoreNum     int      `json:"store_num"`
	AreaIds      string   `json:"area_ids"`
}
