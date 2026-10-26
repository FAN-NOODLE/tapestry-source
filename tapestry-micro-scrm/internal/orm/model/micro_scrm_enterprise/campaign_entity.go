// campaign活动接收返回实体

package micro_scrm_enterprise

type ResultSignInDetail struct {
	ActivityID     int64  `json:"activity_id"`
	UserID         int64  `json:"user_id"`
	UnionID        string `json:"union_id"`
	WxPhone        string `json:"wx_phone"`
	SignInDay      int64  `json:"sign_in_day"`
	LastSignInTime string `json:"last_sign_in_time"`
}
