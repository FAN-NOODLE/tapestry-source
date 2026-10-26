package activityQRCodeReq

import "micro-scrm/internal/app/http/request"

type FrontBaseReq struct {
	request.EnterpriseBaseReq
}

type GetQrCodeReq struct {
	FrontBaseReq
	Unionid string  `json:"unionid" comment:"用户unionid"`
	State   string  `json:"state"  binding:"required"  comment:"企微state"`
	ThirdId int     `json:"third_id"`
	Lat     float64 `json:"lat"`
	Lng     float64 `json:"lng"`
}
