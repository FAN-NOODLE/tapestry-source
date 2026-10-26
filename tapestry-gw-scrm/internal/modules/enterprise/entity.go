package enterprise

// ResponseBase 基础返回结构
type ResponseBase struct {
	App       string `json:"app"`
	Code      int    `json:"code"`
	Msg       string `json:"msg"`
	RequestId string `json:"request_id"`
}

// Enterprise 获取组织信息请求结构
type Enterprise struct {
	EnterpriseId int64 `json:"enterprise_id"`
}

// EnterpriseResponse 企业信息返回
type EnterpriseResponse struct {
	ID                int64  `json:"id"`
	EnterpriseID      int64  `json:"enterprise_id"`
	Hash              string `json:"hash"`
	Name              string `json:"name"`
	Logo              string `json:"logo"`
	Contact           string `json:"contact"`
	Mobile            string `json:"mobile"`
	Address           string `json:"address"`
	IndustryID        int64  `json:"industry_id"`
	SeparateSwitch    int64  `json:"separate_switch"`
	LockTime          int32  `json:"lock_time"`
	Proportion        int64  `json:"proportion"`
	StartFee          int64  `json:"start_fee"`
	CrsStatus         int    `json:"crs_status"`
	ContactStatus     int    `json:"contact_status"`
	CorpId            string `json:"corp_id"`
	Domain            string `json:"domain"`
	ExternalSecret    string `json:"external_secret"`
	Status            int    `json:"status"`
	BusinessSource    int    `json:"business_source"`
	Version           int    `json:"version"`
	VersionTag        string `json:"version_tag"`
	AccountExpireAt   int64  `json:"account_expire_at"`
	AccountStartAt    int64  `json:"account_start_at"`
	AccountLimit      int32  `json:"account_limit"`
	ErpStatus         int    `json:"erp_status"`
	CrmStatus         int    `json:"crm_status"`
	ReturnAddress     string `json:"return_address"`
	IsSplitFlow       int    `json:"is_split_flow"`
	AvailableSystem   string `json:"available_system"`
	ServiceMerchantId int64  `json:"service_merchant_id"`
	EditAt            string `json:"edit_at"`
	CreatedAt         string `json:"created_at"`
	UpdatedAt         int64  `json:"updated_at"`
	DeletedAt         int64  `json:"deleted_at"`
	Identifier        string `json:"identifier"`
}
