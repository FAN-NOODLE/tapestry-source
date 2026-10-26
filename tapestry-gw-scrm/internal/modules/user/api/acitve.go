package api

import (
	"encoding/json"
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/clients"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/helper"
	"gw-scrm/internal/modules/user"
	"gw-scrm/internal/pkg/common"
	"time"
)

//var client = clients.NewHTTPClient()
//var client = common.NewHttpClient()
var client = common.NewHttpClient(clients.Timeout(59 * time.Second))

type ActiveApiImpl struct{}

// 手机号验证
func (a ActiveApiImpl) CheckMobile(req user.GetCodeRequest) (interface{}, error) {
	ClientMode := &common.ClientMode{}
	ClientUrl := ClientMode.GetClientUrl(common.GetClientUrl{"micro-authority"})
	url := ClientUrl + "/micro-authority/v1/account/check-account"
	res, err := client.Post(ctx.New(), url, clients.JSONBody(req))
	helper.Must(err)
	return res.Read(func(bytes []byte) (interface{}, error) {
		rt := fmt.Sprint(string(bytes))
		return rt, nil
	})
}

// 获取验证码
func (a ActiveApiImpl) GetCode(req user.GetCodeRequest) (interface{}, error) {
	ClientMode := &common.ClientMode{}
	ClientUrl := ClientMode.GetClientUrl(common.GetClientUrl{"micro-authority"})
	url := ClientUrl + "/micro-authority/v1/common/send-code"
	res, err := client.Post(ctx.New(), url, clients.JSONBody(req))
	helper.Must(err)
	return res.Read(func(bytes []byte) (interface{}, error) {
		rt := fmt.Sprint(string(bytes))
		return rt, nil
	})
}

// 绑定手机号验证
func (a ActiveApiImpl) CheckBindMobile(req user.GetBindCodeRequest) (interface{}, error) {
	ClientMode := &common.ClientMode{}
	ClientUrl := ClientMode.GetClientUrl(common.GetClientUrl{"boss"})
	url := ClientUrl + "/service/v1/enterprise_admin/detail"
	res, err := client.Post(ctx.New(), url, clients.JSONBody(req))
	helper.Must(err)
	return res.Read(func(bytes []byte) (interface{}, error) {
		rt := fmt.Sprint(string(bytes))
		return rt, nil
	})
}

// 获取绑定验证码
func (a ActiveApiImpl) GetBindCode(req user.GetBindCodeRequest) (interface{}, error) {
	ClientMode := &common.ClientMode{}
	ClientUrl := ClientMode.GetClientUrl(common.GetClientUrl{"micro-authority"})
	url := ClientUrl + "/micro-authority/v1/common/send-sms-code"
	res, err := client.Post(ctx.New(), url, clients.JSONBody(req))
	helper.Must(err)
	return res.Read(func(bytes []byte) (interface{}, error) {
		rt := fmt.Sprint(string(bytes))
		return rt, nil
	})
}

// 登录
func (a ActiveApiImpl) Login(req user.LoginRequest) (interface{}, error) {
	ClientMode := &common.ClientMode{}
	ClientUrl := ClientMode.GetClientUrl(common.GetClientUrl{"micro-authority"})
	url := ClientUrl + "/micro-authority/v1/account/scrm-password-login"
	if req.LoginType == 2 {
		url = ClientUrl + "/micro-authority/v1/account/scrm-code-login"
	}
	res, err := client.Post(ctx.New(), url, clients.JSONBody(req))
	helper.Must(err)
	return res.Read(func(bytes []byte) (interface{}, error) {
		rt := fmt.Sprint(string(bytes))
		return rt, nil
	})
}

// 登录
func (a ActiveApiImpl) ConsoleIpWhiteLogin(req user.LoginRequest) (interface{}, error) {
	ClientMode := &common.ClientMode{}
	ClientUrl := ClientMode.GetClientUrl(common.GetClientUrl{"micro-authority"})
	url := ClientUrl + "/micro-authority/v1/account/scrm-password-ip-white-login"
	res, err := client.Post(ctx.New(), url, clients.JSONBody(req))
	helper.Must(err)
	return res.Read(func(bytes []byte) (interface{}, error) {
		rt := fmt.Sprint(string(bytes))
		return rt, nil
	})
}

// 获取图片验证码
func (a ActiveApiImpl) GetValidateCode() (interface{}, error) {
	ClientMode := &common.ClientMode{}
	ClientUrl := ClientMode.GetClientUrl(common.GetClientUrl{"micro-authority"})
	url := ClientUrl + "/micro-authority/v1/common/validate-code"
	res, err := client.Post(ctx.New(), url)
	helper.Must(err)
	return res.Read(func(bytes []byte) (interface{}, error) {
		rt := fmt.Sprint(string(bytes))
		return rt, nil
	})
}

// 获取token
func (a ActiveApiImpl) GetToken(req user.GetTokenRequest) (interface{}, error) {
	ClientMode := &common.ClientMode{}
	ClientUrl := ClientMode.GetClientUrl(common.GetClientUrl{"micro-authority"})
	url := ClientUrl + "/micro-authority/v1/account/get-token"
	res, err := client.Post(ctx.New(), url, clients.JSONBody(req))
	helper.Must(err)
	return res.Read(func(bytes []byte) (interface{}, error) {
		rt := fmt.Sprint(string(bytes))
		return rt, nil
	})
}

// 获取旧版token
func (a ActiveApiImpl) GetTokenV1(req user.GetTokenRequestV1) (interface{}, error) {
	ClientMode := &common.ClientMode{}
	ClientUrl := ClientMode.GetClientUrl(common.GetClientUrl{"console"})
	url := ClientUrl + "/gw-console/v1/admin/get-tokenv1"
	res, err := client.Post(ctx.New(), url, clients.JSONBody(req))
	helper.Must(err)
	return res.Read(func(bytes []byte) (interface{}, error) {
		var getTokenResponse user.TokenResponse
		errs := json.Unmarshal(bytes, &getTokenResponse)
		helper.Must(errs)
		return getTokenResponse.Data.Data, nil
	})
}

// v1版本用户数据
func (a ActiveApiImpl) GetAdminInfo(req user.GetTokenRequest) (interface{}, error) {
	ClientMode := &common.ClientMode{}
	ClientUrl := ClientMode.GetClientUrl(common.GetClientUrl{"console"})
	url := ClientUrl + "/gw-console/v1/admin/adminInfo"
	res, err := client.Post(ctx.New(), url, clients.JSONBody(req))
	helper.Must(err)
	return res.Read(func(bytes []byte) (interface{}, error) {
		rt := fmt.Sprint(string(bytes))
		var getTokenResponse user.TokenResponse
		errs := json.Unmarshal([]byte(rt), &getTokenResponse)
		helper.Must(errs)
		return getTokenResponse.Data, nil
	})
}

// 获取js api
func (a ActiveApiImpl) GetJsApi(req user.JsApiRequest) (interface{}, error) {
	ClientMode := &common.ClientMode{}
	ClientUrl := ClientMode.GetClientUrl(common.GetClientUrl{"work"})
	url := ClientUrl + "/service/v1/auth/agent_config"
	res, err := client.Post(ctx.New(), url, clients.JSONBody(req))
	helper.Must(err)
	return res.Read(func(bytes []byte) (interface{}, error) {
		rt := fmt.Sprint(string(bytes))
		return rt, nil
	})
}

// 企业微信二维码
func (a ActiveApiImpl) QyCode(req user.QyCodeRequest) string {

	url := "https://www.vchangyi.com/sso/index.php?corp_id=wxa7044ee8255576b0&ref=" + req.Ref + "&user_type=member"
	return url
}

func (a ActiveApiImpl) GetIndustry(params user.GetIndustry) (interface{}, error) {
	ClientMode := &common.ClientMode{}
	ClientUrl := ClientMode.GetClientUrl(common.GetClientUrl{"boss"})
	url := ClientUrl + "/service/v1/industry/list"
	res, err := client.Post(ctx.New(), url, clients.JSONBody(params))
	helper.Must(err)
	return res.Read(func(bytes []byte) (interface{}, error) {
		ret := fmt.Sprint(string(bytes))
		return ret, nil
	})
}

func (a ActiveApiImpl) CheckEmail(params user.CheckEmail) (interface{}, error) {

	ClientMode := &common.ClientMode{}
	ClientUrl := ClientMode.GetClientUrl(common.GetClientUrl{"boss"})

	url := ClientUrl + "/service/v1/enterprise-admin/check-email"
	res, err := client.Post(ctx.New(), url, clients.JSONBody(params))
	helper.Must(err)
	return res.Read(func(bytes []byte) (interface{}, error) {
		rt := fmt.Sprint(string(bytes))
		return rt, nil
	})
}

func (a ActiveApiImpl) CheckMobileCode(params user.CheckMobileCode) (interface{}, error) {
	ClientMode := &common.ClientMode{}
	ClientUrl := ClientMode.GetClientUrl(common.GetClientUrl{"micro-authority"})
	url := ClientUrl + "/micro-authority/v1/common/check-code"
	res, err := client.Post(ctx.New(), url, clients.JSONBody(params))
	helper.Must(err)
	return res.Read(func(bytes []byte) (interface{}, error) {
		ret := fmt.Sprint(string(bytes))
		return ret, nil
	})
}

func (a ActiveApiImpl) AutoAccess(params user.AutoAccess) (interface{}, error) {

	ClientMode := &common.ClientMode{}
	ClientUrl := ClientMode.GetClientUrl(common.GetClientUrl{"boss"})

	url := ClientUrl + "/service/v1/enterprise/add"
	res, err := client.Post(ctx.New(), url, clients.JSONBody(params))
	helper.Must(err)
	return res.Read(func(bytes []byte) (interface{}, error) {
		rt := fmt.Sprint(string(bytes))
		return rt, nil
	})
}

/**
获取可安装附件列表
*/
func (a ActiveApiImpl) AdditionalList(params user.AdditionalList) (interface{}, error) {

	ClientMode := &common.ClientMode{}
	ClientUrl := ClientMode.GetClientUrl(common.GetClientUrl{"micro-authority"})

	url := ClientUrl + "/micro-authority/v1/rbac/additional-list"
	res, err := client.Post(ctx.New(), url, clients.JSONBody(params))
	helper.Must(err)
	return res.Read(func(bytes []byte) (interface{}, error) {
		rt := fmt.Sprint(string(bytes))
		return rt, nil
	})
}

func (a ActiveApiImpl) TmpToken(params user.TmpToken) (interface{}, error) {

	ClientMode := &common.ClientMode{}
	ClientUrl := ClientMode.GetClientUrl(common.GetClientUrl{"console"})

	url := ClientUrl + "/service/v1/admin/auto-account"
	res, err := client.Post(ctx.New(), url, clients.JSONBody(params))
	helper.Must(err)
	return res.Read(func(bytes []byte) (interface{}, error) {
		rt := fmt.Sprint(string(bytes))
		return rt, nil
	})
}

func (a ActiveApiImpl) InstallThirdApp(params user.InstallThirdApp) string {

	ClientMode := &common.ClientMode{}
	UcUrl := ClientMode.GetUcUrl()

	suiteId := ClientMode.GetAppSuiteId()

	installUrl := UcUrl + "/s/suite/install?suiteId=" + suiteId + "&callbackUrl=" + params.BackUrl

	return installUrl
}

func (a ActiveApiImpl) GetRegisterCode(params user.GetRegisterCode) (interface{}, error) {

	ClientMode := &common.ClientMode{}
	UcUrl := ClientMode.GetUcUrl()
	url := UcUrl + "/s/enterprise/get-register-code"

	res, err := client.Post(ctx.New(), url, clients.JSONBody(params))
	helper.Must(err)
	return res.Read(func(bytes []byte) (interface{}, error) {
		rt := fmt.Sprint(string(bytes))
		return rt, nil
	})
}

// WorkWechatLogin 企微登录
func (a ActiveApiImpl) WorkWechatLogin(params user.WorkWechatLoginRequestData) (interface{}, error) {
	ClientMode := &common.ClientMode{}

	ClientUrl := ClientMode.GetClientUrl(common.GetClientUrl{Domain: "console"})
	url := ClientUrl + "/gw-console/v1/admin/work-wechat-login"

	res, err := client.Post(ctx.New(), url, clients.JSONBody(params))
	helper.Must(err)
	return res.Read(func(bytes []byte) (interface{}, error) {
		rt := fmt.Sprint(string(bytes))
		return rt, nil
	})
}

// SendTemplateMail 发送邮件
func (a ActiveApiImpl) SendTemplateMail(ctx ctx.Context, req map[string]string) (*user.SendCloudResponse, error) {
	ClientMode := &common.ClientMode{}

	ClientUrl := ClientMode.GetClientUrl(common.GetClientUrl{"sendcloud"})
	url := ClientUrl + "/apiv2/mail/sendtemplate"

	res, err := client.Post(ctx, url, clients.Form(req))
	helper.Must(err)
	rResp, rErr := res.Read(func(bytes []byte) (interface{}, error) {
		rt := fmt.Sprint(string(bytes))
		return rt, nil
	})
	if rErr != nil {
		return nil, rErr
	}
	sResp := &user.SendCloudResponse{}
	jErr := json.Unmarshal([]byte(rResp.(string)), sResp)
	if jErr != nil {
		return nil, jErr
	}
	return sResp, nil
}
