package smshttp

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"micro-scrm/internal/modules/sms"
)

var smsService *sms.Service

// createTemplate 创建短信模板
func createTemplate(c *gin.Context) {
	var createTemplateRequest sms.CreateTemplateRequest
	jsonError := c.BindJSON(&createTemplateRequest)

	if jsonError != nil {
		ginplus.ResError(c, jsonError, http.StatusOK)
		return
	}
	createTemplateRequest.TemplateContent = createTemplateRequest.TextContent
	res, err := smsService.CreateSmsTemplate(ctx.Wrap(c), createTemplateRequest)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, res)
}

// createSign 创建短信签名
func createSign(c *gin.Context) {
	var createSignRequest sms.SmsSignInfoRequest
	jsonError := c.BindJSON(&createSignRequest)

	if jsonError != nil {
		ginplus.ResError(c, jsonError, http.StatusOK)
		return
	}

	res, err := smsService.CreateSmsTemplateSign(ctx.Wrap(c), createSignRequest)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, res)
}

// arraignment 短信模板提审
func arraignment(c *gin.Context) {
	var smsTemplateArraignmentRequest sms.SmsTemplateArraignmentRequest
	jsonError := c.BindJSON(&smsTemplateArraignmentRequest)

	if jsonError != nil {
		ginplus.ResError(c, jsonError, http.StatusOK)
		return
	}

	res, err := smsService.SmsTemplateArraignment(ctx.Wrap(c), smsTemplateArraignmentRequest)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, res)
}

// paramsList 短信模板参数
func paramsList(c *gin.Context) {
	res, err := smsService.SmsTemplateParamsList(ctx.Wrap(c))

	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	ginplus.ResSuccess(c, res)
}

// templateSignList 获取签名列表
func templateSignList(c *gin.Context) {

	var templateListRequest sms.TemplateListRequest
	jsonError := c.BindJSON(&templateListRequest)
	if jsonError != nil {
		ginplus.ResError(c, jsonError)
		return
	}

	res, err := smsService.SmsTemplateSignList(ctx.Wrap(c), templateListRequest)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	ginplus.ResSuccess(c, res)
}

// templateList 获取短信模板列表
func templateList(c *gin.Context) {

	var templateListRequest sms.TemplateListRequest
	jsonError := c.BindJSON(&templateListRequest)
	if jsonError != nil {
		ginplus.ResError(c, jsonError)
		return
	}

	res, err := smsService.SmsTemplateList(ctx.Wrap(c), templateListRequest)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	ginplus.ResSuccess(c, res)
}

// templateDel 短信删除
func templateDel(c *gin.Context) {

	var updataSmsTemplateStatus sms.UpdateSmsTemplateStatus
	jsonError := c.BindJSON(&updataSmsTemplateStatus)

	if jsonError != nil {
		ginplus.ResError(c, jsonError, http.StatusOK)
		return
	}

	res, err := smsService.SmsTemplateDelete(ctx.Wrap(c), updataSmsTemplateStatus)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, res)
}

// templateEdit 编辑短信模板
func templateEdit(c *gin.Context) {

	var updateSmsTemplateStatus sms.UpdateSmsTemplateStatus
	jsonError := c.BindJSON(&updateSmsTemplateStatus)

	if jsonError != nil {
		ginplus.ResError(c, jsonError, http.StatusOK)
		return
	}
	updateSmsTemplateStatus.TemplateContent = updateSmsTemplateStatus.TextContent
	res, err := smsService.SmsTemplateEdit(ctx.Wrap(c), updateSmsTemplateStatus)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, res)
}

// templateInfo 短信模板详情
func templateInfo(c *gin.Context) {
	var smsTemplate sms.SmsTemplate
	jsonError := c.BindJSON(&smsTemplate)

	if jsonError != nil {
		ginplus.ResError(c, jsonError, http.StatusOK)
		return
	}

	res, err := smsService.SmsTemplateInfo(ctx.Wrap(c), smsTemplate)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, res)
}

// templateInfo 短信模板详情
func templateSignDelete(c *gin.Context) {
	var smsTemplateSignRequest sms.SmsSignInfoRequest
	jsonError := c.BindJSON(&smsTemplateSignRequest)

	if jsonError != nil {
		ginplus.ResError(c, jsonError, http.StatusOK)
		return
	}

	err := smsService.SmsTemplateSignDelete(ctx.Wrap(c), smsTemplateSignRequest)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, nil)
}

// synSmsStatus 同步短信模板标签状态
func synSmsStatus(c *gin.Context) {

	/*err := smsService.SynTemplateStatus(ctx.Wrap(c))

	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	*/

	ginplus.ResSuccess(c, nil)
}

// send 发送短信模板
func send(c *gin.Context) {
	var sendRequest sms.SendRequest
	jsonError := c.BindJSON(&sendRequest)

	if jsonError != nil {
		ginplus.ResError(c, jsonError, http.StatusOK)
		return
	}

	err := smsService.SmsSend(ctx.Wrap(c), sendRequest)

	if err != nil {
		ginplus.ResError(c, errorx.NewErrorMsg(700038, "短信平台服务器错误或短信次数已达上线"))
		return
	}
	ginplus.ResSuccess(c, nil)
}

// total 企业短信数据统计
func total(c *gin.Context) {
	var request sms.Request
	requestError := c.BindJSON(&request)

	if requestError != nil {
		ginplus.ResError(c, requestError, http.StatusOK)
		return
	}

	templateTotal, err := smsService.SmsTemplateTotal(ctx.Wrap(c), request)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, templateTotal)
}
