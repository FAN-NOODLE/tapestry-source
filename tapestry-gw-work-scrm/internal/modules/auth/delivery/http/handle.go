package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"gw-work-scrm/internal/modules/auth"
	"gw-work-scrm/internal/pkg/common"
	"gw-work-scrm/internal/pkg/errs"
	"gw-work-scrm/internal/pkg/jwt"
	microScrm "gw-work-scrm/internal/pkg/sdk/microscrm"
	"net/url"
	"strconv"
)

// QyRedirectUrl 重定向地址
const QyRedirectUrl = "https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_base&state=STATE#wechat_redirect"

// memberSdk 请求人员SDK
var memberSdk *microScrm.MemberApi

// enterpriseSdk 企业SDK
var enterpriseSdk *microScrm.EnterpriseApi

// authLogin 用户授权登录
func authLogin(c *gin.Context) {
	// 获取参数并绑定
	request := auth.LoginAuthRequest{}
	if err := c.ShouldBind(&request); err != nil {
		ginplus.ResError(c, err)
		return
	}
	ct := ctx.Wrap(c)

	// 查询企业
	enterprise, err := common.GetEnterpriseInfo(ct, request.Hash)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	request.EnterpriseID = enterprise.EnterpriseID

	// 根据code去请求企业微信接口进行获取用户信息
	infoRequest := microScrm.AuthInfoRequest{
		Code:         request.Code,
		EnterpriseID: request.EnterpriseID,
		LoginType:    request.LoginType,
	}

	authInfo, err := memberSdk.AuthInfo(ct, infoRequest)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	if authInfo == nil {
		ginplus.ResError(c, errs.AuthInfoNotFind)
		return
	}

	// 生成 token
	token, err := jwt.New().CreateToken(
		jwt.NewCustomClaims(
			ct,
			authInfo.Member.UserId,
			authInfo.External.OpenId,
			authInfo.External.ExternalUserid,
			c.FullPath(),
			request.EnterpriseID,
		),
	)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	ginplus.ResSuccess(c, auth.LoginAuthResponse{AuthInfoResponse: authInfo, Token: token, BrandCode: enterprise.BrandCode})
}

// redirectUrl 返回授权的地址
func redirectUrl(c *gin.Context) {
	// 获取参数并绑定
	request := auth.RedirectUrlRequest{}
	if err := c.BindJSON(&request); err != nil {
		ginplus.ResError(c, err)
		return
	}

	enterprise, err := common.GetEnterpriseInfo(ctx.Wrap(c), request.Hash)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	// 构建授权地址
	baseUrl := url.QueryEscape(request.RedirectUri)
	redirectUrl := fmt.Sprintf(QyRedirectUrl, enterprise.CorpID, baseUrl)
	response := &auth.RedirectUrlResponse{Url: redirectUrl}
	ginplus.ResSuccess(c, response)
}

// jsConfig 企业签名
func jsConfig(c *gin.Context) {
	request := microScrm.JsConfigRequest{}
	if err := c.BindJSON(&request); err != nil {
		ginplus.ResError(c, err)
		return
	}

	enterprise, err := common.GetEnterpriseInfo(ctx.Wrap(c), request.Hash)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	request.EnterpriseID = enterprise.EnterpriseID

	ticket, err := enterpriseSdk.JsConfig(ctx.Wrap(c), request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	ginplus.ResSuccess(c, ticket)
}

// jsAgentConfig 应用签名
func jsAgentConfig(c *gin.Context) {
	request := microScrm.JsConfigRequest{}
	if err := c.BindJSON(&request); err != nil {
		ginplus.ResError(c, err)
		return
	}

	enterprise, err := common.GetEnterpriseInfo(ctx.Wrap(c), request.Hash)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	request.EnterpriseID = enterprise.EnterpriseID

	ticket, err := enterpriseSdk.JsAgentConfig(ctx.Wrap(c), request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	ginplus.ResSuccess(c, ticket)
}

// mock 模型登录地址
func mock(c *gin.Context) {
	userId, ok := c.GetQuery("userid")
	if !ok {
		ginplus.ResError(c, errorx.NewErrorMsg(10086, "userid必传"))
		return
	}

	enterpriseId, ok := c.GetQuery("enterprise_id")
	if !ok {
		ginplus.ResError(c, errorx.NewErrorMsg(10086, "enterpriseId必传"))
		return
	}

	enterprise, err := strconv.ParseInt(enterpriseId, 10, 64)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	// 生成 token
	token, err := jwt.New().CreateToken(
		jwt.NewCustomClaims(
			ctx.Wrap(c),
			userId,
			"",
			"",
			c.FullPath(),
			enterprise,
		),
	)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	ginplus.ResSuccess(c, gin.H{
		"token":  token,
		"userid": userId,
	})
}
