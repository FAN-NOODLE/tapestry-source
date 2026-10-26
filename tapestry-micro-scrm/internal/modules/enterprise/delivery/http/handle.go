package http

import (
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"micro-scrm/internal/modules/enterprise"
	"micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/sdk/boss"
	"micro-scrm/internal/pkg/sdk/work"
	"time"
)

var service *enterprise.Service

// jsConfig 应用签名
func jsConfig(c *gin.Context) {
	request := enterprise.JsConfigRequest{}
	err := c.BindJSON(&request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	if request.Nonce == "" {
		str, err := helper.RandStr(16)
		if err != nil {
			ginplus.ResError(c, err)
			return
		}

		request.Nonce = str
	}

	if request.Timestamp == 0 {
		request.Timestamp = time.Now().Unix()
	}

	// 查询企业信息
	ctx := ctx.Wrap(c)
	_, corpId, _, err := service.GetInfo(ctx, request.EnterpriseID)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	// 请求企业微信生成企业ticket
	ticket, err := service.EnterpriseJsTicket(ctx, request.EnterpriseID)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	// 生成签名
	sign := work.JsSign(request.Nonce, ticket, request.Url, request.Timestamp)
	res := new(enterprise.JsConfigResponse)
	res.AppId = corpId
	res.Url = request.Url
	res.Timestamp = request.Timestamp
	res.Nonce = request.Nonce
	res.Sign = sign
	ginplus.ResSuccess(c, res)

}

// jsAgentConfig 应用签名
func jsAgentConfig(c *gin.Context) {
	request := enterprise.JsConfigRequest{}
	err := c.BindJSON(&request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	if request.Nonce == "" {
		str, err := helper.RandStr(16)
		if err != nil {
			ginplus.ResError(c, err)
			return
		}

		request.Nonce = str
	}

	if request.Timestamp == 0 {
		request.Timestamp = time.Now().Unix()
	}

	// 查询企业信息
	ctx := ctx.Wrap(c)
	_, corpId, _, err := service.GetInfo(ctx, request.EnterpriseID)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	// 请求企业微信生成企业ticket
	ticket, err := service.AgentJsTicket(ctx, request.EnterpriseID)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	// 获取自建应用信息
	req := boss.WorkPlugin{
		EnterpriseID: request.EnterpriseID,
	}
	pluginInfo, err := boss.WorkPluginInfo(ctx, req)
	if err != nil {
		return
	}

	// 生成签名
	sign := work.JsSign(request.Nonce, ticket, request.Url, request.Timestamp)
	res := new(enterprise.AgentConfigResponse)
	res.CorpId = corpId
	res.Url = request.Url
	res.Timestamp = request.Timestamp
	res.Nonce = request.Nonce
	res.Sign = sign
	res.AgentId = pluginInfo.Data.AgentID
	ginplus.ResSuccess(c, res)
}
