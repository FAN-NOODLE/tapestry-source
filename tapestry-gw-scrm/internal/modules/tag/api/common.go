package api

import (
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/clients"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/helper"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	scrm_errs "gw-scrm/internal/pkg/errs"
	"gw-scrm/internal/pkg/http"
)

// ResponseBase 统一返回结构体
type ResponseBase struct {
	App       string `json:"app"`
	Code      int    `json:"code"`
	Msg       string `json:"msg"`
	RequestId string `json:"request_id"`
}

// GetUrl 请求链接
func GetUrl(uri string, domain string) string {
	return fmt.Sprintf(uri, domain)
}

//  CheckResponse 接口返回结果
func CheckResponse(response ResponseBase) error {
	if response.Code == 0 {
		return nil
	}

	return errorx.NewErrorMsg(response.Code, response.Msg)
}

// SendPost 发送POST请求接口
func SendPost(c ctx.Context, uri string, domain string, body interface{}) string {
	url := GetUrl(uri, domain)

	res, err := http.ClientObject.Post(c, url, clients.JSONBody(body))
	helper.Must(err)
	if res.Response().StatusCode != 200 {
		panic(errorx.NewErrorMsg(scrm_errs.DefaultErrorCode, fmt.Sprintf("接口请求错误：%s", res.Response().Status)))
	}
	resStr, err := res.Read(clients.Text)
	helper.Must(err)
	log.Info(c).Msgf("request interface, url:%s, response:%s", url, resStr.(string))
	return resStr.(string)
}
