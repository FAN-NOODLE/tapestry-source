package notificationReq

import (
	"micro-scrm/internal/app/http/request"
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
	Content      []InfoMergeContentReq `json:"content"`
	PdfUrl       string                `json:"pdf_url"`
	EmailLoopUrl string                `json:"email_loop_url"`
}

type InfoMergeContentReq struct {
	DepartmentName string  `json:"department_name"`
	InfoIds        []int64 `json:"info_ids"`
}

type InfoMergeDetailReq struct {
	request.EnterpriseBaseReq
	IsH5     bool   `json:"is_h5"`
	Id       int64  `json:"id"`
	MemberId int64  `json:"member_id"`
	Mobile   string `json:"mobile"`
}

type InfoMergeListReq struct {
	AdminManager
	IsH5       bool   `json:"is_h5"`
	Title      string `json:"title"`
	MemberName string `json:"member_name"`
	StartTime  int64  `json:"start_time"`
	EndTime    int64  `json:"end_time"`
	Type       int64  `json:"type"`
	Page       int    `json:"page"`
	Size       int    `json:"size"`
}

type SavePdfReq struct {
	request.EnterpriseBaseReq
	MergeId int64  `json:"merge_id"`
	InfoId  int64  `json:"info_id"`
	PdfUrl  string `json:"pdf_url"`
}

type SendFileByEmailReq struct {
	request.EnterpriseBaseReq
	Id         int64  `json:"id"`
	MergeId    int64  `json:"merge_id"`
	FileUrl    string `json:"file_url"`
	Email      string `json:"email" binding:"required"`
	Rename     string `json:"rename"`
	Attachment int    `json:"attachment"`
}

const IsAttachment = 1 //附件下载
