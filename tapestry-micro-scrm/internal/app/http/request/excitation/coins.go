package excitationReq

import "micro-scrm/internal/app/http/request"

const ChangeTypeInc = 1
const ChangeTypeDec = 2

type ManualCoinsReq struct {
	ManualCoinsData []ManualCoinsRow `json:"manual_coins_data"`
	EnterpriseId    int64            `json:"enterprise_id"`
}

type ManualCoinsRow struct {
	MemberId    int64   `json:"member_id"`
	UserId      string  `json:"user_id"`
	CoinsRemark string  `json:"coins_remark"`
	Coins       float64 `json:"coins"`
}

type CoinsRecordListReq struct {
	request.MemberSearchReq
	SourceType int64 `json:"source_type"`
	ChangeType int   `json:"change_type"`
	StartTime  int64 `json:"start_time"`
	EndTime    int64 `json:"end_time"`
	Page       int   `json:"page"`
	Size       int   `json:"size"`
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
type GetCoinSettingInfoReq struct {
	EnterpriseId int64 `json:"enterprise_id"`
	HasRecord    bool  `json:"has_record"`
}

type CoinExpiredReq struct {
	EnterpriseId int64 `json:"enterprise_id"`
	DelaySeconds int64 `json:"delay_seconds"`
}
