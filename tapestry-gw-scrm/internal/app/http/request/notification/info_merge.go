package notificationReq

import (
	"gw-scrm/internal/app/http/request"
)

type InfoMergeReq struct {
	request.EnterpriseBaseReq
	Id           int64                 `json:"id"`
	Type         int                   `json:"type"`
	WeekNum      int                   `json:"week_num"`
	Title        string                `json:"title"`
	Cover        string                `json:"cover" `
	CrowdType    int                   `json:"crowd_type"`
	CrowdIds     []int64               `json:"crowd_ids"`
	SendType     int                   `json:"send_type"`
	SendTime     int64                 `json:"send_time"`
	AdminId      int64                 `json:"admin_id"`
	Mobile       string                `json:"mobile"`
	EmailLoopUrl string                `json:"email_loop_url"`
	Content      []InfoMergeContentReq `json:"content"`
}

type InfoMergeContentReq struct {
	DepartmentName string  `json:"department_name"`
	InfoIds        []int64 `json:"info_ids"`
}

type InfoMergeDetailReq struct {
	request.EnterpriseBaseReq
	AdminID int64  `json:"admin_id"`
	Id      int64  `json:"id"`
	Mobile  string `json:"mobile"`
}

type InfoMergeListReq struct {
	request.EnterpriseBaseReq
	AdminID    int64  `json:"admin_id"`
	Title      string `json:"title"`
	MemberName string `json:"member_name"`
	StartTime  int64  `json:"start_time"`
	EndTime    int64  `json:"end_time"`
	Page       int    `json:"page"`
	Size       int    `json:"size"`
}

type SavePdfReq struct {
	request.EnterpriseBaseReq
	MergeId int64  `json:"merge_id"`
	InfoId  int64  `json:"info_id"`
	PdfUrl  string `json:"pdf_url"`
}
