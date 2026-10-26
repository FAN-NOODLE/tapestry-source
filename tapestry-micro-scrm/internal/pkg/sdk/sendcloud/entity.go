package sendcloud

// Response 返回值接口
type Response interface {
	ResponseCheck
}

// ResponseCheck 返回值校验接口
type ResponseCheck interface {
	CheckCode() error
}

// CreateSmsSignRequest 添加短信签名请求结构体
type CreateSmsSignRequest struct {
	SignType  int    `json:"signType"`
	SignName  string `json:"signName"`
	SmsUser   string `json:"smsUser"`
	Signature string `json:"signature"`
}

// CreateSmsSignRequest 查询短信短信列表请求结构体
type TemplateListRequest struct {
	SmsUser       string `json:"smsUser"`
	TemplateIdStr string `json:"templateIdStr"`
	Start         int64  `json:"start"`
	Limit         int64  `json:"limit"`
	Signature     string `json:"signature"`
}

// CreateSmsSignRequest 查询短信短信列表请求结构体
type SendRequest struct {
	SmsUser    string `json:"smsUser"`
	TemplateId int64  `json:"templateId"`
	Phone      string `json:"phone"`
	Vars       string `json:"vars"`
	Signature  string `json:"signature"`
}

type SendMsgRequest struct {
	Mobile       string `json:"mobile"`
	Template     int64  `json:"template"`
	EnterpriseId int64  `json:"enterpriseId"`
	Content      string `json:"content"`
}

// TemplateSignListRequest 查询短信签名列表请求结构体
type TemplateSignListRequest struct {
	SmsUser   string `json:"smsUser"`
	Signature string `json:"signature"`
}

// CreateSmsSignRequest 短信审核请求结构体
type SmsTemplateSubmitsmsRequest struct {
	SmsUser       string `json:"smsUser"`
	TemplateIdStr string `json:"templateIdStr"`
	Signature     string `json:"signature"`
}

// CreateSmsTemplatequest 添加短信模板
type CreateSmsTemplatequest struct {
	SmsUser         string `json:"smsUser"`
	TemplateName    string `json:"templateName"`
	TemplateText    string `json:"templateText"`
	SignName        string `json:"signName"`
	SignId          int64  `json:"signId"`
	SignPositionStr string `json:"signPositionStr"`
	SmsTypeStr      string `json:"smsTypeStr"`
	Signature       string `json:"signature"`
}

// UpdateSmsTemplatequest 编辑短信模板
type UpdateSmsTemplatequest struct {
	SmsUser         string `json:"smsUser"`
	TemplateIdStr   string `json:"templateIdStr"`
	TemplateName    string `json:"templateName"`
	TemplateText    string `json:"templateText"`
	SignName        string `json:"signName"`
	SignId          int64  `json:"signId"`
	SignPositionStr string `json:"signPositionStr"`
	SmsTypeStr      string `json:"smsTypeStr"`
	Signature       string `json:"signature"`
}
