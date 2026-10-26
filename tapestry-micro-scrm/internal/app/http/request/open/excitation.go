package openReq

type MemberInfoReq struct {
	EnterpriseId int64  `json:"enterprise_id"`
	UserId       string `json:"user_id"`
}

type AddOrderReq struct {
	EnterpriseId   int64   `json:"enterprise_id"`
	UserId         string  `json:"user_id"`
	OrderNum       string  `json:"order_num"`
	DeductionCoins float64 `json:"deduction_coins"`
	Info           string  `json:"info"`
}

type CancelOrderReq struct {
	OrderInfoReq
}

type OrderInfoReq struct {
	EnterpriseId int64  `json:"enterprise_id"`
	UserId       string `json:"user_id"`
	OrderNum     string `json:"order_num"`
}
