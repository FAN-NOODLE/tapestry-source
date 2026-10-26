package customerReq

const (
	// SowTaskListTemplateCustom SowTaskListTemplateAll      = iota // 全部
	SowTaskListTemplateCustom   = iota + 1 // 自定义
	SowTaskListTemplateTemplate            // 模版
)

type SowTask struct {
	EnterpriseID       int64  `json:"enterprise_id" binding:"required" required:"企业Id"`
	MemberMobile       string `json:"member_mobile"`
	WorkExternalUserID int64  `json:"work_external_user_id"`
	Status             int    `json:"status"`
}

type SowTaskList struct {
	SowTask
	Time                  int64    `json:"time"`
	SelfManager           uint8    `json:"self_manager"`
	Username              string   `json:"username"`
	Nickname              string   `json:"nickname"`
	ImpressionTagIds      []int64  `json:"impression_tag_ids"`
	WorkExternalTagIds    []int64  `json:"work_external_tag_ids"`
	CltType               []string `json:"clt_type"  comment:"clt类型"`
	CardType              []string `json:"card_type"  comment:"卡等"`
	Page                  int      `json:"page" binding:"min=1" min:"页码"`
	Size                  int      `json:"size" binding:"min=10,max=9999" min:"分页数"`
	BirthdayStart         int      `json:"birthday_start"`
	BirthdayEnd           int      `json:"birthday_end"`
	LastPurchaseStartTime int      `json:"last_purchase_start_time"`
	LastPurchaseEndTime   int      `json:"last_purchase_end_time"`
}

type SowTaskContentBase struct {
	Content    string `json:"content" binding:"required" required:"内容"`
	RemindTime int64  `json:"remind_time" binding:"required" required:"提醒时间"`
}

type SowTaskCreate struct {
	SowTask
	SowTaskContentBase
}

type SowTaskUpdate struct {
	SowTaskCreate
	Id     int64 `json:"id" binding:"required" required:"id"`
	Status int   `json:"status"`
}

type SowTaskFinished struct {
	Id     int64  `json:"id" binding:"required" required:"id"`
	Mobile string `json:"mobile"`
	SowTask
}

type SowTemplateSaveReq struct {
	ID                int64               `json:"id"`
	EnterpriseID      int64               `json:"enterprise_id" binding:"required" required:"企业Id"`
	Mobile            string              `json:"mobile"`
	Name              string              `json:"name" binding:"required" required:"模版名称"`
	Theme             string              `json:"theme"`
	Status            int                 `json:"status"`
	AdminName         string              `json:"admin_name"`
	SowTemplateNotice []SowTemplateNotice `json:"sow_template_notice"`
}

type SowTemplateNotice struct {
	Id           int64  `json:"id" binding:"required" required:"id"`
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
	Status       int   `json:"status"`
	EnterpriseID int64 `json:"enterprise_id"`
}

type SowAllTaskListReq struct {
	Page            int     `json:"page"`
	Size            int     `json:"size"`
	EnterpriseID    int64   `json:"enterprise_id" binding:"required" required:"企业Id"`
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

type SowTaskNewCreate struct {
	SowTask
	TaskContent []SowTaskContentBase `json:"task_content"`
	TemplateId  int64                `json:"template_id"`
}
