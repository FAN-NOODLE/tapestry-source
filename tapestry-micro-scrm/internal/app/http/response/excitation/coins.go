package excitationResp

import excitationModels "micro-scrm/internal/app/models/excitation"

type GetCoinsListResp struct {
	Id          int64  `json:"id"`
	CoinsRemark string `json:"coins_remark"`
	Coins       string `json:"coins"`
	CreatedAt   int64  `json:"created_at"`
}

type CoinsSettingInfoResp struct {
	MallIsOpen              int32                                                     `json:"mall_is_open"`
	RatioPoints             int                                                       `json:"ratio_points"`
	RatioCoins              int                                                       `json:"ratio_coins"`
	CoinsExpireIsOpen       int32                                                     `json:"coins_expire_is_open"`
	CoinsExpireTime         int64                                                     `json:"coins_expire_time"`
	CoinsExpireRecordList   []excitationModels.ExcitationCoinsExpireRecordModel       `json:"coins_expire_record_list"`
	CoinsSettingsEditRecord []excitationModels.ExcitationCoinsSettingsEditRecordModel `json:"coins_settings_edit_record"`
}
type CoinsExpireRecord struct {
	ID              int64 `json:"id"`
	CoinsExpireTime int64 `json:"coins_expire_time"`
}

type CoinsRecordList struct {
	Id             int64   `json:"id"`
	CoinsRemark    string  `json:"coins_remark"`
	Coins          float64 `json:"coins"`
	SourceType     int     `json:"source_type"`
	CreatedAt      int64   `json:"created_at"`
	CreatedTimeStr string  `json:"created_time_str"`
	SourceTypeStr  string  `json:"source_type_str"`
	OrderNum       string  `json:"order_num"`
	ChangeType     int     `json:"change_type"`
	ChangeTypeStr  string  `json:"change_type_str"`
	MemberInfo
}

type SpendCoinsResp struct {
	MemberId   int64   `json:"member_id"`
	SpendCoins float64 `json:"spend_coins"`
}
