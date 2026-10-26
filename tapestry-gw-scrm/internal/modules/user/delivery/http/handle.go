package user_http

import (
	"encoding/json"
	"errors"
	"fmt"
	"gw-scrm/internal/modules/version"
	helper2 "gw-scrm/internal/pkg/helper"
	"gw-scrm/internal/pkg/sdk/microscrm"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/helper"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"gw-scrm/internal/config"
	"gw-scrm/internal/modules/user"
	"gw-scrm/internal/modules/user/api"
	"gw-scrm/internal/pkg/common"
	scrm_errs "gw-scrm/internal/pkg/errs"
	"gw-scrm/internal/pkg/redis"
)

// 获取客户端标识
func getClient() string {
	ClientMode := &common.ClientMode{}
	client := ClientMode.GetThisClient()
	return client
}

const client = "console"

// 获取验证码
func getCode(c *gin.Context) {
	var requestInfo user.GetCodeRequest
	err := c.BindJSON(&requestInfo)
	helper.Must(err)
	requestInfo.Account = "+86-" + requestInfo.Mobile
	requestInfo.Client = getClient()
	act := &api.ActiveApiImpl{}
	if requestInfo.Type == 1 {
		cres, cerr := act.CheckMobile(requestInfo)
		helper.Must(cerr)
		cstr := cres.(string)
		var checkMobileResponse user.GwResponse
		err := json.Unmarshal([]byte(cstr), &checkMobileResponse)
		helper.Must(err)
		if checkMobileResponse.Code != 200 {
			if checkMobileResponse.Code == 700005 {
				checkMobileResponse.Message = scrm_errs.MobileErr.Message
			}
			ginplus.ResError(c, errorx.NewErrorMsg(int(checkMobileResponse.Code), checkMobileResponse.Message))
			return
		}
	}
	gres, gerr := act.GetCode(requestInfo)
	helper.Must(gerr)
	gstr := gres.(string)
	var getCodeResponse user.GwResponse
	gerrs := json.Unmarshal([]byte(gstr), &getCodeResponse)
	helper.Must(gerrs)
	if getCodeResponse.Code != 200 {
		ginplus.ResError(c, errorx.NewErrorMsg(int(getCodeResponse.Code), getCodeResponse.Message))
		return
	}
	ginplus.ResSuccess(c, getCodeResponse.Data)
}

// 获取短信验证码
func getBindCode(c *gin.Context) {
	var requestInfo user.GetBindCodeRequest
	err := c.BindJSON(&requestInfo)
	helper.Must(err)
	act := &api.ActiveApiImpl{}
	requestInfo.Account = "+86-" + requestInfo.Mobile
	requestInfo.Client = getClient()
	requestInfo.Type = 1
	cres, cerr := act.CheckBindMobile(requestInfo)
	helper.Must(cerr)
	cstr := cres.(string)
	var checkMobileResponse user.GwResponse
	cerrs := json.Unmarshal([]byte(cstr), &checkMobileResponse)
	helper.Must(cerrs)
	if checkMobileResponse.Code != 0 {
		checkMobileResponse.Message = scrm_errs.BindMobileErr.Message
		ginplus.ResSuccess(c, checkMobileResponse.Data)
		return
	}
	gres, gerr := act.GetBindCode(requestInfo)
	helper.Must(gerr)
	gstr := gres.(string)
	var getCodeResponse user.GwResponse
	gerrs := json.Unmarshal([]byte(gstr), &getCodeResponse)
	helper.Must(gerrs)
	ginplus.ResSuccess(c, getCodeResponse.Data)
}

// 获取验证码
func codeForOauth(requestInfo user.GetCodeRequest) error {
	requestInfo.Account = "+86-" + requestInfo.Mobile
	requestInfo.Client = getClient()
	act := &api.ActiveApiImpl{}
	if requestInfo.Type == 1 {
		cres, cerr := act.CheckMobile(requestInfo)
		if cerr != nil {
			return cerr
		}
		cstr := cres.(string)
		var checkMobileResponse user.GwResponse
		cerr = json.Unmarshal([]byte(cstr), &checkMobileResponse)
		if cerr != nil {
			return cerr
		}
		if checkMobileResponse.Code != 200 {
			checkMobileResponse.Message = scrm_errs.AccountErr.Message
			return errorx.NewErrorMsg(int(checkMobileResponse.Code), checkMobileResponse.Message)
		}
	}
	gres, cerr := act.GetCode(requestInfo)
	if cerr != nil {
		return cerr
	}
	gstr := gres.(string)
	var getCodeResponse user.GwResponse
	gerrs := json.Unmarshal([]byte(gstr), &getCodeResponse)
	if gerrs != nil {
		return gerrs
	}
	if getCodeResponse.Code != 200 {
		getCodeResponse.Message = "操作频繁,请五分钟后再试"
		return errorx.NewErrorMsg(int(getCodeResponse.Code), getCodeResponse.Message)
	}
	return nil
}

func openValidate(c *gin.Context) {
	ginplus.ResSuccess(c, nil)
}
func oauth(c *gin.Context) {
	var requestInfo user.LoginRequest
	jErr := c.BindJSON(&requestInfo)
	if jErr != nil {
		ginplus.ResError(c, jErr)
		return
	}
	requestInfo.LoginType = 1
	requestInfo.Code = "9118"
	requestInfo.Client = client
	act := &api.ActiveApiImpl{}
	res, rErr := act.Login(requestInfo)

	if rErr != nil {
		ginplus.ResError(c, rErr)
		return
	}
	str := res.(string)
	var loginResponse user.AdminResponse
	errs := json.Unmarshal([]byte(str), &loginResponse)
	if errs != nil {
		ginplus.ResError(c, errs)
		return
	}
	codeReq := user.GetCodeRequest{
		Mobile:       loginResponse.Data.Token,
		Client:       "",
		Type:         requestInfo.LoginType,
		EnterpriseId: requestInfo.EnterpriseId,
	}

	errs = codeForOauth(codeReq)
	if errs != nil {
		ginplus.ResError(c, errs)
		return
	}
	// 通过token和短信验证码登录
	logReq := user.LoginRequest{
		Account:      loginResponse.Data.Token,
		Password:     "REDACTED_FOR_SOURCE_SHARING",
		LoginType:    2,
		EnterpriseId: loginResponse.Data.EnterpriseId,
		Client:       client,
	}
	res, rErr = act.Login(logReq)
	if rErr != nil {
		ginplus.ResError(c, errs)
		return
	}
	var tokenResponse user.AdminResponse
	errs = json.Unmarshal([]byte(res.(string)), &tokenResponse)
	if errs != nil {
		ginplus.ResError(c, errs)
		return
	}
	if tokenResponse.Code == 700005 || tokenResponse.Code == 700006 {
		if requestInfo.LoginType == 1 {
			tokenResponse.Message = scrm_errs.AccountErr.Message
		} else {
			tokenResponse.Message = scrm_errs.MobileErr.Message
		}
		ginplus.ResError(c, errorx.NewErrorMsg(int(tokenResponse.Code), tokenResponse.Message))
		return
	}
	var adminId = int64(0)
	for _, data := range tokenResponse.Data.EnterpriseAdminList {
		if data.EnterpriseId == requestInfo.EnterpriseId {
			adminId = data.AdminId
		}
	}
	adminReq := user.GetTokenRequest{
		AdminId:      adminId,
		EnterpriseId: requestInfo.EnterpriseId,
		Client:       client,
	}
	res, rerr := act.GetToken(adminReq)
	if rerr != nil {
		ginplus.ResError(c, errs)
		return
	}
	var getTokenResponse user.TokenResponse
	errs = json.Unmarshal([]byte(res.(string)), &getTokenResponse)
	if errs != nil {
		ginplus.ResError(c, errs)
		return
	}
	// 成功后存储token
	if getTokenResponse.Code == http.StatusOK {
		tokenInfo := getTokenResponse.Data.TokenBytes
		key := fmt.Sprintf("%s:token:%s", config.CfgData.App, getTokenResponse.Data.Token)
		jsonStr, err := common.JsonEncode(tokenInfo)
		if err != nil {
			ginplus.ResError(c, err)
			return
		}
		redis.Client().Set(c, key, jsonStr, 15*60*10*time.Second)
	}
	td := user.OauthData{
		Token: getTokenResponse.Data.Token,
	}
	ginplus.ResJSON(c, http.StatusOK, td)
}

// login 登录
func login(c *gin.Context) {
	var requestInfo user.LoginRequest
	jErr := c.BindJSON(&requestInfo)
	if jErr != nil {
		ginplus.ResError(c, jErr)
		return
	}
	if requestInfo.EnterpriseId == 0 {
		ginplus.ResError(c, scrm_errs.EnterpriseIdNull)
		return
	}
	if requestInfo.LoginType == 1 {
		//改为rsa加密密码
		requestInfo.Password = helper2.RsaDecode(requestInfo.Password)
	}
	if requestInfo.Password == "" {
		ginplus.ResError(c, errors.New("密码错误1"))
		return
	}
	requestInfo.Client = getClient()
	act := &api.ActiveApiImpl{}
	res, rErr := act.Login(requestInfo)

	if rErr != nil {
		ginplus.ResError(c, rErr)
		return
	}
	str := res.(string)
	var loginResponse user.AdminResponse
	errs := json.Unmarshal([]byte(str), &loginResponse)
	if errs != nil {
		ginplus.ResError(c, errs)
		return
	}
	if loginResponse.Code == 700005 || loginResponse.Code == 700006 {
		if requestInfo.LoginType == 1 {
			loginResponse.Message = scrm_errs.AccountErr.Message
		} else {
			loginResponse.Message = scrm_errs.MobileErr.Message
		}
		ginplus.ResError(c, errorx.NewErrorMsg(int(loginResponse.Code), loginResponse.Message))
		return
	}

	if len(loginResponse.Data.EnterpriseAdminList) != 0 {
		loginResponse.Data.Version = loginResponse.Data.EnterpriseAdminList[0].Version
		loginResponse.Data.EnterpriseTotal = len(loginResponse.Data.EnterpriseAdminList)
		if loginResponse.Data.EnterpriseAdminList[0].Version == 1 {
			// 获取V1版本管理员数据
			var requestV1Info user.GetTokenRequestV1
			requestV1Info.AdminId = loginResponse.Data.EnterpriseAdminList[0].AdminId
			requestV1Info.ServeType = loginResponse.Data.EnterpriseAdminList[0].ServeType
			requestV1Info.EnterpriseId = loginResponse.Data.EnterpriseAdminList[0].EnterpriseId
			requestV1Info.Role = loginResponse.Data.EnterpriseAdminList[0].Role
			requestV1Info.Account = loginResponse.Data.EnterpriseAdminList[0].Account
			requestV1Info.AccountMobile = loginResponse.Data.EnterpriseAdminList[0].AccountMobile
			ares, aerr := act.GetTokenV1(requestV1Info)
			if aerr != nil {
				ginplus.ResError(c, aerr)
				return
			}
			loginResponse.Data.Data = ares
		}
		if loginResponse.Data.EnterpriseAdminList[0].IsSplitFlow == 1 {
			loginResponse.Data.Hash = loginResponse.Data.EnterpriseAdminList[0].Hash
		} else {
			loginResponse.Data.Hash = "store"
		}
	}
	loginResponse.Data.EnterpriseList = loginResponse.Data.EnterpriseAdminList
	loginResponse.Data.EnterpriseAdminList = nil

	ginplus.ResJSON(c, http.StatusOK, loginResponse)
}

// ip白名单登录，不需要验证码
func whiteLogin(c *gin.Context) {
	var requestInfo user.LoginRequest
	jErr := c.BindJSON(&requestInfo)
	if jErr != nil {
		ginplus.ResError(c, jErr)
		return
	}
	if requestInfo.EnterpriseId == 0 {
		ginplus.ResError(c, scrm_errs.EnterpriseIdNull)
		return
	}
	//改为rsa加密密码
	requestInfo.Password = helper2.RsaDecode(requestInfo.Password)
	requestInfo.Client = getClient()
	act := &api.ActiveApiImpl{}
	res, rErr := act.ConsoleIpWhiteLogin(requestInfo)

	if rErr != nil {
		ginplus.ResError(c, rErr)
		return
	}
	str := res.(string)
	var loginResponse user.AdminResponse
	errs := json.Unmarshal([]byte(str), &loginResponse)
	if errs != nil {
		ginplus.ResError(c, errs)
		return
	}
	if loginResponse.Code == 700005 || loginResponse.Code == 700006 {
		if requestInfo.LoginType == 1 {
			loginResponse.Message = scrm_errs.AccountErr.Message
		} else {
			loginResponse.Message = scrm_errs.MobileErr.Message
		}
		ginplus.ResError(c, errorx.NewErrorMsg(int(loginResponse.Code), loginResponse.Message))
		return
	}

	if len(loginResponse.Data.EnterpriseAdminList) != 0 {
		loginResponse.Data.Version = loginResponse.Data.EnterpriseAdminList[0].Version
		loginResponse.Data.EnterpriseTotal = len(loginResponse.Data.EnterpriseAdminList)
		if loginResponse.Data.EnterpriseAdminList[0].Version == 1 {
			// 获取V1版本管理员数据
			var requestV1Info user.GetTokenRequestV1
			requestV1Info.AdminId = loginResponse.Data.EnterpriseAdminList[0].AdminId
			requestV1Info.ServeType = loginResponse.Data.EnterpriseAdminList[0].ServeType
			requestV1Info.EnterpriseId = loginResponse.Data.EnterpriseAdminList[0].EnterpriseId
			requestV1Info.Role = loginResponse.Data.EnterpriseAdminList[0].Role
			requestV1Info.Account = loginResponse.Data.EnterpriseAdminList[0].Account
			requestV1Info.AccountMobile = loginResponse.Data.EnterpriseAdminList[0].AccountMobile
			ares, aerr := act.GetTokenV1(requestV1Info)
			if aerr != nil {
				ginplus.ResError(c, aerr)
				return
			}
			loginResponse.Data.Data = ares
		}
		if loginResponse.Data.EnterpriseAdminList[0].IsSplitFlow == 1 {
			loginResponse.Data.Hash = loginResponse.Data.EnterpriseAdminList[0].Hash
		} else {
			loginResponse.Data.Hash = "store"
		}
	}
	loginResponse.Data.EnterpriseList = loginResponse.Data.EnterpriseAdminList
	loginResponse.Data.EnterpriseAdminList = nil
	log.Info(ctx.Wrap(c)).Msgf("loginAdmin: %+v", loginResponse)
	ginplus.ResJSON(c, http.StatusOK, loginResponse)
}

// 用户退出
func logout(c *gin.Context) {
	token := c.GetHeader(common.HeaderKeyEnum.AuthToken)
	key := fmt.Sprintf("%s:token:%s", config.CfgData.App, token)
	redis.Client().Del(c, key)
	ginplus.ResOK(c)
	return
}

// 获取图片验证码
func getVCode(c *gin.Context) {
	act := &api.ActiveApiImpl{}
	res, rErr := act.GetValidateCode()
	if rErr != nil {
		ginplus.ResError(c, rErr)
		return
	}
	c.String(http.StatusOK, res.(string))
}

func changeEnterpriseList(c *gin.Context) {
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	adminId := tokenInfo.AdminId
	req := map[string]interface{}{
		"admin_id": adminId,
	}
	res, err := api.AccountBaseInfo(ctx.Wrap(c), req)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	accountBaseRes := res.(map[string]interface{})
	if nil == accountBaseRes["account"] {
		ginplus.ResError(c, scrm_errs.AccountNotExists)
		return
	}

	newReq := map[string]interface{}{
		"account": accountBaseRes["account"],
		"version": "v2",
		"client":  "console",
	}
	newRes, err := api.EnterpriseList(ctx.Wrap(c), newReq)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	ginplus.ResSuccess(c, newRes)
}

func getCurrentVersionInfo(c *gin.Context) {
	var param version.GetCurrentVersionInfoRequest
	if err := c.ShouldBindJSON(&param); err != nil {
		ginplus.ResError(c, err)
		return
	}

	versionInfo, err := microscrm.GetCurrentVersionInfo(ctx.Wrap(c), param)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	ginplus.ResSuccess(c, versionInfo)
}

// 获取token
func getToken(c *gin.Context) {
	var requestInfo user.GetTokenRequest
	err := c.BindJSON(&requestInfo)
	helper.Must(err)
	if requestInfo.EnterpriseId == 0 {
		ginplus.ResError(c, scrm_errs.EnterpriseIdNull)
		return
	}
	appCode, ok := user.NormalizeAppCode(requestInfo.AppCode)
	if !ok {
		ginplus.ResError(c, errors.New("不支持的app_code"))
		return
	}
	requestInfo.AppCode = appCode
	requestInfo.Client = user.ClientForAppCode(appCode)
	act := &api.ActiveApiImpl{}
	res, rerr := act.GetToken(requestInfo)
	helper.Must(rerr)
	str := res.(string)
	// c.String(http.StatusOK, str)

	var getTokenResponse user.TokenResponse
	errs := json.Unmarshal([]byte(str), &getTokenResponse)
	helper.Must(errs)
	// 沿用console方法
	if getTokenResponse.Data.EnterpriseAdminInfo.EnterpriseId != requestInfo.EnterpriseId {
		ginplus.ResError(c, scrm_errs.AccountNotExists)
		return
	}
	if getTokenResponse.Data.EnterpriseAdminInfo.IsSplitFlow == 1 {
		getTokenResponse.Data.Hash = getTokenResponse.Data.EnterpriseAdminInfo.Hash
	} else {
		if getTokenResponse.Data.EnterpriseAdminInfo.Version == 1 {
			getTokenResponse.Data.Hash = "REDACTED_FOR_SOURCE_SHARING"
		} else {
			getTokenResponse.Data.Hash = "REDACTED_FOR_SOURCE_SHARING"
		}
	}
	if getTokenResponse.Data.EnterpriseAdminInfo.Version == 1 {
		// 获取V1版本管理员数据
		var getTokenRequestV1 user.GetTokenRequestV1
		getTokenRequestV1.EnterpriseId = getTokenResponse.Data.EnterpriseAdminInfo.EnterpriseId
		getTokenRequestV1.ServeType = getTokenResponse.Data.EnterpriseAdminInfo.ServeType
		getTokenRequestV1.Account = getTokenResponse.Data.EnterpriseAdminInfo.Account
		getTokenRequestV1.AccountMobile = getTokenResponse.Data.EnterpriseAdminInfo.AccountMobile
		getTokenRequestV1.Role = getTokenResponse.Data.EnterpriseAdminInfo.Role
		getTokenRequestV1.AdminId = getTokenResponse.Data.EnterpriseAdminInfo.AdminId
		ares, aerr := act.GetTokenV1(getTokenRequestV1)
		helper.Must(aerr)
		getTokenResponse.Data.Data = ares
	}
	getTokenResponse.Data.Version = getTokenResponse.Data.EnterpriseAdminInfo.Version
	getTokenResponse.Data.Total = 1
	// 兼容其他版本，冗余数据
	adminList := make(map[int8]interface{})
	adminList[0] = getTokenResponse.Data.EnterpriseAdminInfo
	getTokenResponse.Data.AdminList = adminList
	jses, jserr := act.GetJsApi(user.JsApiRequest{getTokenResponse.Data.EnterpriseAdminInfo.Domain, "crs_store_assistant"})
	helper.Must(jserr)
	getTokenResponse.Data.JsApiConfig = jses

	// 成功后存储token
	if getTokenResponse.Code == http.StatusOK {
		tokenInfo := getTokenResponse.Data.TokenBytes
		// 计算token剩余时间
		key := fmt.Sprintf("%s:token:%s", config.CfgData.App, getTokenResponse.Data.Token)
		jsonStr, err := common.JsonEncode(tokenInfo)
		helper.Must(err)
		redis.Client().Set(c, key, jsonStr, 15*60*10*time.Second)

		// 删除短信验证码
		// redis.Client().Do(c,"select",0)
		// mobile := getTokenResponse.Data.EnterpriseAdminInfo.AccountMobile
		// smsKey := fmt.Sprintf("micro:authority:easy-sms:console:mobile:%s",mobile)
		// redis.Client().Del(c,smsKey)
		// redis.Client().Do(c,"select",config.CfgData.Redis.DB)
	}

	ginplus.ResSuccess(c, getTokenResponse.Data)
}

// 企业微信二维码
func QyCode(c *gin.Context) {
	ref := c.Query("ref")
	act := &api.ActiveApiImpl{}
	req := user.QyCodeRequest{
		ref,
	}
	res := act.QyCode(req)
	c.Header("location", res)
	c.JSON(http.StatusFound, res)
}

/*
*
获取行业列表
*/
func getIndustryList(c *gin.Context) {
	var params user.GetIndustry
	err := c.BindJSON(&params)
	helper.Must(err)

	act := &api.ActiveApiImpl{}
	res, err := act.GetIndustry(params)
	helper.Must(err)

	str := res.(string)
	c.Header("Content-Type", "application/json")
	c.String(http.StatusOK, str)
}

/*
*
验证邮箱是否已存在
*/
func CheckEmail(c *gin.Context) {
	var params user.CheckEmail
	err := c.BindJSON(&params)
	helper.Must(err)

	act := &api.ActiveApiImpl{}
	res, err := act.CheckEmail(params)

	helper.Must(err)
	str := res.(string)
	c.Header("Content-Type", "application/json")
	c.String(http.StatusOK, str)
}

type ResponseData struct {
	Code int `json:"code"`
}

type ResponseautoAccessres struct {
	Data AutoAccessresData `json:"data"`
}

type AutoAccessresData struct {
	Admin Admin `json:"admin"`
}

type Admin struct {
	Id        int64  `json:"id"`
	OriginPwd string `json:"originPwd"`
}

type AutoAccessresRes struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Data    AutoAccessresAdminData `json:"data"`
}

type AutoAccessresAdminData struct {
	Admin AutoAccessresAdmin `json:"admin"`
}

type AutoAccessresAdmin struct {
	AdminId int64 `json:"admin_id"`
}

type Register struct {
	Data RegisterData `json:"data"`
}
type RegisterData struct {
	RegisterCode string `json:"registerCode"`
}

// AutoAccess 自助开户
func AutoAccess(c *gin.Context) {
	act := &api.ActiveApiImpl{}

	var autoAccessParams user.AutoAccessParams
	autoAccessParamsErr := c.BindJSON(&autoAccessParams)
	helper.Must(autoAccessParamsErr)

	// 手机号验证验证
	var checkMobileParams user.CheckMobileCode
	checkMobileParams.Client = getClient()
	checkMobileParams.Mobile = autoAccessParams.Mobile
	checkMobileParams.CerificationCode = autoAccessParams.VerificationCode
	checkMobileParams.Code = autoAccessParams.VerificationCode

	res, mobileCodeErr := act.CheckMobileCode(checkMobileParams)
	if mobileCodeErr != nil {
		ginplus.ResError(c, mobileCodeErr)
		return
	}
	CheckMobileStr := res.(string)

	var responseData ResponseData
	json.Unmarshal([]byte(CheckMobileStr), &responseData)

	if responseData.Code != 200 {
		c.Header("Content-Type", "application/json")
		c.String(http.StatusOK, CheckMobileStr)
		return
	}

	// 自助注册数据提交
	startTime := time.Now().Unix() - 10
	endTime := startTime + 3600*24*60

	var params user.AutoAccess
	params.Version = autoAccessParams.Version
	params.CorpId = autoAccessParams.CorpId
	params.Name = autoAccessParams.Name
	params.Contact = autoAccessParams.Contact
	params.Mobile = autoAccessParams.Mobile
	params.Email = autoAccessParams.Email
	params.IndustryId = autoAccessParams.IndustryId
	params.AccountLimit = user.AccountDefaultLimit
	params.BusinessSource = user.BusinessSourceDefault
	params.MemUserId = autoAccessParams.MemUserId
	params.AvailableSystem = make(map[int8]int8)
	params.AvailableSystem[0] = user.AvailableSystemDefaultWorker
	params.AvailableSystem[1] = user.AvailableSystemDefaultShopping

	// 获取可使用附加功能
	var additionalListParams user.AdditionalList
	additionalListParams.Version = params.Version
	additionalList, additionalListErr := act.AdditionalList(additionalListParams)
	helper.Must(additionalListErr)
	additionalListStr := additionalList.(string)

	var gwResponse user.GwTmpResponse
	json.Unmarshal([]byte(additionalListStr), &gwResponse)
	var idList []int
	for _, id := range gwResponse.Data {
		idList = append(idList, id.Id)
	}

	// 组装开通企业必要参数
	params.VersionTag = "试用版"
	params.AccountStartAt = startTime
	params.AccountExpireAt = endTime
	params.AermissionList = idList

	autoAccessres, autoAccessresErr := act.AutoAccess(params)
	if autoAccessresErr != nil {
		ginplus.ResError(c, autoAccessresErr)
		return
	}
	autoAccessresStr := autoAccessres.(string)

	var responseautoAccessres ResponseautoAccessres
	json.Unmarshal([]byte(autoAccessresStr), &responseautoAccessres)
	if responseautoAccessres.Data.Admin.Id == 0 {
		c.Header("Content-Type", "application/json")
		c.String(http.StatusOK, autoAccessresStr)
		return
	}

	var autoAccessresRes AutoAccessresRes
	autoAccessresRes.Data.Admin.AdminId = responseautoAccessres.Data.Admin.Id
	autoAccessresRes.Code = 200
	autoAccessresRes.Message = "注册成功"
	ginplus.ResSuccess(c, autoAccessresRes.Data)
}

/*
*

	获取安装智慧导购二维码
*/
func installThirdApp(c *gin.Context) {
	var params user.InstallThirdApp
	err := c.BindJSON(&params)
	helper.Must(err)

	act := api.ActiveApiImpl{}

	installUrl := act.InstallThirdApp(params)

	var GetRegisterCode user.GetRegisterCode
	GetRegisterCode.TemplateId = user.TemplateDefaultId
	GetRegisterCode.ProviderCorpId = user.ProviderCorpId
	res, getRegisterCodeErr := act.GetRegisterCode(GetRegisterCode)
	helper.Must(getRegisterCodeErr)
	GetRegisterStr := res.(string)

	var Register Register
	json.Unmarshal([]byte(GetRegisterStr), &Register)

	c.JSON(http.StatusOK, map[string]string{"url": installUrl, "qy_register_url": "https://open.work.weixin.qq.com/3rdservice/wework/register?register_code=" + Register.Data.RegisterCode})
}

// 企业微信扫码登录
func WorkWechatLogin(c *gin.Context) {
	var params user.WorkWechatLoginRequestData
	params.BindAccount = make(map[string]string)

	paramsErr := c.BindJSON(&params)
	helper.Must(paramsErr)

	act := &api.ActiveApiImpl{}
	params.RegisterEnterprise = make(map[string]string)
	res, WorkWechatLoginErr := act.WorkWechatLogin(params)
	helper.Must(WorkWechatLoginErr)

	str := res.(string)
	c.Header("Content-Type", "application/json")
	c.String(http.StatusOK, str)

	// 初始化企业微信扫码登录返回结果
	/*workWechatLoginResponse := &entity.WorkWechatLoginResponse{}
	workWechatLoginResponse.Code = "success"
	workWechatLoginResponse.Desc = "可以登录"

	c.BindJSON(workWechatLoginRequest)
	if workWechatLoginRequest.Type == "scan_code" {
		// 根据auth_code查询企业信息
		authInfo, error := authInfoByCode(workWechatLoginRequest, workWechatLoginResponse)
		if error == nil {
			ginplus.ResJSON(c, http.StatusOK, workWechatLoginResponse)
			return
		}

		// 根据epId 查询Uc企业信息
		res, err := api.EnSettingQuery(authInfo.EnterpriseInfo.EpId)
		helper.Must(err)
		str := res.(string)
		ucEnSettingQuery := entity.UcEnSettingQuery{}
		json.Unmarshal([]byte(str), ucEnSettingQuery)
		// 查询企业信息
		_, err = enterpriseDetail(ucEnSettingQuery.WxMpAppId, workWechatLoginRequest, workWechatLoginResponse)
		helper.Must(err)
		ginplus.ResJSON(c, http.StatusOK, workWechatLoginResponse)
		// uc超级管理员信息
		res, err = api.AdminDetail(authInfo.EnterpriseInfo.EpId, authInfo.AdminInfo.EaId)
		helper.Must(err)
		str = res.(string)
		ucSuperAdmin := &user.UcAdminDetail{}
		json.Unmarshal([]byte(str), ucSuperAdmin)

		// 查询Boss系统管理员信息
		req := &entity.BossEnterpriseAdminDetailRequest{MemberId: ucSuperAdmin.EaId}
		res, err = api.BossAdminDetail(req)
		helper.Must(err)
		str = res.(string)
		bossEnterpriseAdmin := &entity.BossEnterpriseAdminDetail{}
		json.Unmarshal([]byte(str), bossEnterpriseAdmin)

		// 没有查询到管理员信息
		if bossEnterpriseAdmin.AdminId == 0 {
			if workWechatLoginRequest.BindAccount.Mobile != "" {
				req := &entity.AuthorityCheckSmsCodeRequest{
					Mobile: workWechatLoginRequest.BindAccount.Mobile,
					Code:   workWechatLoginRequest.BindAccount.Code,
					Event:  "word_wechat_bind_account",
				}
				if res, err = api.AuthorityCheckSmsCode(req); err == nil {
					workWechatLoginResponse.Code = "verify_code_error"
					workWechatLoginResponse.Desc = "验证码错误"

					ginplus.ResJSON(c, http.StatusOK, workWechatLoginResponse)
					return
				}

				// 绑定管理员信息
				bindAdmin()
			} else {
				workWechatLoginResponse.Code = "admin_not_found"
				workWechatLoginResponse.Desc = "管理员不存在"

				ginplus.ResJSON(c, http.StatusOK, workWechatLoginResponse)
				return
			}
		}

		// 生成Token
	}*/
}

// func authInfoByCode(workWechatLoginRequest *user.WorkWechatLoginRequest, workWechatLoginResponse *entity.WorkWechatLoginResponse) (*user.UcAuthInfo, error) {
//	// 先根据authCode查询企业信息是否存在
//	res, err := api.AuthCodeInfo(workWechatLoginRequest.AuthCode)
//	helper.Must(err)
//	str := res.(string)
//	authCodeInfo := &user.UcAuthInfo{}
//	json.Unmarshal([]byte(str), authCodeInfo)
//	if authCodeInfo.EnterpriseInfo.EpId == "" {
//		workWechatLoginResponse.Code = "auth_invalid"
//		workWechatLoginResponse.Desc = "授权无效"
//	}
//	return authCodeInfo, nil
// }

// func enterpriseDetail(wxMpAppId string, workWechatLoginRequest *user.WorkWechatLoginRequest, workWechatLoginResponse *entity.WorkWechatLoginResponse) (*entity.BossEnterpriseDetailResponse, error) {
//	res, err := api.EnterpriseDetail(wxMpAppId)
//	helper.Must(err)
//	detailStr := res.(string)
//	bossEnterpriseDetailRes := &entity.BossEnterpriseDetailResponse{}
//	json.Unmarshal([]byte(detailStr), bossEnterpriseDetailRes)
//	// boss 企业信息不存在
//	if bossEnterpriseDetailRes.ID != 0 {
//		if workWechatLoginRequest.RegisterEnterprise != "" {
//			// 注册企业信息,再返回
//		} else {
//			workWechatLoginResponse.Code = "enterprise_not_found"
//			workWechatLoginResponse.Desc = "需要注册企业"
//			workWechatLoginResponse.AuthEnterprise = bossEnterpriseDetailRes
//			return bossEnterpriseDetailRes, nil
//		}
//		// 检查是否过期
//		if true == checkExpire(bossEnterpriseDetailRes) {
//			workWechatLoginResponse.Code = "enterprise_expired"
//			workWechatLoginResponse.Desc = "企业已过试用期"
//			return bossEnterpriseDetailRes, nil
//		}
//	}
//	return bossEnterpriseDetailRes, nil
// }
//
// func checkExpire(bossEnterpriseDetailRes *entity.BossEnterpriseDetailResponse) bool {
//	if bossEnterpriseDetailRes.ID == 0 {
//		return false
//	}
//	nowTime := time.Now().Unix()
//	if bossEnterpriseDetailRes.VersionTag == "试用版" && (bossEnterpriseDetailRes.AccountStartAt >= nowTime || bossEnterpriseDetailRes.AccountExpireAt <= nowTime) {
//		return true
//	}
//	return false
// }
