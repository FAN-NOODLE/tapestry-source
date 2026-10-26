package microauthority

import (
	"encoding/json"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/pkg/sdk"
)

var urlEnum = struct {
	DecryptToken    string // 解密token接口
	BaseConfigCheck string
	CheckAccount    string
	EnterpriseList  string
	AccessControl   string
	LoginByMobile   string
}{
	DecryptToken:    "%s/micro-authority/v1/account/decrypt-token",
	BaseConfigCheck: "%s/micro-authority/v1/account/base-info-check",
	CheckAccount:    "%s/micro-authority/v1/account/check-account",
	EnterpriseList:  "%s/micro-authority/v1/account/account-enterprise-list",
	AccessControl:   "%s/micro-authority/v1/account/account-access-control",
	LoginByMobile:   "%s/micro-authority/v1/account/sso-login",
}

// 解密token接口调用
func DecryptToken(c ctx.Context, request interface{}) (interface{}, error) {
	response, err := sdk.SendPost(c, GetUrl(urlEnum.DecryptToken), request)
	return response, err
}

func BaseInfoCheck(c ctx.Context, request interface{}) (interface{}, error) {
	response, err := sdk.SendPost(c, GetUrl(urlEnum.BaseConfigCheck), request)
	return response, err
}

func CheckAccount(c ctx.Context, request interface{}) (interface{}, error) {
	response, err := sdk.SendPost(c, GetUrl(urlEnum.CheckAccount), request)
	return response, err
}

func LoginByMobile(c ctx.Context, request interface{}) (interface{}, error) {
	response, err := sdk.SendPost(c, GetUrl(urlEnum.LoginByMobile), request)
	return response, err
}

func EnterpriseList(c ctx.Context, request interface{}) (interface{}, error) {
	response, err := sdk.SendPost(c, GetUrl(urlEnum.EnterpriseList), request)
	return response, err
}

type AccessControlReq struct {
	AdminID      int64  `json:"admin_id"`
	EnterpriseID int64  `json:"enterprise_id"`
	URL          string `json:"url"`
	AppCode      string `json:"app_code"`
}

type AccessControlRep struct {
	Access bool `json:"access"`
}

func AccessControl(c ctx.Context, request interface{}) (bool, error) {
	response, err := sdk.SendPost(c, GetUrl(urlEnum.AccessControl), request)
	if err != nil {
		return false, err
	}

	arr, err := json.Marshal(response)
	if err != nil {
		return false, err
	}
	result := new(AccessControlRep)
	err = json.Unmarshal(arr, result)
	if err != nil {
		return false, err
	}

	return result.Access, nil
}
