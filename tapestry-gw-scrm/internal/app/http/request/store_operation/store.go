package store_operation

type StoreListReq struct {
	EnterpriseID  int64    `json:"enterprise_id"`
	Page          int      `json:"page"`
	Size          int      `json:"size"`
	CountryCodes  []string `json:"country_codes"`
	StoreName     string   `json:"store_name"`
	StoreCode     string   `json:"store_code"`
	Channel       string   `json:"channel"`
	Status        string   `json:"status"`
	OperatingMode string   `json:"operating_mode"`
	PaymentMethod string   `json:"payment_method"`
}

type StoreDetailReq struct {
	EnterpriseID int64 `json:"enterprise_id"`
	ID           int64 `json:"id" binding:"required"`
}

type StoreUpdateReq struct {
	EnterpriseID            int64    `json:"enterprise_id"`
	ID                      int64    `json:"id" binding:"required"`
	StoreNameZh             string   `json:"store_name_zh" binding:"max=255"`
	AddressZh               string   `json:"address_zh" binding:"max=500"`
	Phone                   string   `json:"phone" binding:"max=64"`
	Email                   string   `json:"email" binding:"max=255"`
	AppointmentServiceTypes []string `json:"appointment_service_types"`
	OperatorAdminID         int64    `json:"operator_admin_id"`
	OperatorName            string   `json:"operator_name"`
}
