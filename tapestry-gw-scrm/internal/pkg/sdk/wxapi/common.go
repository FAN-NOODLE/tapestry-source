package wxapi

import (
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
	"gw-scrm/internal/config"
	"gw-scrm/internal/pkg/sdk"
)

var urlEnum = struct {
	GetWxUrl string
}{
	// 获取微信跳转链接
	GetWxUrl: "%s/wxa/generatescheme?access_token=%s",
}

// 获取cdp请求链接
func GetUrl(uri string, accessToken interface{}) string {
	return fmt.Sprintf(uri, config.CfgData.Service.WxApiUrl.Domain, accessToken)
}

var httpClient = sdk.NewHttpClient(sdk.ParseJson{})

type GwShopResponse struct {
	App       string      `json:"app"`
	Code      int         `json:"code"`
	Msg       string      `json:"message"`
	RequestId string      `json:"request_id"`
	Data      interface{} `json:"data"`
}

func (r *GwShopResponse) CheckCode() error {
	if r.Code == 200 {
		return nil
	}
	return errorx.NewErrorMsg(r.Code, r.Msg)

}

// CommonResponse 公共默认返回值
type CommonResponse struct {
	GwShopResponse
	Data interface{} `json:"data"`
}

var scrmClient = sdk.NewHttpClient(sdk.ParseJson{})
