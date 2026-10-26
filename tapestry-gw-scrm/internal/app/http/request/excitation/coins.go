package excitationReq

import "gw-scrm/internal/app/http/request"

type CoinsRecordReq struct {
	UserId       string `json:"user_id"`
	UserName     string `json:"user_name"`
	ThirdId      int64  `json:"third_id"`
	EnterpriseId int64  `json:"enterprise_id"`
	SourceType   int64  `json:"source_type"`
	ChangeType   int    `json:"change_type"`
	StartTime    int64  `json:"start_time"`
	EndTime      int64  `json:"end_time"`
	Page         int    `json:"page"`
	Size         int    `json:"size"`
}

type GetCoinSettingInfoReq struct {
	EnterpriseId int64 `json:"enterprise_id"`
	HasRecord    bool  `json:"has_record"`
}

type SaveCoinsSettingsReq struct {
	request.EnterpriseBaseReq
	MallIsOpen        int32 `json:"mall_is_open"`
	RatioPoints       int   `json:"ratio_points"`
	RatioCoins        int   `json:"ratio_coins"`
	CoinsExpireIsOpen int32 `json:"coins_expire_is_open"`
	CoinsExpireTime   int64 `json:"coins_expire_time"`
	AdminId           int64 `json:"admin_id"`
}
