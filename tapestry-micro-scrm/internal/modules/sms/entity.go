package sms

const (
	// 加粉短信模板类型
	TemplateTypeIsFans = 1
	// 普通模板类型
	TemplateTypeNormal = 2
	// 畅移科技 签名 ID
	DefaultSignID = 9465
	// 默认的短信签名
	DefaultSignName = "畅移科技"
	// 默认短信签名在系统中的ID
	DefaultSignSysID = 117567463708480
	// 短信加粉模板名称
	DefaultName = "短信加粉"
	// 默认的短信模板内容
	DefaultTemplateContent = "亲爱的小主，您的订单%order_id%已受理。偷偷跟你讲喔~两小时内添加%member_name%（%member_mobile%），送%number%%classifier%%goods_name%。"
	// 默认短信模板内容
	DefaultTemplateContentStr = "<p>亲爱的小主，您的订单<input type=\"button\" class=\"status-orange\" attr=\"order_id\" value=\"%订单号%\">已受理。偷偷跟你讲喔~两小时内添加<input type=\"button\" class=\"status-blue\" attr=\"member_name\" value=\"%员工微信昵称%\">（<input type=\"button\" class=\"status-blue\" attr=\"member_mobile\" value=\"%员工手机号%\">），送<input type=\"button\" class=\"status-blue\" attr=\"number\" value=\"%数字%\"><input type=\"button\" class=\"status-blue\" attr=\"classifier\" value=\"%量词%\"><input type=\"button\" class=\"status-blue\" attr=\"goods_name\" value=\"%商品名称%\"></p>"
	// 默认短信模板参数
	DefaultTemplateParams = ""
)

const (
	// 异常模板名称
	AbnormalName = "异常策略提醒"
	// 异常模板内容
	AbnormalTemplateContent = "尊敬的用户，您在SCRM系统中的短信条数不足100条，为了您的正常使用，请联系商务张艳芳（%member_mobile%）充值短信。"
	// 异常模板内容（前端使用）
	AbnormalTemplateContentStr = "<p>尊敬的用户，您在SCRM系统中的短信条数不足100条，为了您的正常使用，请联系商务张艳芳（<input type=\"button\" class=\"status-blue\" attr=\"member_mobile\" value=\"%13186008012%\">）充值短信</p>"
	// 异常模板默认参数
	AbnormalTemplateParams = "{\"member_mobile\":\"13186008012\"}"
)

const (
	// 通知名称
	ReappearingName = "长期未使用短信被置为未提审策略通知"
	// 通知模板内容
	ReappearingTemplateContent = "尊敬的用户，您配置的营销策略因长期未触发，导致配置的短信模板被第三方下架，如仍需使用请再次提审。"
	// 通知模板内容（前端使用）
	ReappearingTemplateContentStr = "<p>尊敬的用户，您配置的营销策略因长期未触发，导致配置的短信模板被第三方下架，如仍需使用请再次提审。</p>"
	// 通知模板默认参数
	ReappearingTemplateParams = ""
)

const (
	// 通知名称
	OverdueName = "策略到期提醒"
	// 通知模板内容
	OverdueTemplateContent = "尊敬的用户，您有营销策略将在24小时内结束，请及时登录SCRM系统查看。"
	// 通知模板内容（前端使用）
	OverdueTemplateContentStr = "<p>尊敬的用户，您有营销策略将在24小时内结束，请及时登录SCRM系统查看。</p>"
	// 通知模板默认参数
	OverdueTemplateParams = ""
)

const (
	// 通知名称
	CouponName = "优惠券不足100条"
	// 通知模板内容
	CouponTemplateContent = "尊敬的用户，您有自动发券的营销策略中优惠券条数不足100条，请及时查看。退订回T"
	// 通知模板内容（前端使用）
	CouponTemplateContentStr = "<p>尊敬的用户，您有自动发券的营销策略中优惠券条数不足100条，请及时查看。</p>"
	// 通知模板默认参数
	CouponTemplateParams = ""
)

// sms_sign 短信模板签名状态
const (
	// 审核中
	ThirdSignStatusUnderReview = 1 + iota
	// 审核失败
	ThirdSignStatusFailure
	// 审核成功
	ThirdSignStatusSuccessful
)

// sms_template表 短信模板状态
const (
	// 待审核
	TemplateStatusWaitVerify = 1 + iota
	// 审核中
	TemplateStatusUnderReview
	// 审核失败
	TemplateStatusAuditFailure
	// 审核成功
	TemplateStatusSuccessfulReview
	// 已赋值
	TemplateStatusHasValue
)

// 每次查询条数
var LimitNumber int64 = 100

// 参数类型
const (
	ParamsTypeDynamic = 1 + iota
	ParamsTypeStatic
)

// 可用短信条数
var AvailableMessageNumberKey = "availableMessageNumber"

// 企业剩余提醒数量
var MessageNoticeNumber int64 = 100

// CreateTemplateRequest 创建短信模板
type CreateTemplateRequest struct {
	EnterpriseId    int64  `json:"enterprise_id" desc:"企业Id"`
	AdminId         int64  `json:"admin_id" desc:"管理员Id"`
	TemplateName    string `json:"template_name" desc:"短信模板名称"`
	TemplateSignId  int64  `json:"template_sign_id" desc:"短信签名Id"`
	MessageTemplate string `json:"message_template" desc:"短信模板"`
	TemplateHtml    string `json:"template_html" desc:"前端需要的html" db:"template_html"`
	TemplateContent string `json:"template_content" db:"template_content"`
	TextContent     string `json:"text_content"`
	MessageStr      string `json:"message_str" desc:"前端列表展示字符串" db:"message_str"`
	IsExamine       int64  `json:"is_examine" desc:"是否提交审核"`
}

// SmsTemplate 短信模板状态码
type SmsTemplate struct {
	Id              int64  `json:"id" desc:"短信模板ID" db:"id"`
	EnterpriseId    int64  `json:"enterprise_id" desc:"企业ID" db:"enterprise_id"`
	AdminId         int64  `json:"admin_id" desc:"操作人ID" db:"admin_id"`
	TemplateId      int64  `json:"template_id" desc:"第三方短信模板ID" db:"template_id"`
	TemplateName    string `json:"template_name" desc:"模板名称" db:"template_name"`
	TemplateSignId  int64  `json:"template_sign_id" desc:"模板签名Id" db:"template_sign_id"`
	MessageTemplate string `json:"message_template" desc:"短信模板" db:"message_template"`
	TemplateContent string `json:"template_content" desc:"短信插入参数后内容体" db:"template_content"`
	TemplateHtml    string `json:"template_html" desc:"前端需要的html" db:"template_html"`
	MessageStr      string `json:"message_str" desc:"前端列表展示字符串" db:"message_str"`
	Status          int64  `json:"status" desc:"审核状态（1：未审核、2：审核中、3：审核失败、4：审核成功、5：已赋值）" db:"status"`
	Type            int64  `json:"type" desc:"模板类型（1：默认模板，2：管理员自建模板）" db:"type"`
	Occupancy       int64  `json:"occupancy"`
	Fail_reason     string `json:"fail_reason" desc:"失败原因" db:"fail_reason"`
	CreatedAt       int64  `json:"created_at" desc:"创建时间" db:"created_at"`
	UpdatedAt       int64  `json:"updated_at" desc:"更新时间" db:"updated_at"`
	DeletedAt       int64  `json:"deleted_at" desc:"删除时间" db:"deleted_at"`

	TemplateShortName string  `json:"template_short_name" desc:"短信模板展示模板名称"`
	TemplateIDs       []int64 `json:"template_ids" desc:"模版IDs"`
	StrategyID        int64   `json:"strategy_id" desc:"策略ID" db:"strategy_id"`
	StrategyIDs       []int64 `json:"strategy_ids" desc:"策略ID" `
}

// SendLog 短信日志表结构体
type SendLog struct {
	Id             int64 `json:"id" desc:"短信日志ID" db:"id"`
	EnterpriseId   int64 `json:"enterprise_id" desc:"企业ID" db:"enterprise_id"`
	UserId         int64 `json:"user_id" desc:"用户ID" db:"user_id"`
	SmsTemplateId  int64 `json:"sms_template_id" desc:"短信模板主键" db:"sms_template_id"`
	ChannelType    int64 `json:"channel_type" desc:"短信日志ID" db:"channel_type"`
	ChannelValueId int64 `json:"channel_value_id" desc:"短信触发渠道（1：营销策略）" db:"channel_value_id"`
	Status         int64 `json:"status" desc:"短信发送状态（1：成功 2：失败）" db:"status"`
	CreatedAt      int64 `json:"created_at" desc:"数据创建时间" db:"created_at"`
	UpdatedAt      int64 `json:"updated_at" desc:"数据更新时间" db:"updated_at"`
	DeletedAt      int64 `json:"deleted_at" desc:"数据删除时间" db:"deleted_at"`
}

// SmsParams 短信模板状态码
type SmsParams struct {
	Id           int64  `json:"id" desc:"短信模板ID" db:"id"`
	EnterpriseId int64  `json:"enterprise_id" desc:"企业ID" db:"enterprise_id"`
	AdminId      int64  `json:"admin_id" desc:"操作人ID" db:"admin_id"`
	Key          string `json:"key" desc:"第三方短信模板ID" db:"key"`
	Name         string `json:"name" desc:"模板名称" db:"name"`
	ParamsType   int64  `json:"params_type" desc:"参数类型（1：动态参数，2：静态参数）" db:"params_type"`
	CreatedAt    int64  `json:"created_at" desc:"创建时间" db:"created_at"`
	UpdatedAt    int64  `json:"updated_at" desc:"更新时间" db:"updated_at"`
	DeletedAt    int64  `json:"deleted_at" desc:"删除时间" db:"deleted_at"`
}

// SmsParamsRequest 列表请求数据
type SmsParamsRequest struct {
	List []SmsParams `json:"list"`
}

// SmsParamsRequest 列表请求数据
type SmsSignListRequest struct {
	List       []SmsSignInfo `json:"list"`
	Page       int64         `json:"page"`
	Size       int64         `json:"size"`
	Count      int64         `json:"total"`
	UpdateTime int64         `json:"updateTime"`
}

// SmsParamsRequest 列表请求数据
type SmsTemplateListRequest struct {
	List       []SmsTemplate `json:"list"`
	Page       int64         `json:"page"`
	Size       int64         `json:"size"`
	Count      int64         `json:"total"`
	TpdateTime int64         `json:"updateTime"`
}

// UpdataSmsTemplateStatus 修改短信模板状态
type UpdateSmsTemplateStatus struct {
	Id              int64  `json:"id" desc:"短信模板ID" db:"id"`
	Status          int64  `json:"status" desc:"审核状态（1：未审核、2：审核中、3：审核失败、4：审核成功、5：已赋值）" db:"status"`
	IsExamine       int64  `json:"is_examine"`
	EnterpriseId    int64  `json:"enterprise_id" db:"enterprise_id"`
	AdminId         int64  `json:"admin_id" db:"admin_id"`
	TemplateId      int64  `json:"template_id" db:"template_id"`
	TemplateName    string `json:"template_name" db:"template_name"`
	TemplateSignId  int64  `json:"template_sign_id" db:"template_sign_id"`
	MessageTemplate string `json:"message_template" db:"message_template"`
	TemplateContent string `json:"template_content" db:"template_content"`
	TemplateHtml    string `json:"template_html" db:"template_html"`
	TextContent     string `json:"text_content"`
	MessageStr      string `json:"message_str" desc:"前端列表展示字符串" db:"message_str"`
	FailReason      string `json:"fail_reason" db:"fail_reason"`
	CreatedAt       int64  `json:"created_at" db:"created_at"`
	UpdatedAt       int64  `json:"updated_at" desc:"更新时间" db:"updated_at"`
	DeletedAt       int64  `json:"deleted_at" desc:"更新时间" db:"deleted_at"`
}

// SmsTemplateArraignmentReponse 提审成功返回短信模板Id
type SmsTemplateArraignmentReponse struct {
	Id int64 `json:"id" desc:"短信模板ID" db:"id"`
}

// SmsTemplate 短信模板状态码
type SmsTemplateArraignmentRequest struct {
	Id int64 `json:"id" desc:"短信模板ID" db:"id"`
}

// SmsSignInfoRequest 获取短信签名
type SmsSignInfoRequest struct {
	Id              int64  `json:"id" db:"id"`
	EnterpriseId    int64  `json:"enterprise_id" db:"enterprise_id"`
	AdminId         int64  `json:"admin_id" db:"admin_id"`
	ThirdSignId     int64  `json:"third_sign_id" db:"third_sign_id"`
	ThirdSignStatus int64  `json:"third_sign_status" db:"third_sign_status"`
	ThirdSignType   int64  `json:"third_sign_type" db:"third_sign_type"`
	SignName        string `json:"sign_name" db:"sign_name"`
	UpdatedAt       int64  `json:"updated_at" db:"updated_at"`
	DeletedAt       int64  `json:"deleted_at" db:"deleted_at"`
}

// SmsSignInfoResponse 短信签名返回结构体
type SmsSignInfo struct {
	Id              int64  `json:"id" db:"id"`
	EnterpriseId    int64  `json:"enterprise_id" db:"enterprise_id"`
	AdminId         int64  `json:"admin_id" db:"admin_id"`
	ThirdSignId     int64  `json:"third_sign_id" db:"third_sign_id"`
	ThirdSignStatus int64  `json:"third_sign_status" db:"third_sign_status"`
	ThirdSignType   int64  `json:"third_sign_type" db:"third_sign_type"`
	SignName        string `json:"sign_name" db:"sign_name"`
	CreatedAt       int64  `json:"created_at" db:"created_at"`
	UpdatedAt       int64  `json:"updated_at" db:"updated_at"`
	DeletedAt       int64  `json:"deleted_at" db:"deleted_at"`
}

// CreateTagResponse 创建标签返回参数
type CreateTemplateResponse struct {
	TemplateId int64 `json:"template_id"`
}

// EnterpriseRequest 通过企业信息查询条件
type EnterpriseRequest struct {
	EnterpriseId int64 `json:"enterprise_id" desc:"企业Id" db:"enterprise_id"`
}

type TemplateListRequest struct {
	EnterpriseId  int64   `json:"enterprise_id" db:"enterprise_id"`
	Page          int64   `json:"page" db:"page"`
	Size          int64   `json:"size" db:"size"`
	SearchType    int64   `json:"search_type" db:"search_type"`
	SearchText    string  `json:"search_text" db:"search_text"`
	Status        int64   `json:"status" db:"third_sign_status"`
	TemplateIds   []int64 `json:"ids"`
	TemplateIdCnt int     `json:"template_id_cnt"`
	Type          int64   `json:"type"`
}

// SendRequest 发送短信请求参数
type SendRequest struct {
	UserId         int64  `json:"user_id"`
	EnterpriseId   int64  `json:"enterprise_id"`
	TemplateId     int64  `json:"template_id"`
	Phone          string `json:"phone"`
	ChannelType    int64  `json:"channel_type"`
	ChannelValueId int64  `json:"channel_value_id"`
	//Vars           struct {
	//	UserName    string `json:"user_name"`
	//	AddressName string `json:"address_name"`
	//	OrderId     string `json:"order_id"`
	//} `json:"vars"`
	Vars VarsInfo `json:"vars"`
}

type VarsInfo struct {
	UserName string `json:"user_name"`
	UserNick string `json:"user_nick"`
}

// TemplateTotalRequest 获取模板条数返回参数
type TemplateTotalRequest struct {
	AvailableNum int64 `json:"available_num"`
	UsedNum      int64 `json:"used_num"`
}

type Request struct {
	EnterpriseId int64 `json:"enterprise_id" db:"enterprise_id"`
	UserId       int64 `json:"user_id"`
}

// 总数返回值
type CountResponse struct {
	Count int64 `json:"type" db:"count"`
}
