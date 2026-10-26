package customerReq

type ListReq struct {
	SowTask
	Page int `json:"page"`
	Size int `json:"size"`
}

type SowTask struct {
	EnterpriseID       int64  `json:"enterprise_id"`
	MemberMobile       string `json:"member_mobile"`
	WorkExternalUserID int64  `json:"work_external_user_id" required:"联系人id"`
	Status             int    `json:"status"`
}

type SowTaskCreate struct {
	SowTask
	Content    string `json:"content" binding:"required" required:"内容"`
	RemindTime int64  `json:"remind_time" binding:"required" required:"提醒时间"`
}

type SowTaskUpdate struct {
	SowTaskCreate
	Id     int64 `json:"id" binding:"required" required:"id"`
	Status int   `json:"status"`
}

type SowTaskFinished struct {
	Id int64 `json:"id" binding:"required" required:"id"`
	SowTask
}
