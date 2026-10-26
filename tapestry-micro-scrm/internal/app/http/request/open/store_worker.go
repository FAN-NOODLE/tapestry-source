package openReq

type StoreWorkerSaListReq struct {
	EnterpriseId int64  `json:"enterprise_id"`
	Unionid      string `json:"unionid" binding:"required" required:"用户unionid"`
}

type StoreListReq struct {
	EnterpriseId       int64    `json:"enterprise_id"`
	Page               int      `json:"page"`
	Size               int      `json:"size"`
	StoreCodes         []string `json:"store_codes"`
	StoreName          string   `json:"store_name"`
	StoreType          int      `json:"store_type"`
	WechatDepartmentID int      `json:"wechat_department_id"`
}

type SaListReq struct {
	StoreListReq
}

type StoreWorkListReq struct {
	EnterpriseId int64  `json:"enterprise_id"`
	UserId       string `json:"user_id"`
	StartTime    int64  `json:"start_time"`
	EndTime      int64  `json:"end_time"`
	Page         int    `json:"page"`
	Size         int    `json:"size"`
}

type CustomerStoreWorkListReq struct {
	EnterpriseId int64  `json:"enterprise_id"`
	UnionId      string `json:"unionid" binding:"required" required:"用户unionid必填"`
	Page         int    `json:"page"`
	Size         int    `json:"size"`
}

type GenQrCodeReq struct {
	EnterpriseId int64  `json:"enterprise_id"`
	UserId       string `json:"user_id" binding:"required" required:"用户user_id必填"`
	State        string `json:"state"`
}

type AreaTreeReq struct {
	EnterpriseId  int64 `json:"enterprise_id"`
	ParentId      int   `json:"parent_id"`
	HasChildren   bool  `json:"has_children"`
	ExcludeClosed bool  `json:"exclude_closed"`
}

type AreaManagerReq struct {
	EnterpriseId int64  `json:"enterprise_id"`
	UserId       string `json:"user_id"`
	MemberId     int64  `json:"member_id"`
}

type HkStoreListReq struct {
	EnterpriseId       int64    `json:"enterprise_id"`
	Page               int      `json:"page"`
	Size               int      `json:"size"`
	StoreCodes         []string `json:"store_codes"`
	StoreName          string   `json:"store_name"`
	Area               string   `json:"area"`
	WechatDepartmentID int      `json:"wechat_department_id"`
	ExcludeClosed      bool     `json:"exclude_closed"`
}

type HkStoreWorkListReq struct {
	EnterpriseId       int64    `json:"enterprise_id"`
	UserId             string   `json:"user_id"`
	StartTime          int64    `json:"start_time"`
	EndTime            int64    `json:"end_time"`
	StoreCodes         []string `json:"store_codes"`
	Area               string   `json:"area"`
	WechatDepartmentID int      `json:"wechat_department_id"`
	Page               int      `json:"page"`
	Size               int      `json:"size"`
}

type HkAreaManagerReq struct {
	EnterpriseId int64  `json:"enterprise_id"`
	UserId       string `json:"user_id"`
	MemberId     int64  `json:"member_id"`
}
