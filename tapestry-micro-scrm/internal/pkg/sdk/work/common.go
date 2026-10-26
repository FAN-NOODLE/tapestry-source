package work

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/clients"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"micro-scrm/internal/pkg/errs"
	"micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/http"
	"micro-scrm/internal/pkg/sdk"
	"time"
)

// 请求方式常量
const (
	GET  = "GET"
	POST = "POST"
)

func WorkRequest(ctx ctx.Context, url string, req interface{}, response interface{}, options ...clients.Option) error {
	return workRequest(ctx, url, req, response, options...)
}

// workRequest 请求企业微信
func workRequest(ctx ctx.Context, url string, req interface{}, response interface{}, options ...clients.Option) error {
	// 请求参数
	body := helper.MarshalToString(req)
	// 返回值
	var respStr string
	// 企微提示消息标题
	var messageTitle string

	err := helper.WithRetry(2, time.Second, func() error {
		if req != nil {
			options = append(options, clients.JSONBody(req))
		}
		options = append(options, clients.Header(sdk.HeaderKeyEnum.XRequestID, ctx.ID()))
		client := http.NewHTTPClient(clients.Timeout(50 * time.Second))
		res, err := client.Post(ctx, url, options...)
		if err != nil {
			messageTitle = fmt.Sprintf("企微-接口请求失败")
			log.Error(ctx).Msgf("企微-接口请求失败：%+v, url:%s", err, url)
			return errorx.NewErrorMsg(res.Response().StatusCode, fmt.Sprintf("企微-接口请求失败：%+v", err))
		}

		if res.Response().StatusCode != 200 {
			messageTitle = fmt.Sprintf("企微-接口请求状态码异常：状态码:%d", res.Response().StatusCode)
			log.Error(ctx).Msgf("企微-接口请求状态码异常：%s, url:%s", res.Response().Status, url)
			return errorx.NewErrorMsg(res.Response().StatusCode, fmt.Sprintf("企微-接口请求状态码异常：%s", res.Response().Status))
		}

		rt, rErr := res.Read(clients.Text)
		if rErr != nil {
			messageTitle = fmt.Sprintf("企微-接口返回数据获取失败")
			log.Error(ctx).Msgf("企微-接口返回数据获取失败：%s, url:%s", res.Response().Status, url)
			return errorx.NewErrorMsg(res.Response().StatusCode, fmt.Sprintf("企微-接口返回数据获取失败：%+v", rErr))
		}

		log.Info(ctx).Msgf("企微, 地址：%s, 参数：%s, 返回结果：%s", url, body, rt.(string))

		respStr = rt.(string)
		responseInt := QyResponse{}
		uIntErr := json.Unmarshal([]byte(respStr), &responseInt)
		if uIntErr != nil {
			log.Error(ctx).Msgf("企微-返回数据解析失败：%+v, url:%s", err, url)
			messageTitle = fmt.Sprintf("企微-返回数据解析失败")
			return errorx.NewErrorMsg(res.Response().StatusCode, fmt.Sprintf("企微-返回数据解析失败：%+v", uIntErr))
		}

		if responseInt.ErrCode != 0 {
			log.Error(ctx).Msgf("企微-状态码异常：状态码:%d, 信息:%s", responseInt.ErrCode, responseInt.ErrMsg)
			messageTitle = fmt.Sprintf("企微-状态码异常：状态码:%d", responseInt.ErrCode)
			return errorx.NewErrorMsg(responseInt.ErrCode, responseInt.ErrMsg)
		}
		return nil
	})

	if err != nil {
		resErr := err.(*errorx.Error)
		log.Error(ctx).Msgf("%s：%+v", resErr.Error(), resErr)
		// 41051 客户和服务人员已经开始聊天了
		// 40098 成员尚未实名认证
		// 84061 不存在外部联系人的关系
		// 41096 欢迎语已在发送中
		// 60011 导购已离职
		// 41093 群发消息已取消
		if helper.InArray(resErr.Code, []int{41051, 40098, 84061, 41096, 60011, 41093}) {
			return err
		}
		go sdk.SendMessage(ctx, messageTitle, url, body, respStr, resErr)
		return err
	}

	if response != nil {
		jErr := json.Unmarshal([]byte(respStr), response)
		if jErr != nil {
			log.Error(ctx).Msgf("企微-自定义返回数据解析失败：%+v, url:%s", jErr, url)
			go sdk.SendMessage(ctx, "企微-自定义返回数据解析失败", url, body, respStr, jErr)
			return errorx.NewErrorMsg(errs.ErrorStatus, "企微-自定义返回数据解析失败")
		}
	}

	return nil
}

func Request(ctx ctx.Context, url string, req interface{}, options ...clients.Option) (string, error) {
	options = append(options, clients.JSONBody(req))
	options = append(options, clients.Header(sdk.HeaderKeyEnum.XRequestID, ctx.ID()))

	var result string

	client := http.NewHTTPClient(clients.Timeout(60 * time.Second))
	res, err := client.Post(ctx, url, options...)
	if err != nil {
		return result, err
	}

	rt, rErr := res.Read(clients.Text)
	if rErr != nil {
		return result, rErr
	}

	result = rt.(string)

	body, _ := json.Marshal(req)
	log.Info(ctx).Msgf("请求地址：%s, 请求参数：%s, 返回结果：%s", url, string(body), result)

	var responseInt QyResponse
	uIntErr := json.Unmarshal([]byte(result), &responseInt)
	if uIntErr != nil {
		return result, errs.ResponseNotJson
	}

	if responseInt.ErrCode != 0 {
		return result, errorx.NewErrorMsg(responseInt.ErrCode, responseInt.ErrMsg)
	}

	return result, nil
}

// js config 签名
func JsSign(nonce, ticket, url string, timestamp int64) string {
	demo := "jsapi_ticket=%s&noncestr=%s&timestamp=%d&url=%s"
	string := fmt.Sprintf(demo, ticket, nonce, timestamp, url)
	h := sha1.New()
	h.Write([]byte(string))
	l := fmt.Sprintf("%x", h.Sum(nil))
	return l
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
	responseData := new(SBaseResponse)
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
