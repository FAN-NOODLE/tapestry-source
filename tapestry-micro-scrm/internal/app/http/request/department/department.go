package department

type BaseReq struct {
	EnterpriseID int64 `json:"enterprise_id" comment:"企业ID"`
}

type SAStoreListReq struct {
	BaseReq
	StoreType int `json:"store_type"`
}
type DepartmentsLists struct {
	OCList  []NewMomentDepartments `json:"oc_list" desc:"零售"`
	OCFList []NewMomentDepartments `json:"ocf_list" desc:"奥莱"`
}

type MomentDepartments struct {
	ID                int64  `json:"dp_id" db:"id" desc:"部门ID"`
	DpName            string `json:"dp_name" db:"name" desc:"部门名称"`
	ThirdID           int64  `json:"third_id" db:"third_id" desc:"企业微信中组织ID"`
	ThirdParentID     int64  `json:"third_parent_id" db:"third_parent_id" desc:"企业微信中的父ID"`
	Order             int32  `json:"dp_display_order" db:"order" desc:"排序"`
	IsChildDepartment int    `json:"is_child_department" db:"is_child_department" desc:"是不是子部门"`
	StoreType         int    `json:"store_type" db:"store_type" desc:"部门类型,1为零售,2为奥莱,0为默认"`
}

type NewMomentDepartments struct {
	ID                 int64  `json:"id" db:"id" desc:"部门ID"`
	Name               string `json:"name" db:"name" desc:"部门名称"`
	ThirdID            int64  `json:"third_id" db:"third_id" desc:"企业微信中组织ID"`
	IsDp               bool   `json:"is_dp" `
	Leaf               int32  `json:"leaf"`
	MemSubscribeStatus int    `json:"mem_subscribe_status"`
	StoreType          int    `json:"store_type" db:"store_type" desc:"部门类型,1为零售,2为奥莱,0为默认"`
}
