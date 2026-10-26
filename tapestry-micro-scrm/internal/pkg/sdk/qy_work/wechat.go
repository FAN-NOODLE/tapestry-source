package qyWork

import (
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
	"os"
	"time"
)

const (
	// CommunalContact Contact 通讯录标识
	CommunalContact    string = "communalContact"
	Contact            string = "Contact"
	storeAssistant     string = "crs_store_assistant"
	StoreAssistantSelf string = "crs_store_assistant_self"
	// ManuCustomer 外部联系人标识
	ManuCustomer    string = "crs-manu-customer"
	VM              string = "vm"
	Notification    string = "notification"     //传讯
	OnlineWareHouse string = "online-warehouse" //网仓订单
	StoreVisit      string = "store-visit"      //巡店

)

var wechatDomain = "https://qyapi.weixin.qq.com"

type QyWork struct {
	enterpriseId int64
}

// NewQyWork 实例化授权信息
func NewQyWork(enterpriseId int64) *QyWork {
	return &QyWork{
		enterpriseId: enterpriseId,
	}
}

func (w *QyWork) getBaseUrl(path string) string {
	return fmt.Sprintf("%s/%s", wechatDomain, path)
}

func (w *QyWork) getTokenUrl(path, token string) string {
	return fmt.Sprintf("%s/%s?access_token=%s", wechatDomain, path, token)
}

// QyResponse 企业微信返回结果
type QyResponse struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

// requestPost 请求企业微信
func (w *QyWork) requestPost(ctx ctx.Context, url string, req any, resp any, opts ...clients.Option) error {
	// 请求参数
	body := helper.MarshalToString(req)
	// 返回值
	var respStr string
	// 企微提示消息标题
	var messageTitle string

	opts = append(opts, clients.Header(sdk.HeaderKeyEnum.XRequestID, ctx.ID()))
	client := http.NewHTTPClient(clients.Timeout(50 * time.Second))
	res, err := client.Post(ctx, url, opts...)
	var requestErr error
	if err != nil {
		messageTitle = fmt.Sprintf("企微-接口请求失败")
		log.Error(ctx).Msgf("企微-接口请求失败：%+v, url:%s", err, url)
		requestErr = errorx.NewErrorMsg(res.Response().StatusCode, fmt.Sprintf("企微-接口请求失败：%+v", err))
	} else {
		if res.Response().StatusCode != 200 {
			messageTitle = fmt.Sprintf("企微-接口请求状态码异常：状态码:%d", res.Response().StatusCode)
			log.Error(ctx).Msgf("企微-接口请求状态码异常：%s, url:%s", res.Response().Status, url)
			requestErr = errorx.NewErrorMsg(res.Response().StatusCode, fmt.Sprintf("企微-接口请求状态码异常：%s", res.Response().Status))
		} else {
			rt, rErr := res.Read(clients.Text)
			if rErr != nil {
				messageTitle = fmt.Sprintf("企微-接口返回数据获取失败")
				log.Error(ctx).Msgf("企微-接口返回数据获取失败：%s, url:%s", res.Response().Status, url)
				requestErr = errorx.NewErrorMsg(res.Response().StatusCode, fmt.Sprintf("企微-接口返回数据获取失败：%+v", rErr))
			} else {
				log.Info(ctx).Msgf("企微, 地址：%s, 参数：%s, 返回结果：%s", url, body, rt.(string))

				respStr = rt.(string)
				var qyResp QyResponse
				uIntErr := json.Unmarshal([]byte(respStr), &qyResp)
				if uIntErr != nil {
					log.Error(ctx).Msgf("企微-返回数据解析失败：%+v, url:%s", err, url)
					messageTitle = fmt.Sprintf("企微-返回数据解析失败")
					requestErr = errorx.NewErrorMsg(res.Response().StatusCode, fmt.Sprintf("企微-返回数据解析失败：%+v", uIntErr))
				} else {
					if qyResp.ErrCode != 0 {
						log.Error(ctx).Msgf("企微-状态码异常：状态码:%d, 信息:%s", qyResp.ErrCode, qyResp.ErrMsg)
						messageTitle = fmt.Sprintf("企微-状态码异常：状态码:%d", qyResp.ErrCode)
						requestErr = errorx.NewErrorMsg(qyResp.ErrCode, qyResp.ErrMsg)
					}
				}
			}
		}

	}

	if requestErr != nil {
		resErr := requestErr.(*errorx.Error)
		log.Error(ctx).Msgf("%s：%+v", resErr.Error(), resErr)
		// 41051 客户和服务人员已经开始聊天了
		// 40098 成员尚未实名认证
		// 84061 不存在外部联系人的关系
		// 41096 欢迎语已在发送中
		if helper.InArray(resErr.Code, []int{41051, 40098, 84061, 41096}) {
			return requestErr
		}
		go sdk.SendMessage(ctx, messageTitle, url, body, respStr, resErr)
		return requestErr
	}

	if resp != nil {
		jErr := json.Unmarshal([]byte(respStr), resp)
		if jErr != nil {
			log.Error(ctx).Msgf("企微-自定义返回数据解析失败：%+v, url:%s", jErr, url)
			go sdk.SendMessage(ctx, "企微-自定义返回数据解析失败", url, body, respStr, jErr)
			return errorx.NewErrorMsg(errs.ErrorStatus, "企微-自定义返回数据解析失败")
		}
	}

	return nil
}

// json请求
func (w *QyWork) requestPostJson(ctx ctx.Context, url string, req any, resp any, opts ...clients.Option) error {
	if req != nil {
		opts = append(opts, clients.JSONBody(req))
	}

	return w.requestPost(ctx, url, req, resp, opts...)
}

// form表单
func (w *QyWork) requestPostForm(ctx ctx.Context, url string, req map[string]string, files []*os.File, resp any, opts ...clients.Option) error {
	if req != nil {
		opts = append(opts, clients.Form(req, files...))
	}

	return w.requestPost(ctx, url, req, resp, opts...)
}
