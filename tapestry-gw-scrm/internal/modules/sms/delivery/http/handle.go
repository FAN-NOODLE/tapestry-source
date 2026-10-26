package sms_http

import (
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"gw-scrm/internal/pkg/common"
	"gw-scrm/internal/pkg/sdk/microscrm"
)

// templateList 短信模板列表
func templateList(c *gin.Context) {

	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}

	var request map[string]interface{}
	c.BindJSON(&request)
	request["enterprise_id"] = tokenInfo.EnterpriseId
	request["admin_id"] = tokenInfo.AdminId

	response, err := microscrm.TemplateList(ctx.Wrap(c), request)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)
}

// templateList 短信模板列表
func createSign(c *gin.Context) {
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}

	var request map[string]interface{}
	c.BindJSON(&request)
	request["enterprise_id"] = tokenInfo.EnterpriseId
	request["admin_id"] = tokenInfo.AdminId
	response, err := microscrm.CreateSign(ctx.Wrap(c), request)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)
}

// templateList 短信模板列表
func createTemplate(c *gin.Context) {
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}

	var request map[string]interface{}
	c.BindJSON(&request)
	request["enterprise_id"] = tokenInfo.EnterpriseId
	request["admin_id"] = tokenInfo.AdminId
	response, err := microscrm.CreateTemplate(ctx.Wrap(c), request)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)
}

// arraignment 短信模板提审
func arraignment(c *gin.Context) {
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}

	var request map[string]interface{}
	c.BindJSON(&request)
	request["enterprise_id"] = tokenInfo.EnterpriseId
	request["admin_id"] = tokenInfo.AdminId
	response, err := microscrm.Arraignment(ctx.Wrap(c), request)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)
}

// paramsList 短信参数列表
func paramsList(c *gin.Context) {
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}

	var request map[string]interface{}
	c.BindJSON(&request)
	request["enterprise_id"] = tokenInfo.EnterpriseId
	request["admin_id"] = tokenInfo.AdminId
	response, err := microscrm.ParamsList(ctx.Wrap(c), request)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)
}

// paramsList 短信参数列表
func signList(c *gin.Context) {
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}

	var request map[string]interface{}
	c.BindJSON(&request)
	request["enterprise_id"] = tokenInfo.EnterpriseId
	request["admin_id"] = tokenInfo.AdminId
	response, err := microscrm.SignList(ctx.Wrap(c), request)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)
}

// paramsList 短信参数列表
func templateDelete(c *gin.Context) {
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}

	var request map[string]interface{}
	c.BindJSON(&request)
	request["enterprise_id"] = tokenInfo.EnterpriseId
	request["admin_id"] = tokenInfo.AdminId
	response, err := microscrm.TemplateDelete(ctx.Wrap(c), request)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)
}

// templateEdit 删除短信模板
func templateEdit(c *gin.Context) {
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}

	var request map[string]interface{}
	c.BindJSON(&request)
	request["enterprise_id"] = tokenInfo.EnterpriseId
	request["admin_id"] = tokenInfo.AdminId
	response, err := microscrm.TemplateEdit(ctx.Wrap(c), request)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)
}

// templateEdit 删除短信模板
func templateInfo(c *gin.Context) {
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}

	var request map[string]interface{}
	c.BindJSON(&request)
	request["enterprise_id"] = tokenInfo.EnterpriseId
	request["admin_id"] = tokenInfo.AdminId
	response, err := microscrm.TemplateInfo(ctx.Wrap(c), request)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)
}

// templateSignDelete 删除短信签名
func templateSignDelete(c *gin.Context) {
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}

	var request map[string]interface{}
	c.BindJSON(&request)
	request["enterprise_id"] = tokenInfo.EnterpriseId
	request["admin_id"] = tokenInfo.AdminId
	response, err := microscrm.TemplateSignDelete(ctx.Wrap(c), request)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)
}

// send 删除短信签名
func send(c *gin.Context) {
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}

	var request map[string]interface{}
	c.BindJSON(&request)
	request["enterprise_id"] = tokenInfo.EnterpriseId
	request["admin_id"] = tokenInfo.AdminId
	response, err := microscrm.Send(ctx.Wrap(c), request)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)
}

// total 企业短信统计
func total(c *gin.Context) {
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}

	var request map[string]interface{}
	c.BindJSON(&request)
	request["enterprise_id"] = tokenInfo.EnterpriseId
	request["admin_id"] = tokenInfo.AdminId
	response, err := microscrm.Total(ctx.Wrap(c), request)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)
}
