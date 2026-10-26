package gwshop

import (
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
	"gw-scrm/internal/config"
	"gw-scrm/internal/pkg/sdk"
)

var urlEnum = struct {
	GetAccessToken string
}{
	// 获取短链
	GetAccessToken: "%s/gw-shop/app/v1/mini-program/get-access-token",
}

// 获取cdp请求链接
func GetUrl(uri string) string {
	return fmt.Sprintf(uri, config.CfgData.Service.ShopToken.Domain)
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
