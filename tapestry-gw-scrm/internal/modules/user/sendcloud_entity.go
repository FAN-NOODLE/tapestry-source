package user

// SendCloudApiUser apiUser
const SendCloudApiUser = "crs_notice"

// SendCloudApiKey apiKey
const SendCloudApiKey = "REDACTED_FOR_SOURCE_SHARING"

// SendCloudFrom from
const SendCloudFrom = "changxun@chufa.vchangyi.com"

// SendCloudFromName from name
const SendCloudFromName = "畅鲟"

// RegisterEnterpriseNoticeTemplate 邮件模板
const RegisterEnterpriseNoticeTemplate = "crs_register_enterprise_notice"

// SendTemplateMailRequest 发送邮件的请求
type SendTemplateMailRequest struct {
	ApiUser            string `json:"apiUser"`
	ApiKey             string `json:"apiKey"`
	From               string `json:"from"`
	To                 string `json:"to" desc:"地址列表. 在 useAddressList=true 时使用"`
	XSmtpApi           string `json:"xsmtpapi"`
	Subject            string `json:"subject"`
	TemplateInvokeName string `json:"templateInvokeName"`
	ContentSummary     string `json:"contentSummary"`
	FromName           string `json:"fromName" desc:"发件人显示名称"`
	ReplyTo            string `json:"replyTo"`
	LabelId            int    `json:"labelId"`
	Headers            string `json:"headers"`
	UseAddressList     string `json:"use_address_list" desc:"默认是false"`
}

// XSmtpApiRequest
type XSmtpApiRequest struct {
	To  []string            `json:"to"`
	Sub map[string][]string `json:"sub"`
}

// SendCloudResponse 返回结果
type SendCloudResponse struct {
	Result     bool   `json:"result"`
	StatusCode int64  `json:"statusCode"`
	Message    string `json:"message"`
}
