package marketing

const (
	CouponMatchFailWeixin           = 101
	CouponMatchFailUserGroup        = 102
	CouponMatchFailWeixinNoData     = 104
	CouponMatchFailCustomerChannel  = 105
	CouponMatchFailTriggerCondition = 106
	CouponMatchFailSourceNoData     = 107
	CouponMatchFailTimeOut          = 500
)

var (
	GetAllShortUrlSql = "SELECT id,one_level_name,two_level_name,page_name,operator,short_url,weixin_url,created_at,url_type FROM external_links WHERE enterprise_id = ?"
)
