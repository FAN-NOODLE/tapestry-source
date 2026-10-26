package openReq

// StoreMasterListReq 门店主数据分页查询参数。
type StoreMasterListReq struct {
	EnterpriseID      int64    `json:"-"`
	Page              int      `json:"page"`
	Size              int      `json:"size"`
	UpdatedAt         int64    `json:"updated_at"`
	StoreCountryCodes []string `json:"store_country_codes"`
}

// EmployeeMasterListReq 员工主数据分页查询参数。
type EmployeeMasterListReq struct {
	EnterpriseID      int64    `json:"-"`
	Page              int      `json:"page"`
	Size              int      `json:"size"`
	UpdatedAt         int64    `json:"updated_at"`
	StoreScope        string   `json:"store_scope"`
	StoreCountryCodes []string `json:"store_country_codes"`
	EmploymentScope   string   `json:"employment_scope"`
}
