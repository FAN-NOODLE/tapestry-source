package customerReq

type SowTemplateSaveReq struct {
	ID                int64               `json:"id"`
	EnterpriseID      int64               `json:"enterprise_id"`
	Mobile            string              `json:"mobile"`
	Name              string              `json:"name" binding:"required" required:"模版名称"`
	Theme             string              `json:"theme"`
	Status            int                 `json:"status"`
	AdminName         string              `json:"admin_name"`
	SowTemplateNotice []SowTemplateNotice `json:"sow_template_notice"`
}

type SowTemplateNotice struct {
	Id           int64  `json:"id"`
	EnterpriseID int64  `json:"enterprise_id"`
	TemplateID   int64  `json:"template_id"`
	DelayType    int    `json:"delay_type" binding:"required" required:"迟提醒类型：1天 2周 3月"`
	DelayValue   int    `json:"delay_value" binding:"required" required:"延迟提醒的时间数值"`
	DelayTime    string `json:"delay_time" binding:"required" required:"延迟提醒当天的提醒时间，时:分"`
	Content      string `json:"content" binding:"required" required:"提醒文案"`
}

type SowTemplateListReq struct {
	EnterpriseID int64  `json:"enterprise_id"`
	Status       int    `json:"status"`
	Name         string `json:"name"`
	StartTime    int64  `json:"start_time"`
	EndTime      int64  `json:"end_time"`
	Page         int    `json:"page"`
	Size         int    `json:"size"`
}

type SowTemplateDetailReq struct {
	Id           int64 `json:"id" binding:"required" required:"id"`
	EnterpriseID int64 `json:"enterprise_id"`
	Status       int   `json:"status"`
}

type SowAllTaskListReq struct {
	Page            int     `json:"page"`
	Size            int     `json:"size"`
	EnterpriseID    int64   `json:"enterprise_id"`
	TemplateName    string  `json:"template_name"`
	TemplateId      int64   `json:"template_id"`
	TemplateType    int     `json:"template_type"`
	SaName          string  `json:"sa_name"`
	Username        string  `json:"username"`
	UnionCode       string  `json:"union_code"`
	SowStatus       int     `json:"sow_status"`
	StartCreatedAt  int64   `json:"start_created_at"`
	EndCreatedAt    int64   `json:"end_created_at"`
	StartNoticeAt   int64   `json:"start_notice_at"`
	EndNoticeAt     int64   `json:"end_notice_at"`
	StartFinishedAt int64   `json:"start_finished_at"`
	EndFinishedAt   int64   `json:"end_finished_at"`
	MemberIds       []int64 `json:"member_ids"`
	ExternalIds     []int64 `json:"external_ids"`
	TemplateIds     []int64 `json:"template_ids"`
}

type SowTemplateList struct {
	Total int64             `json:"total"`
	Page  int               `json:"page"`
	Size  int               `json:"size"`
	List  []SowTemplateInfo `json:"list"`
}
type SowTemplateInfo struct {
	Id           int64  `json:"id"`
	Name         string `json:"name"`
	TaskNum      int    `json:"task_num"`
	UseNum       int    `json:"use_num"`
	UseSaNum     int    `json:"use_sa_num"`
	AdminName    string `json:"admin_name"`
	CreatedAt    int64  `json:"created_at"`
	CreatedAtStr string `json:"created_at_str"`
	UpdatedAt    int64  `json:"updated_at"`
	UpdatedAtStr string `json:"updated_at_str"`
}
