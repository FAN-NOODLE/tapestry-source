package http

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"micro-scrm/internal/modules/enterprise"
	"micro-scrm/internal/orm"
	enterpriseConfigModel "micro-scrm/internal/orm/model/micro_scrm/enterprise"
	"micro-scrm/internal/orm/model/micro_scrm/fission"
	"micro-scrm/internal/pkg/common"
	"micro-scrm/internal/pkg/errs"
	"micro-scrm/internal/pkg/sdk/wechat/officialaccount"
	"micro-scrm/internal/pkg/uid"
	"net/url"
	"time"
)

func redirectURL(c *gin.Context) {
	request := enterprise.OfficialRedirectURLReq{}
	if err := c.BindJSON(&request); err != nil {
		ginplus.ResError(c, err)
		return
	}

	cc := ctx.Wrap(c)
	authService, serErr := NewAuthService(cc, request.EnterpriseID)
	if serErr != nil {
		ginplus.ResError(c, serErr)
		return
	}

	redirectURL, redErr := authService.RedirectURL(request)
	if redErr != nil {
		ginplus.ResError(c, redErr)
		return
	}

	ginplus.ResSuccess(c, enterprise.OfficialRedirectURLRes{URL: redirectURL})
}

// officialAuth 公众号授权登录
func officialAuth(c *gin.Context) {
	request := enterprise.OfficialAuthReq{}
	if err := c.BindJSON(&request); err != nil {
		ginplus.ResError(c, err)
		return
	}

	cc := ctx.Wrap(c)
	authService, serErr := NewAuthService(cc, request.EnterpriseID)
	if serErr != nil {
		ginplus.ResError(c, serErr)
		return
	}

	isExist, userInfo, userErr := authService.GetUserInfo(request)
	if userErr != nil {
		ginplus.ResError(c, userErr)
		return
	}

	if !isExist {
		// 用户不存在时，新增用户
		userInfo.ID = uid.ID(cc)
		userInfo.EnterpriseID = request.EnterpriseID
		userInfo.CreatedAt = time.Now().Unix()
		userInfo.UpdatedAt = time.Now().Unix()
		insertErr := orm.Service.Fission.InsertFissionMpUser(userInfo)
		if insertErr != nil {
			ginplus.ResError(c, insertErr)
			return
		}
	}

	// 还原原先的返回值，让前端不动
	ginplus.ResSuccess(c, enterprise.OfficialUser{
		ID:           userInfo.ID,
		OpenID:       userInfo.OpenID,
		Nickname:     userInfo.NickName,
		Sex:          userInfo.Sex,
		Province:     userInfo.Province,
		City:         userInfo.City,
		Country:      userInfo.Country,
		HeadImgURL:   userInfo.ImgURL,
		Privilege:    userInfo.Privilege,
		UnionID:      userInfo.UnionID,
		EnterpriseID: userInfo.EnterpriseID,
		CreatedAt:    userInfo.CreatedAt,
		UpdatedAt:    userInfo.UpdatedAt,
		DeletedAt:    userInfo.DeletedAt,
	})
	return
}

// 重构登陆授权这块，客户可能是两种配置
// 1、微信公众号
// 2、第三方平台
type configValue struct {
	AppID   string            `json:"app_id"`
	Secret  string            `json:"secret"`
	Encrypt map[string]string `json:"encrypt"`
	URLS    map[string]string `json:"urls"`
}

type authService struct {
	enterpriseID int64
	key          string
	config       enterpriseConfigModel.EnterpriseConfigModel
	value        configValue
	CC           ctx.Context
}

// 参与该功能的enterprise_config.key值集合
var (
	EnterpriseConfigKeyMap = map[string]bool{
		"official_info":   true,
		"component":       true,
		"component_no_js": true,
	}
)

const (
	OfficialInfo  = "official_info"
	Component     = "component"
	ComponentNoJs = "component_no_js"
	CYFlag        = "cy_flag"
)

func NewAuthService(cc ctx.Context, enterpriseID int64) (*authService, error) {
	// 根据企业ID获取认证的配置
	configList, queryErr := orm.Service.Enterprise.GetConfigListOneEnterprise(enterpriseID)
	if queryErr != nil {
		return nil, errs.SettingsInfoErr
	}

	// 根据配置，获取公众号/第三方配置。
	// 默认访问优先级：(期望)公众号 > 第三方配置，其实是根据数据库的数据排序优先级的。所以，两者只能使用其一。配置enterprise_config.is_use 0 未使用，1使用
	// enterprise_config.key 公众号：official_info 第三方方：component
	for _, config := range configList {
		if "" == config.Value {
			continue
		}

		if _, ok := EnterpriseConfigKeyMap[config.Key]; ok {
			conf := configValue{}
			jsonErr := json.Unmarshal([]byte(config.Value), &conf)
			if jsonErr != nil {
				return nil, jsonErr
			}

			if 1 == config.IsUse {
				return &authService{
					key:          config.Key,
					enterpriseID: enterpriseID,
					config:       config,
					value:        conf,
					CC:           cc,
				}, nil
			}
		}
	}

	return nil, errs.SettingsInfoErr
}

// 获取配置中的重定向连接
func (a *authService) RedirectURL(param enterprise.OfficialRedirectURLReq) (redirectUrl string, err error) {
	switch a.key {
	case OfficialInfo:
		oauth := officialaccount.NewOauth(&officialaccount.Config{
			AppID:     a.value.AppID,
			AppSecret: a.value.Secret,
		})

		newUrl, _ := a.URLAddParam(param.URL, CYFlag, OfficialInfo)
		return oauth.GetRedirectURL(newUrl, param.Scope, param.State)

	case Component, ComponentNoJs:
		lorealAuthUrl := a.value.URLS["auth"]
		newUrl, _ := a.URLAddParam(param.URL, CYFlag, Component)
		return fmt.Sprintf(lorealAuthUrl, a.value.AppID, param.Scope, url.QueryEscape(newUrl)), nil
	}

	err = errors.New("异常操作")
	return
}

func (a authService) URLAddParam(urlStr, key, value string) (string, error) {
	url, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}

	if "" == url.RawQuery {
		return fmt.Sprintf("%s?%s=%s", urlStr, key, value), nil
	} else {
		return fmt.Sprintf("%s&%s=%s", urlStr, key, value), nil
	}
}

// 获取用户基本信息, 并检测用户是否存在
func (a *authService) GetUserInfo(param enterprise.OfficialAuthReq) (isExist bool, data fission.FissionMpUserModel, err error) {
	switch a.key {
	case OfficialInfo:
		// 获取 access_token
		oauth := officialaccount.NewOauth(&officialaccount.Config{
			AppID:     a.value.AppID,
			AppSecret: a.value.Secret,
		})
		token, tokenErr := oauth.GetUserAccessToken(a.CC, param.Code)
		if tokenErr != nil {
			err = tokenErr
			return
		}
		// 查询该企业用户是否存在
		isExist, data, err = orm.Service.Fission.GetFissionMpUser(a.enterpriseID, token.OpenID)
		if err != nil {
			return
		}

		// 不存在时，需要获取
		if !isExist {
			userInfo, getUserErr := oauth.GetUserInfo(a.CC, token.AccessToken, token.OpenID, "zh_CN")
			if getUserErr != nil {
				err = getUserErr
				return
			}

			data.OpenID = userInfo.OpenID
			data.NickName = userInfo.Nickname
			data.Sex = userInfo.Sex
			data.Province = userInfo.Province
			data.City = userInfo.City
			data.Country = userInfo.Country
			data.ImgURL = userInfo.HeadImgURL
			privilegeListJson, _ := json.Marshal(userInfo.Privilege)
			data.Privilege = string(privilegeListJson)
			data.UnionID = userInfo.Unionid
		}

		return

	case Component, ComponentNoJs:
		// 如果是开放平台，咋用户信息放到code中
		encryptUserInfo, _ := hex.DecodeString(param.Code)
		userInfoStr := common.AES_CBC_Decrypt(encryptUserInfo, []byte(a.value.Encrypt["key"]), []byte(a.value.Encrypt["iv"]))

		// 原先公众号的用户信息，和欧莱雅第三方的返回值中的Privilege类型不一致，这块做一个兼容
		// 存库时，还是用原先的公众号的结构存储
		var userInfo enterprise.ThirdUserInfo
		err = json.Unmarshal(userInfoStr, &userInfo)
		if err != nil {
			return
		}

		// 查询该企业用户是否存在
		isExist, data, err = orm.Service.Fission.GetFissionMpUser(a.enterpriseID, userInfo.OpenID)
		if err != nil {
			return
		}

		data.OpenID = userInfo.OpenID
		data.NickName = userInfo.Nickname
		data.Sex = userInfo.Sex
		data.Province = userInfo.Province
		data.City = userInfo.City
		data.Country = userInfo.Country
		data.ImgURL = userInfo.HeadImgURL
		privilegeListJson, _ := json.Marshal(userInfo.Privilege)
		data.Privilege = string(privilegeListJson)
		data.UnionID = userInfo.UnionID

		return
	}

	err = errors.New("异常操作")
	return
}
