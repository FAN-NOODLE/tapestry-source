package notificationResp

type ListRes struct {
	ID             int64  `json:"id"`
	Title          string `json:"title"`
	Type           int32  `json:"type"`
	DepartmentName string `json:"department_name"`
	SubmitName     string `json:"submit_name"`
	UpdateAt       int64  `json:"update_at"`
	ReviewName     string `json:"review_name"`
	ReviewTime     int    `json:"review_time"`
	Status         int32  `json:"status"`
	IsDraft        int32  `json:"is_draft"`
	RelStatus      int    `json:"rel_status"`
	IsRead         bool   `json:"is_read"`
}

type InfoViewNumResp struct {
	ViewPeopleNum int `json:"view_people_num"`
	ViewNum       int `json:"view_num"`
}

type InfoRes struct {
	Id             int64            `json:"id"`
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
	AdminID        int64            `json:"admin_id"`
	Status         int32            `json:"status"`
	ReviewTime     int              `json:"review_time"`
	ReviewName     string           `json:"review_name"`
	SubmitTime     int              `json:"submit_time"`
	MemberID       int64            `json:"member_id"`
	MemberName     string           `json:"member_name"`
	Comment        string           `json:"comment"`
	PageUrl        string           `json:"page_url"`
	IsDraft        int32            `json:"is_draft"`
	TypeStr        string           `json:"type_str"`
	StatusStr      string           `json:"status_str"`
	PdfUrl         string           `json:"pdf_url"`
	InfoViewNumResp
}

type ContactList struct {
	DepartmentName string `json:"department_name"`
	Name           string `json:"name"`
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

type ManagerInfo struct {
	IsFT     int32 `json:"is_ft"`
	IsRo     int32 `json:"is_ro"`
	MemberId int64 `json:"member_id"`
	AdminId  int64 `json:"admin_id"`
}

type ListReadNum struct {
	InfoCheckingNum int64 `json:"info_checking_num"`
	InfoFailNum     int64 `json:"info_fail_num"`
	MergeInfoNum    int64 `json:"merge_info_num"`
}
