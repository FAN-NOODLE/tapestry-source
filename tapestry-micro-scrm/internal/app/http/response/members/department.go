package membersResp

import departmentsModels "micro-scrm/internal/app/models/departments"

type DepListResp struct {
	departmentsModels.DepartmentsModel
	StoreName string `json:"store_name"`
	StoreCode string `json:"store_code"`
}

type SaListResp struct {
	departmentsModels.DepartmentsMembersData
	StoreName string `json:"store_name"`
	StoreCode string `json:"store_code"`
}
