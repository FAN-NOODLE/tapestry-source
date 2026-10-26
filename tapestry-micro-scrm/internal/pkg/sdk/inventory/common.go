package inventory

import (
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/clients"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"micro-scrm/internal/config"
	"micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/notice"
	"time"
)

const (
	DESC = 2
)

func reportInventoryError(c ctx.Context, url, reqStr, respStr string, err error) {
	log.Error(c).Msgf("三方库存接口请求异常：%+v", err)
	if !IsContractError(err) {
		return
	}
	content := fmt.Sprintf("库存接口数据契约异常\n请求地址：%s\n请求参数：%s\n返回信息：%s\n错误信息：%v\nRequestId：%s", url, reqStr, respStr, err, c.ID())
	if len(content) > 2048 {
		content = content[:2048]
	}
	go notice.NewWechat().SendTextMessage(c, content)
}

func Post(ctx ctx.Context, isDcl bool, request ParamsReq, resp *Response, opts ...clients.Option) error {
	// 获取url
	conf := config.CfgData.Inventory
	url := conf.Url
	if isDcl {
		url = conf.DclUrl
	}
	// 请求参数
	request.Country = "CN"
	request.SortOption.OnHand = DESC
	body := clients.JSONBody(request)
	opts = append(opts, body)
	// 返回值
	var respStr string
	reqStr := helper.MarshalToString(request)
	if handled, mockResp, mockErr := tryInventoryMock(ctx, isDcl, request, resp); handled {
		if mockErr != nil {
			reportInventoryError(ctx, url, reqStr, mockResp, mockErr)
		}
		return mockErr
	}
	err := helper.WithRetry(1, time.Second, func() error {
		opts = append(opts, clients.Header("X-IBM-Client-Id", conf.ClientId))
		client := clients.NewHTTPClient(clients.Timeout(120 * time.Second))
		res, err := client.Post(ctx, url, opts...)
		if err != nil {
			return errorx.NewErrorMsg(500, fmt.Sprintf("三方库存接口-接口请求失败：%+v，res:%v", err, res))
		}

		if res.Response().StatusCode != 200 {
			return errorx.NewErrorMsg(res.Response().StatusCode, fmt.Sprintf("三方库存接口-接口请求状态码异常：%s", res.Response().Status))
		}

		rt, rErr := res.Read(clients.Text)
		if rErr != nil {
			return errorx.NewErrorMsg(res.Response().StatusCode, fmt.Sprintf("三方库存接口-接口返回数据获取失败：%+v", rErr))
		}

		log.Info(ctx).Msgf("三方库存接口, 地址：%s, 参数：%+v, 返回结果：%s", url, reqStr, rt.(string))

		respStr = rt.(string)
		if decodeErr := decodeAndValidateResponse(isDcl, request, respStr, resp); decodeErr != nil {
			return decodeErr
		}
		return nil
	})

	if err != nil {
		reportInventoryError(ctx, url, reqStr, respStr, err)
		return err
	}

	return nil
}

func PostOnce(ctx ctx.Context, isDcl bool, request ParamsReq, resp *Response, opts ...clients.Option) error {
	// 获取url
	conf := config.CfgData.Inventory
	url := conf.Url
	if isDcl {
		url = conf.DclUrl
	}
	// 请求参数
	request.Country = "CN"
	request.SortOption.OnHand = DESC
	body := clients.JSONBody(request)
	opts = append(opts, body)
	// 返回值
	var respStr string
	reqStr := helper.MarshalToString(request)
	if handled, mockResp, mockErr := tryInventoryMock(ctx, isDcl, request, resp); handled {
		if IsContractError(mockErr) {
			reportInventoryError(ctx, url, reqStr, mockResp, mockErr)
		}
		return mockErr
	}
	opts = append(opts, clients.Header("X-IBM-Client-Id", conf.ClientId))
	client := clients.NewHTTPClient(clients.Timeout(120 * time.Second))
	res, err := client.Post(ctx, url, opts...)
	if err != nil {
		return errorx.NewErrorMsg(500, fmt.Sprintf("三方库存接口-接口请求失败：%+v，res:%v", err, res))
	}

	if res.Response().StatusCode != 200 {
		return errorx.NewErrorMsg(res.Response().StatusCode, fmt.Sprintf("三方库存接口-接口请求状态码异常：%s", res.Response().Status))
	}

	rt, rErr := res.Read(clients.Text)
	if rErr != nil {
		return errorx.NewErrorMsg(res.Response().StatusCode, fmt.Sprintf("三方库存接口-接口返回数据获取失败：%+v", rErr))
	}

	log.Info(ctx).Msgf("三方库存接口, 地址：%s, 参数：%+v, 返回结果：%s", url, reqStr, rt.(string))

	respStr = rt.(string)
	if decodeErr := decodeAndValidateResponse(isDcl, request, respStr, resp); decodeErr != nil {
		reportInventoryError(ctx, url, reqStr, respStr, decodeErr)
		return decodeErr
	}
	return nil

}
