package notificationReq

import (
	"gw-work-scrm/internal/app/http/request"
)

type BaseReq struct {
	request.EnterpriseBaseReq
	Mobile string `json:"mobile"`
}

type InfoReq struct {
	SearchId
	Status  int32  `json:"status"`
	Comment string `json:"comment"`
	IsH5    bool   `json:"is_h5"`
}

type ContactList struct {
	DepartmentName string `json:"department_name"`
	Name           string `json:"name" binding:"required"`
	Email          string `json:"email"`
	TelPhone       string `json:"tel_phone"`
	Remark         string `json:"remark"`
}

type AttachmentInfo struct {
	Name           string `json:"name"`
	AttachmentName string `json:"attachment_name"`
	Type           string `json:"type"`
	Size           int64  `json:"size"`
	Url            string `json:"url"`
}

type SearchListReq struct {
	request.EnterpriseBaseReq
	MemberId         int64  `json:"member_id"`
	Title            string `json:"title"`
	UpdatedStartTime int    `json:"updated_start_time"`
	UpdatedEndTime   int    `json:"updated_end_time"`
	ReviewStartTime  int    `json:"review_start_time"`
	ReviewEndTime    int    `json:"review_end_time"`
	Type             int32  `json:"type"`
	DepartmentName   string `json:"department_name"`
	Status           int32  `json:"status"`
	IsAll            int32  `json:"is_all"`
	request.PageReq
}
type SearchId struct {
	BaseReq
	MemberId int64 `json:"member_id"`
	Id       int64 `json:"id"`
}

type AdminManager struct {
	request.EnterpriseBaseReq
	MemberId int64 `json:"member_id"`
}

type SendFileByEmailReq struct {
	request.EnterpriseBaseReq
	Id         int64  `json:"id"`
	MergeId    int64  `json:"merge_id"`
	FileUrl    string `json:"file_url"`
	Email      string `json:"email"`
	Rename     string `json:"rename"`
	Attachment int    `json:"attachment"`
}

type SavePdfReq struct {
	request.EnterpriseBaseReq
	MergeId int64  `json:"merge_id"`
	InfoId  int64  `json:"info_id"`
	PdfUrl  string `json:"pdf_url"`
}
