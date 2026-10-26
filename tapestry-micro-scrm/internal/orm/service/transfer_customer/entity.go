package transfer_customer

type JobWorkExternalUser struct {
	Id      int64  `json:"id"`
	Openid  string `json:"openid"`
	Unionid string `json:"unionid"`
}

type JobWorkExternalAdminBaUser struct {
	Id     int64  `json:"id"`
	Name   string `json:"name"`
	UserId string `json:"user_id"`
}
