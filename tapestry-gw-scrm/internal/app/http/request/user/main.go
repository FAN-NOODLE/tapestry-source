package userReq

import "gw-scrm/internal/app/http/request"

type SSOLoginReq struct {
	request.EnterpriseBaseReq
	Code    string `json:"code" binding:"required"`
	AppCode string `json:"app_code"`
}

type SSOLoginResp struct {
	Mobile string `json:"mobile"`
}
