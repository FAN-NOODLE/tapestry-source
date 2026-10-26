package micro_scrm_enterprise

import "micro-scrm/internal/orm/model"

// EmployeeMasterModel Coach Workday 员工主数据。
type EmployeeMasterModel struct {
	model.BranchModel
	EmployeeID        string `gorm:"column:employee_id" json:"employee_id"`
	EmployeeCode      string `gorm:"column:employee_code" json:"employee_code"`
	MobileNumber      string `gorm:"column:mobile_number" json:"mobile_number"`
	Name              string `gorm:"column:name" json:"name"`
	AliasName         string `gorm:"column:alias_name" json:"alias_name"`
	NameZh            string `gorm:"column:name_zh" json:"name_zh"`
	JobLocation       string `gorm:"column:job_location" json:"job_location"`
	JobPosition       string `gorm:"column:job_position" json:"job_position"`
	Gender            string `gorm:"column:gender" json:"gender"`
	EmailAddress      string `gorm:"column:email_address" json:"email_address"`
	Status            string `gorm:"column:status" json:"status"`
	StoreCode         string `gorm:"column:store_code" json:"store_code"`
	Department        string `gorm:"column:department" json:"department"`
	HireDate          string `gorm:"column:hire_date" json:"hire_date"`
	TerminationDate   string `gorm:"column:termination_date" json:"termination_date"`
	JobFunction       string `gorm:"column:job_function" json:"job_function"`
	IsStoreManagement string `gorm:"column:is_store_management" json:"is_store_management"`
}

func (EmployeeMasterModel) TableName() string { return "employee_master" }
