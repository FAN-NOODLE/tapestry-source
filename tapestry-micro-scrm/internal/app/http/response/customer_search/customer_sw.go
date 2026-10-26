package customerSearchResp

type ExternalUserListSw struct {
	ID                    int64  `json:"id"`
	Openid                string `json:"openid"`
	Unionid               string `json:"unionid"`
	Nick                  string `json:"nick"`
	Name                  string `json:"name"`
	Avatar                string `json:"avatar"`
	CardLevel             string `json:"card_level"`
	AmountTotal           int    `json:"amount_total"`
	Remark                string `json:"remark"`
	VipType               int    `json:"vip_type"`
	ValidReserveAt        int    `json:"valid_reserve_at"`
	ValidReserveAtTag     bool   `json:"valid_reserve_at_tag"`
	SuccessReserveAt      int    `json:"success_reserve_at"`
	TaskFinishAt          int    `json:"task_finish_at"`
	TaskFinishAtTag       bool   `json:"task_finish_at_tag"`
	CouponStatus          int    `json:"coupon_status"`
	CouponBy30daysStatus  int    `json:"coupon_by_30days_status"`
	WorkExternalUserId    int64  `json:"work_external_user_id"`
	LastPurchaseDate      int64  `json:"last_purchase_date"`
	LastPurchaseStoreCode string `json:"last_purchase_store_code"`
	LastPurchaseStoreName string `json:"last_purchase_store_name"`
	UnionCode             string `json:"union_code"`
	Birthday              string `json:"birthday"`
	MemberName            string `json:"member_name"`
	CreatedAt             int64  `json:"created_at"`
	SelfRelationStatus    int64  `json:"self_relation_status"`
}

type SearchQuerySwResp struct {
	External     BuildSql
	RelationShip BuildSql
}

type BuildSqlSw struct {
	SqlConn  string        `json:"conn"`
	SqlSlice []interface{} `json:"sql_slice"`
}

type ExternalCrowdSwList struct {
	ExternalUserListSw
	IsContact int `json:"is_contact"`
}
