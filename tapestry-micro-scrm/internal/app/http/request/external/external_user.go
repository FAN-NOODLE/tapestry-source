package external

import microScrmModel "micro-scrm/internal/app/models/micro_scrm"

type CdpExternalInfo struct {
	ID                    string `json:"id" comment:"cdp的customer的id"`
	CustomerCode          string `json:"customer_code" comment:"cdp的customer_code"`
	UnionId               string `json:"union_id" comment:"外部联系人的unionid"`
	Name                  string `json:"name" comment:"外部联系人的姓名"`
	Last12MPurchaseAmount int64  `json:"last_12m_purchase_amount" comment:"用户最近12个月的消费金额"`
	Last36MPurchaseAmount int64  `json:"last_36m_purchase_amount" comment:"用户最近36个月的消费金额"`
	MemberShip            string `json:"member_ship" comment:"用户的会员卡号"`
	ClientType            string `json:"client_type" comment:"用户的client_type类型"`
	UnionIdBindFlag       int    `json:"union_id_bind_flag" comment:"用户的会员是否绑定，1为是，0为否--解绑"`

	CustomerLevel         string `json:"customer_level" comment:"cdp的会员等级"`
	CustomerType          int    `json:"customer_type" comment:"cdp的会员类型"`
	LastPurchaseDate      int    `json:"last_purchase_date" comment:"最近一次订单消费时间"`
	LastPurchaseStoreCode string `json:"last_purchase_store_code" comment:"最近一次订单消费的门店code"`
	UpgradeNeedAmount     int64  `json:"upgrade_need_amount" comment:"升级下一等级所需金额"`
	Birthday              string `json:"birthday" comment:"顾客生日"`
	BirthYear             string `json:"birth_year" comment:"顾客生日"`
	BirthMonth            string `json:"birth_month" comment:"顾客生日"`
	BirthDay              string `json:"birth_day" comment:"顾客生日"`
	VipType               int    `json:"vip_type"  comment:"是否vip"`
	RemindExpiredPoints   int    `json:"remind_expired_points" comment:"下个月有积分将到期 0：不存在下个月积分过期 1：存在下个月积分过期"`
	CustomerTypeUpdatedAt int    `json:"customer_type_updated_at" comment:"cdp的会员类型更新时间"`
	AgeTag                string `json:"age_tag" comment:"‘Z世代’ 或 ‘千禧一代’"`
}

type WxUserCorpExternalInfo struct {
	microScrmModel.WxUserCorpExternalModel
	IsKeySa    int    `json:"is_key_sa" db:"is_key_sa"`   // 是否专属导购：0否，1是
	SaUserId   string `json:"sa_user_id" db:"sa_user_id"` // 是否专属导购：0否，1是
	SaUserName string `json:"sa_user_name" db:"sa_user_name"`
}

type CoachNameInfo struct {
	Unionid string `json:"unionid"`
	Nick    string `json:"nick"`
	Remark  string `json:"remark"`
}
