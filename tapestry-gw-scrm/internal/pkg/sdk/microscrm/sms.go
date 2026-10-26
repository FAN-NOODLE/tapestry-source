package microscrm

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/pkg/sdk"
)

// TemplateList 短信模板列表
func TemplateList(c ctx.Context, request interface{}) (interface{}, error) {
	var response CommonResponse
	err := scrmClient.Post(c, GetUrl(urlEnum.TemplateListUrl), request, &response)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
	//response, err := sdk.SendPost(c, GetUrl(urlEnum.TemplateListUrl), request)
	//return response, err
}

// CreateSign 创建短信签名
func CreateSign(c ctx.Context, request interface{}) (interface{}, error) {
	response, err := sdk.SendPost(c, GetUrl(urlEnum.CreateSignUrl), request)
	return response, err
}

// CreateTemplate 创建短信模板
func CreateTemplate(c ctx.Context, request interface{}) (interface{}, error) {
	response, err := sdk.SendPost(c, GetUrl(urlEnum.TemplateCreateUrl), request)
	return response, err
}

// arraignment 短信模板提审
func Arraignment(c ctx.Context, request interface{}) (interface{}, error) {
	var response CommonResponse
	err := scrmClient.Post(c, GetUrl(urlEnum.TemplateArraignmentUrl), request, &response)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

// ParamsList 短信模板参数
func ParamsList(c ctx.Context, request interface{}) (interface{}, error) {
	response, err := sdk.SendPost(c, GetUrl(urlEnum.ParamsListUrl), request)
	return response, err
}

// SignList 短信签名列表
func SignList(c ctx.Context, request interface{}) (interface{}, error) {
	response, err := sdk.SendPost(c, GetUrl(urlEnum.SignListUrl), request)
	return response, err
}

// TemplateDelete 短信签名列表
func TemplateDelete(c ctx.Context, request interface{}) (interface{}, error) {
	response, err := sdk.SendPost(c, GetUrl(urlEnum.TemplateDeleteUrl), request)
	return response, err
}

// TemplateDelete 短信签名列表
func TemplateEdit(c ctx.Context, request interface{}) (interface{}, error) {
	response, err := sdk.SendPost(c, GetUrl(urlEnum.TemplateEditUrl), request)
	return response, err
}

// TemplateDelete 短信签名列表
func TemplateInfo(c ctx.Context, request interface{}) (interface{}, error) {
	response, err := sdk.SendPost(c, GetUrl(urlEnum.TemplateInfoUrl), request)
	return response, err
}

// TemplateSignDelete 短信签名删除短信模板
func TemplateSignDelete(c ctx.Context, request interface{}) (interface{}, error) {
	response, err := sdk.SendPost(c, GetUrl(urlEnum.TemplateSignDeleteUrl), request)
	return response, err
}

// Send 短信发送
func Send(c ctx.Context, request interface{}) (interface{}, error) {
	response, err := sdk.SendPost(c, GetUrl(urlEnum.TemplateSendUrl), request)
	return response, err
}

// Total 企业短信统计
func Total(c ctx.Context, request interface{}) (interface{}, error) {
	response, err := sdk.SendPost(c, GetUrl(urlEnum.TemplateTotalUrl), request)
	return response, err
}
