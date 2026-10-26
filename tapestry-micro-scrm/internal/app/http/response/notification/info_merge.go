package notificationResp

import microScrmModel "micro-scrm/internal/app/models/micro_scrm"

type InfoMergeDetailResp struct {
	Id           int64              `json:"id"`
	Type         int                `json:"type"`
	WeekNum      int                `json:"week_num"`
	Title        string             `json:"title"`
	Cover        string             `json:"cover" `
	SendType     int                `json:"send_type"`
	SendTime     int64              `json:"send_time"`
	Status       int                `json:"status"`
	MemberName   string             `json:"member_name"`
	CrowdType    int                `json:"crowd_type"`
	PdfUrl       string             `json:"pdf_url"`
	SendNameList []string           `json:"send_name_list"`
	EmailLoopUrl string             `json:"email_loop_url"`
	CrowdList    []MergeCrowdResp   `json:"crowd_list"`
	ContentList  []MergeContentResp `json:"content_list"`
	InfoMergeMemberNumResp
}

type MergeContentResp struct {
	DepartmentName string    `json:"department_name"`
	InfoList       []InfoRes `json:"info_list"`
	Sort           int       `json:"sort"`
}

type MergeCrowdResp struct {
	Id           int64  `json:"id"`
	Name         string `json:"title"`
	DepartmentId int64  `json:"department_id"`
}

type InfoMergeListResp struct {
	microScrmModel.NotificationInfoMergeModel
	InfoMergeMemberNumResp
}

type InfoMergeMemberNumResp struct {
	MergeId       int64 `json:"merge_id"`
	SendNum       int   `json:"send_num"`
	ViewPeopleNum int   `json:"view_people_num"`
	ViewNum       int   `json:"view_num"`
}

type InfoMergeListByH5Resp struct {
	ID       int64  `json:"id"`
	Type     int    `json:"type"`
	Title    string `json:"title"`
	Cover    string `json:"cover"`
	SendTime int64  `json:"send_time"`
	RoName   string `json:"ro_name"`
	ViewNum  int    `json:"view_num"`
	IsRead   bool   `json:"is_read"`
	PdfUrl   string `json:"pdf_url"`
	InfoMergeMemberNumResp
}
