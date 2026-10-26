package entity

type BossEnterpriseDetailResponse struct {
	Id                int64  `json:"id"`
	Hash              string `json:"hash"`
	Name              string `json:"name"`
	Logo              string `json:"logo"`
	Contact           string `json:"contact"`
	Mobile            string `json:"mobile"`
	Address           string `json:"address"`
	IndustryId        int64  `json:"industry_id"`
	SeparateSwitch    int8   `json:"separate_switch"`
	LockTime          int32  `json:"lock_time"`
	Proportion        int64  `json:"proportion"`
	StartFee          string `json:"start_fee"`
	CrsStatus         string `json:"crs_status"`
	StoreStatus       int8   `json:"store_status"`
	ContactStatus     int8   `json:"contact_status"`
	CorpId            string `json:"corp_id"`
	Domain            string `json:"domain"`
	ExternalSecret    string `json:"external_secret"`
	Status            int8   `json:"status"`
	BusinessSource    int64  `json:"business_source"`
	Version           int    `json:"version"`
	VersionTag        string `json:"version_tag"`
	AccountExpireAt   int64  `json:"account_expire_at"`
	AccountStartAt    int64  `json:"account_start_at"`
	AccountLimit      int32  `json:"account_limit"`
	ErpStatus         int8   `json:"erp_status"`
	CrmStatus         int8   `json:"crm_status"`
	ReturnAddress     string `json:"return_address"`
	IsSplitFlow       int8   `json:"is_split_flow"`
	AvailableSystem   string `json:"available_system"`
	ServiceMerchantId int64  `json:"service_merchant_id"`
	EditAt            int64  `json:"edit_at"`
	CreatedAt         int64  `json:"created_at"`
	UpdatedAt         int64  `json:"updated_at"`
	DeletedAt         int64  `json:"deleted_at"`
	Identifier        string `json:"identifier"`
}
