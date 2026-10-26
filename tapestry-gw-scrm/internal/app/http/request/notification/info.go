package notificationReq

import (
	"gw-scrm/internal/app/http/request"
)

type BaseReq struct {
	AdminManager
	Mobile string `json:"mobile"`
}

type InfoReq struct {
	SearchId
	ThemeTitle     string           `json:"theme_title"`
	Title          string           `json:"title"`
	UntilName      string           `json:"until_name"`
	Content        string           `json:"content"`
	Perform        string           `json:"perform" `
	DateTime       string           `json:"date_time"`
	ContactInfo    []ContactList    `json:"contact_info"`
	Attachment     []AttachmentInfo `json:"attachment"`
	Type           int32            `json:"type"`
	DepartmentName string           `json:"department_name"`
	Status         int32            `json:"status"`
	ReviewTime     int              `json:"review_time"`
	MemberID       int64            `json:"member_id"`
	Comment        string           `json:"comment"`
	PageUrl        string           `json:"page_url"`
	IsDraft        int32            `json:"is_draft"`
	IsH5           bool             `json:"is_h5"`
	PdfUrl         string           `json:"pdf_url"`
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
	AdminManager
	Title            string `json:"title"`
	UpdatedStartTime int    `json:"updated_start_time"`
	UpdatedEndTime   int    `json:"updated_end_time"`
	ReviewStartTime  int    `json:"review_start_time"`
	ReviewEndTime    int    `json:"review_end_time"`
	Type             int32  `json:"type"`
	DepartmentName   string `json:"department_name"`
	Status           int32  `json:"status"`
	NoRelease        int    `json:"no_release"`
	IsAll            int32  `json:"is_all"`
	request.PageReq
}
type SearchId struct {
	BaseReq
	Id int64 `json:"id"`
}

type AdminManager struct {
	request.EnterpriseBaseReq
	AdminID int64 `json:"admin_id"`
}
