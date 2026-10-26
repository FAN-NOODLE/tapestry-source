package uc

import (
	"encoding/json"
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/clients"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"micro-scrm/internal/pkg/errs"
	"micro-scrm/internal/pkg/http"
	"micro-scrm/internal/pkg/sdk"
)

var urlEnum = struct {
	MemberByUserId string
}{
	// 成员详情
	MemberByUserId: "%s/a/%s/%s/qy/member/get-by-userid",
}

// ResponseBase 服务接口返回数据结构体头
type ResponseBase struct {
	Code      string `json:"code"`
	Msg       string `json:"msg"`
	RequestId string `json:"requestId"`
}

// Post 发送Post请求
func Post(c ctx.Context, url string, request interface{}, response interface{}, options ...clients.Option) error {
	log.Info(c).Msgf("请求UC接口：%s，请求参数：%+v", url, request)
	if request != nil {
		options = append(options, clients.JSONBody(request))
	}
	options = append(options, clients.Header(sdk.HeaderKeyEnum.XRequestID, c.ID()))
	res, err := http.ClientObject.Post(c, url, options...)
	if err != nil {
		log.Info(c).Msgf("UC接口请求失败：%s, url:%s", err.Error(), url)
		return errorx.NewErrorMsg(errs.ErrorStatus, fmt.Sprintf("接口请求失败：%s", err.Error()))
	}

	resStr, err := res.Read(clients.Text)
	if err != nil {
		log.Info(c).Msgf("UC 接口请求失败，错误信息：%+v", err)
		return err
	}
	responseData := new(ResponseBase)
	byteStr := []byte(resStr.(string))
	jErr := json.Unmarshal(byteStr, responseData)
	if jErr != nil {
		return jErr
	}
	if responseData.Code != "SUCCESS" {
		return errorx.NewErrorMsg(100000, responseData.Msg)
	}
	return json.Unmarshal(byteStr, &response)
}
