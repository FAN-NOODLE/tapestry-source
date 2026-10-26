package notificationReq

import (
	"gw-work-scrm/internal/app/http/request"
)

type InfoMergeDetailReq struct {
	request.EnterpriseBaseReq
	Id       int64 `json:"id"`
	MemberId int64 `json:"member_id"`
	IsH5     bool  `json:"is_h5"`
}

type InfoMergeListReq struct {
	request.EnterpriseBaseReq
	MemberId int64  `json:"member_id"`
	Title    string `json:"title"`
	Type     int64  `json:"type"`
	IsH5     bool   `json:"is_h5"`
	Page     int    `json:"page"`
	Size     int    `json:"size"`
}
