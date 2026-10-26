package customerResp

type SowTaskList struct {
	Id                 int64  `json:"id"`
	WorkExternalUserID int64  `json:"work_external_user_id"`
	MemberID           int64  `json:"member_id"`
	MemberUserID       string `json:"member_user_id"`
	Content            string `json:"content"`
	Status             int32  `json:"status"`
	RemindTime         int64  `json:"remind_time"`
	Username           string `json:"username"`
	Nickname           string `json:"nickname"`
	Grade              string `json:"grade"`
	Birthday           string `json:"birthday"`
	OpenId             string `json:"open_id"`
	IsKsa              int    `json:"is_ksa"`
}

type SowTemplateInfo struct {
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	TaskNum   int    `json:"task_num"`
	Status    int    `json:"status"`
	UseNum    int    `json:"use_num"`
	UseSaNum  int    `json:"use_sa_num"`
	AdminName string `json:"admin_name"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

type SowSaNum struct {
	TemplateId int64 `json:"template_id"`
	SaNum      int   `json:"sa_num"`
}

type SowInfoResp struct {
	Id            int64  `json:"id"`
	TemplateName  string `json:"template_name"`
	TemplateId    int64  `json:"template_id"`
	TemplateIdStr string `json:"template_id_str"`
	TemplateType  string `json:"template_type"`
	SaName        string `json:"sa_name"`
	SaDpName      string `json:"sa_dp_name"`
	Username      string `json:"username"`
	UnionCode     string `json:"union_code"`
	Status        int    `json:"status"`
	CreatedAt     int64  `json:"created_at"`
	CreatedAtStr  string `json:"created_at_str"`
	RemindTime    int64  `json:"remind_time"`
	RemindTimeStr string `json:"remind_time_str"`
	FinishedAt    int64  `json:"finished_at"`
	FinishedAtStr string `json:"finished_at_str"`
	StatusStr     string `json:"status_str"`
	Content       string `json:"content"`
}

type SowTaskSimple struct {
	Id      int64  `json:"id"`
	Content string `json:"content"`
}
