package openResp

import pkgKMS "micro-scrm/internal/pkg/kms"

type StoreBusinessOrganization struct {
	RegionCode   string `json:"region_code"`
	Region       string `json:"region"`
	DistrictCode string `json:"district_code"`
	District     string `json:"district"`
	AreaCode     string `json:"area_code"`
	Area         string `json:"area"`
}

// StoreMasterResp 门店主数据对外响应。
type StoreMasterResp struct {
	StoreCode               string                    `json:"store_code"`
	CountryCode             string                    `json:"country_code"`
	Channel                 string                    `json:"channel"`
	StoreNameZh             string                    `json:"store_name_zh"`
	StoreNameEn             string                    `json:"store_name_en"`
	Status                  string                    `json:"status"`
	OpenDate                string                    `json:"open_date"`
	CloseDate               string                    `json:"close_date"`
	RegionDesc              string                    `json:"region_desc"`
	State                   string                    `json:"state"`
	City                    string                    `json:"city"`
	AddressZh               string                    `json:"address_zh"`
	Address                 string                    `json:"address"`
	Address2                string                    `json:"address2"`
	ZipCode                 string                    `json:"zip_code"`
	Phone                   string                    `json:"phone"`
	Email                   string                    `json:"email"`
	Longitude               string                    `json:"longitude"`
	Latitude                string                    `json:"latitude"`
	StoreWorkingHours       string                    `json:"store_working_hours"`
	SquareFootage           string                    `json:"square_footage"`
	AppointmentServiceTypes string                    `json:"appointment_service_types"`
	BrandCode               string                    `json:"brand_code"`
	OrganizationCode        string                    `json:"organization_code"`
	PaymentMethods          string                    `json:"payment_methods"`
	CompanyCode             string                    `json:"company_code"`
	StoreClass              string                    `json:"store_class"`
	WechatDepartmentID      int64                     `json:"wechat_department_id"`
	BusinessOrganization    StoreBusinessOrganization `json:"business_organization"`
	AreaManager             pkgKMS.EncryptedValue     `json:"area_manager"`
	StoreDistrictManager    pkgKMS.EncryptedValue     `json:"store_district_mgr"`
	StoreDesignConcept      string                    `json:"store_design_concept"`
	DPAP                    string                    `json:"dp_ap"`
	UpdatedAt               int64                     `json:"updated_at"`
}

type DepartmentChange struct {
	EventType          string `json:"event_type"`
	EventTime          int64  `json:"event_time"`
	WechatDepartmentID int64  `json:"wechat_department_id"`
	DepartmentName     string `json:"department_name"`
}

// EmployeeMasterResp 员工主数据对外响应。
type EmployeeMasterResp struct {
	Name               pkgKMS.EncryptedValue `json:"name"`
	NameZh             pkgKMS.EncryptedValue `json:"name_zh"`
	Gender             string                `json:"gender"`
	EmployeeID         pkgKMS.EncryptedValue `json:"employee_id"`
	EmployeeCode       pkgKMS.EncryptedValue `json:"employee_code"`
	WecomUserid        pkgKMS.EncryptedValue `json:"wecom_userid"`
	StoreCode          string                `json:"store_code"`
	Department         string                `json:"department"`
	OrganizationCode   string                `json:"organization_code"`
	MobileNumber       pkgKMS.EncryptedValue `json:"mobile_number"`
	EmailAddress       pkgKMS.EncryptedValue `json:"email_address"`
	Status             string                `json:"status"`
	JobPosition        string                `json:"job_position"`
	HireDate           string                `json:"hire_date"`
	TerminationDate    string                `json:"termination_date"`
	WecomJoinedAt      int64                 `json:"wecom_joined_at"`
	VerificationStatus int                   `json:"verification_status"`
	WecomLeftAt        int64                 `json:"wecom_left_at"`
	DepartmentChanges  []DepartmentChange    `json:"department_changes"`
	JobFunction        string                `json:"job_function"`
	IsStoreManagement  string                `json:"is_store_management"`
	UpdatedAt          int64                 `json:"updated_at"`
}
