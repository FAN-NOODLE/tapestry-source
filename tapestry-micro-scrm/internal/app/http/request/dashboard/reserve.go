package dashboardReq

type DataKanban struct {
	EnterpriseId  int64  `json:"enterprise_id" binding:"required" required:"企业ID必填"`
	MemberName    string `json:"member_name"`
	MemberUserid  string `json:"member_userid"`
	DepartmentIds []int  `json:"department_ids" common:"部门id集合(查询的范围)"`

	StartTime int   `json:"start_time"`
	EndTime   int   `json:"end_time"`
	StoreIds  []int `json:"store_ids"`
	MemberIds []int `json:"member_ids"`
}

type DataList struct {
	EnterpriseId int64             `json:"enterprise_id" binding:"required" comment:"企业ID"`
	NameType     int               `json:"name_type" binding:"required" comment:"名称类型;1=导购,2=渠道,3=大区,4=分区,5=门店"`
	OrderMap     map[string]string `json:"order_map" comment:"排序"`
	Page         int               `json:"page" binding:"min=1" min:"页码最小为1"`
	Size         int               `json:"size" binding:"min=10,max=1000" min:"分页数最小为1" max:"分页数最大为1000"`

	StartTime int   `json:"start_time"`
	EndTime   int   `json:"end_time"`
	StoreIds  []int `json:"store_ids"`
	MemberIds []int `json:"member_ids"`
}
