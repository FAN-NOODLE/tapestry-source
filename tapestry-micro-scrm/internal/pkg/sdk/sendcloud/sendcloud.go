package sendcloud

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
)

type CreateSmsSignResponse struct {
	CommonResponse
	Info struct {
		SmsSignVo struct {
			Id           int64  `json:"id"`
			SignName     string `json:"signName"`
			SignTypeDesc string `json:"signTypeDesc"`
			VerifyStatus string `json:"verifyStatus"`
			CreateTime   string `json:"createTime"`
			UpdateTime   string `json:"updateTime"`
		} `json:"smsSignVo"`
	} `json:"info"`
}

// 添加短信模板
type CreateSmsTemplateResponse struct {
	CommonResponse
	Info struct {
		TemplateId int64 `json:"templateId"`
	} `json:"info"`
}

// 短信模板提审
type SmsTemplateSubmitsmsResponse struct {
	CommonResponse
	Info struct {
	} `json:"info"`
}

// 短信模板更新
type SmsTemplateUpdateResponse struct {
	CommonResponse
	Info struct {
	} `json:"info"`
}

// 短信模板列表返回数据格式
type SmsTemplateListResponse struct {
	CommonResponse
	Info struct {
		Total int64 `json:"total"`
		List  []struct {
			TemplateContent    string `json:"templateContent"`
			TemplateCreateTime string `json:"templateCreateTime"`
			TemplateId         int64  `json:"templateId"`
			TemplateName       string `json:"templateName"`
			MsgType            string `json:"msgType"`
			SmsType            string `json:"smsType"`
			IsVerify           string `json:"isVerify"`
			TemplateUpdateTime string `json:"templateUpdateTime"`
		}
	} `json:"info"`
}

// 短信签名列表
type SmsTemplateSignListResponse struct {
	CommonResponse
	Info struct {
		VoList []struct {
			CreateTime   string `json:"createTime"`
			Id           int64  `json:"id"`
			SignName     string `json:"signName"`
			SignTypeDesc string `json:"signTypeDesc"`
			UpdateTime   string `json:"updateTime"`
			VerifyStatus string `json:"verifyStatus"`
		} `json:"voList"`
	} `json:"info"`
}

// 短信发送成功返回值
type SendResponse struct {
	CommonResponse
	Info struct {
		SuccessCount int64    `json:"successCount"`
		SmsIds       []string `json:"smsIds"`
	} `json:"info"`
}

// CreateSmsSign 创建短信签名
func CreateSmsSign(c ctx.Context, request interface{}) (*CreateSmsSignResponse, error) {
	var response CreateSmsSignResponse
	err := sendGet(c, GetURL(urlEnum.CreateSmsSign), request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// CreateSmsTemplate 创建短信模板
func CreateSmsTemplate(c ctx.Context, request interface{}) (*CreateSmsTemplateResponse, error) {
	var response CreateSmsTemplateResponse
	err := sendGet(c, GetURL(urlEnum.CreateSmsTemplate), request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// UpdateSmsTemplate 更新短信模板
func UpdateSmsTemplate(c ctx.Context, request interface{}) (*SmsTemplateUpdateResponse, error) {
	var response SmsTemplateUpdateResponse
	err := sendGet(c, GetURL(urlEnum.UpdateTemplate), request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// SmsTemplateSubmitsms 保存提交审核
func SmsTemplateSubmitsms(c ctx.Context, request interface{}) (*SmsTemplateSubmitsmsResponse, error) {
	var response SmsTemplateSubmitsmsResponse
	err := sendGet(c, GetURL(urlEnum.SubmitsmsSmsTemplate), request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// SmsTemplateList 短信模板列表
func SmsTemplateList(c ctx.Context, request interface{}) (*SmsTemplateListResponse, error) {
	var response SmsTemplateListResponse
	err := sendGet(c, GetURL(urlEnum.TemplateList), request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// SmsTemplateSignList 短信签名列表
func SmsTemplateSignList(c ctx.Context, request interface{}) (*SmsTemplateSignListResponse, error) {
	var response SmsTemplateSignListResponse
	err := sendGet(c, GetURL(urlEnum.TemplateSignList), request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// SmsTemplateSignList 短信签名列表
func Send(c ctx.Context, request interface{}) (*SendResponse, error) {
	var response SendResponse
	err := sendGet(c, GetURL(urlEnum.Send), request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func CheckSmsSignStatus(status string) int64 {
	var statusNum int64 = 0
	switch status {
	case "待审核":
		statusNum = 1
	case "审核成功":
		statusNum = 3
	case "审核失败":
		statusNum = 2
	default:
		statusNum = 2
	}
	return statusNum
}

func CheckSmsDesc(desc string) int64 {
	var descNum int64 = 0
	switch desc {
	case "国内短信":
		descNum = 0
	case "国际短信":
		descNum = 1
	default:
		descNum = 0
	}
	return descNum
}
