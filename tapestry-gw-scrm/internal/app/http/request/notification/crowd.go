package notificationReq

import (
	"gw-scrm/internal/app/http/request"
)

type CrowdReq struct {
	request.EnterpriseBaseReq
	Id        int64   `json:"id"`
	Name      string  `json:"name"`
	Type      int     `json:"type"`
	FileUrl   string  `json:"file_url"`
	Mobile    string  `json:"mobile"`
	AdminId   int64   `json:"admin_id"`
	AdminName string  `json:"admin_name"`
	MemberIds []int64 `json:"member_ids"`
}

type CrowdListReq struct {
	request.EnterpriseBaseReq
	Name      string `json:"name"`
	StartTime int64  `json:"start_time"`
	EndTime   int64  `json:"end_time"`
	Page      int    `json:"page"`
	Size      int    `json:"size"`
}

type CrowdDetailReq struct {
	request.EnterpriseBaseReq
	Id int64 `json:"id"`
}
