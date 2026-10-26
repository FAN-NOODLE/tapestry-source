package customerSearchResp

type ExternalUserList struct {
	ID                   int64  `json:"id"`
	Openid               string `json:"openid"`
	Unionid              string `json:"unionid"`
	Nick                 string `json:"nick"`
	Name                 string `json:"name"`
	Avatar               string `json:"avatar"`
	BaCode               string `json:"ba_code"`
	CardLevel            string `json:"card_level"`
	AmountTotal          int    `json:"amount_total"`
	Remark               string `json:"remark"`
	RelationStatus       int    `json:"relation_status"`
	SelfRelationStatus   int64  `json:"self_relation_status"`
	ActiveTime           int64  `json:"active_time"`
	InteractTime         int64  `json:"interact_time"`
	KeySaName            string `json:"key_sa_name"`
	VipType              int    `json:"vip_type"`
	ValidReserveAt       int    `json:"valid_reserve_at"`
	ValidReserveAtTag    bool   `json:"valid_reserve_at_tag"`
	SuccessReserveAt     int    `json:"success_reserve_at"`
	TaskFinishAt         int    `json:"task_finish_at"`
	TaskFinishAtTag      bool   `json:"task_finish_at_tag"`
	CouponStatus         int    `json:"coupon_status"`
	CouponBy30daysStatus int    `json:"coupon_by_30days_status"`
	WorkExternalUserId   int64  `json:"work_external_user_id"`
	MemberId             int64  `json:"member_id"`
	UserId               string `json:"user_id"`
	Createtime           int64  `json:"createtime"`
	IsKeySa              int    `json:"is_key_sa"`
	LastPurchaseDate     int64  `json:"last_purchase_date"`
	CltType              string `json:"clt_type"`
	UnionCode            string `json:"union_code"`
}

type SearchQueryResp struct {
	ExternalExtra BuildSql
	RelationShip  BuildSql
}

type BuildSql struct {
	SqlConn  string        `json:"conn"`
	SqlSlice []interface{} `json:"sql_slice"`
}

type ExternalDetailResp struct {
	ID                          int64      `json:"id"`
	Openid                      string     `json:"openid"`
	Unionid                     string     `json:"unionid"`
	Nick                        string     `json:"nick"`
	Name                        string     `json:"name"`
	CardLevel                   string     `json:"card_level"`
	Remark                      string     `json:"remark"`
	Description                 string     `json:"description"`
	Avatar                      string     `json:"avatar"`
	Birthday                    string     `json:"birthday"`
	UnionCode                   string     `json:"union_code"`
	LastPurchaseStoreCode       string     `json:"last_purchase_store_code"`
	HighWilling                 string     `json:"high_willing" comment:"高意愿客户 Y是 N否"`
	PriceSensitivity            string     `json:"price_sensitivity" comment:"价格敏感度 H高 L低"`
	CategoryPreference          []string   `json:"category_preference"`
	ChatTag                     []ChatList `json:"chat_tag" comment:"顾客沟通意愿标签"`
	RelationStatus              int        `json:"relation_status"`
	ValidReserveAt              int        `json:"valid_reserve_at"`
	TaskFinishAt                int        `json:"task_finish_at"`
	InteractTime                int64      `json:"interact_time" comment:"最近一次联系时间"`
	ActiveTime                  int64      `json:"active_time" comment:"最近一次活跃沟通时间"`
	AmountTotal                 int        `json:"amount_total"`
	HistoryAmountTotal          int        `json:"history_amount_total"`
	LastPurchaseDate            int64      `json:"last_purchase_date"`
	AgeTag                      string     `json:"age_tag"`
	CdpTagBuyCategoryPreference []PieChart `json:"cdp_tag_buy_category_preference"`
	CdpTagPurchasedMonth        []PieChart `json:"cdp_tag_purchased_month"`
}

type ExternalCrowdList struct {
	ExternalUserList
	IsContact int `json:"is_contact"`
}

type PieChart struct {
	Name    string  `json:"name"`
	Percent float64 `json:"percent"`
}

type ChatList struct {
	Key      int64  `json:"key"`
	Value    string `json:"value"`
	IsSelect bool   `json:"is_select"`
}
