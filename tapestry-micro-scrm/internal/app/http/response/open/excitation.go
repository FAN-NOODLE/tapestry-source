package openResp

type MemberInfoResp struct {
	Coins  float64 `json:"coins"`
	Avatar string  `json:"avatar"`
	Name   string  `json:"name"`
}

type OrderInfoResp struct {
	OrderStatus int     `json:"order_status"`
	UserId      string  `json:"user_id"`
	OrderNum    string  `json:"order_num"`
	OrderCoins  float64 `json:"order_coins"`
	CreatedAt   int64   `json:"created_at"`
	UpdatedAt   int64   `json:"updated_at"`
}
