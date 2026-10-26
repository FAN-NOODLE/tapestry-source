package sdk

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/clients"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"micro-scrm/internal/config"
	internalHelper "micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/notice"
	"net/http"
	"strconv"
	"time"
)

// ClientObject 请求
var ClientObject = NewSdkClient()

// ContextKeyEnum ctx.Context中存储的数据key枚举
var ContextKeyEnum = struct {
	EnterpriseHash string
}{
	EnterpriseHash: "enterprise_hash",
}

// HeaderKeyEnum 请求Header中存储的key枚举
var HeaderKeyEnum = struct {
	EnterpriseHash string
	XRequestID     string
}{
	EnterpriseHash: "Enterprise-Hash",
	XRequestID:     "X-Request-ID",
}

// CommonResponse 基础返回参数
type CommonResponse struct {
	App       string      `json:"app"`
	Code      int         `json:"code"`
	Msg       string      `json:"message"`
	RequestID string      `json:"request_id"`
	Data      interface{} `json:"data"`
}

// DefaultClient 调用内部服务接口Client
type DefaultClient struct {
	Client *clients.HTTPClient
}

// HeaderAuthSign 设置Header签名
func HeaderAuthSign() clients.Option {
	return func(request *http.Request) {
		// todo key配置
		key := "REDACTED_FOR_SOURCE_SHARING"
		client := config.CfgData.App
		curTime := strconv.FormatInt(time.Now().Unix(), 10)
		m := md5.New()
		m.Write([]byte(client + key + curTime))
		signStr := fmt.Sprintf("%s", hex.EncodeToString(m.Sum(nil)))
		request.Header.Set("Auth-Client", client)
		request.Header.Set("Auth-Sign", signStr)
		request.Header.Set("Auth-Timestamp", curTime)
	}
}

// HeaderSign 设置请求sdk接口的Header签名
var HeaderSign = HeaderAuthSign()

// NewSdkClient http client
func NewSdkClient(options ...clients.HTTPClientOption) *DefaultClient {
	var client DefaultClient
	client.Client = clients.NewHTTPClient(options...)
	return &client
}

// Post 发送 Post 请求
func (s *DefaultClient) Post(c ctx.Context, url string, options ...clients.Option) (*clients.ResponseProxy, error) {
	options = append(options, HeaderSign)
	return s.Client.Post(c, url, options...)
}

// Get 发送 Get 请求
func (s *DefaultClient) Get(c ctx.Context, url string, options ...clients.Option) (*clients.ResponseProxy, error) {
	options = append(options, HeaderSign)
	return s.Client.Get(c, url, options...)
}

// ParseResponse 解析response接口
type ParseResponse interface {
	Parse([]byte, interface{}) error
}

// ParseJSON json 格式
type ParseJSON struct {
}

// Parse 转换
func (p ParseJSON) Parse(bytes []byte, response interface{}) error {
	return json.Unmarshal(bytes, response)
}

// Client 请求代理结构体
type Client struct {
	parse ParseResponse
}

// NewHTTPClient 创建Client对象
func NewHTTPClient(parse ParseResponse) *Client {
	var client = Client{
		parse: parse,
	}
	return &client
}

// Response 返回结果
type Response interface {
	CheckResponse
}

// CheckResponse 检测返回结果
type CheckResponse interface {
	CheckCode() error
}

// ResponseBase 返回结果定义
type ResponseBase struct {
	App       string `json:"app"`
	Code      int    `json:"code"`
	Msg       string `json:"msg"`
	RequestID string `json:"request_id"`
}

// CheckCode 检测错误码
func (r *ResponseBase) CheckCode() error {
	if r.Code == 200 {
		return nil
	}
	return errorx.NewErrorMsg(r.Code, r.Msg)
}

func SendMessage(ctx ctx.Context, title, url string, body, resp interface{}, err error) {
	content := fmt.Sprintf("%s \n"+
		">请求时间：<font color='comment'>%s</font>\n"+
		">错误信息：<font color='comment'>%v</font>\n"+
		">请求地址：<font color='comment'>%s</font>\n"+
		">请求参数：<font color='comment'>%v</font>\n"+
		">返回信息：<font color='comment'>%v</font>\n"+
		">RequestId：<font color='comment'>%v</font>\n",
		title,
		internalHelper.GetFormatDateByParams(0),
		err,
		url,
		internalHelper.MarshalToString(body),
		internalHelper.MarshalToString(resp),
		ctx.ID(),
	)
	if len(content) > 2048 {
		content = content[:2048]
	}
	notice.NewWechat().SendMarkdownMessage(ctx, content)
}

// Post 发送Post请求
func (c *Client) Post(ctx ctx.Context, url string, req interface{}, response Response, options ...clients.Option) error {
	log.Info(ctx).Msgf("请求微服务, url:%s, request:%+v", url, req)
	// 请求参数
	body := internalHelper.MarshalToString(req)
	// 返回值
	var respStr string
	// 企微提示消息标题
	var messageTitle string

	err := internalHelper.WithRetry(3, time.Second, func() error {
		if req != nil {
			options = append(options, clients.JSONBody(req))
		}

		client := NewSdkClient(clients.Timeout(50 * time.Second))
		res, err := client.Post(ctx, url, options...)

		if err != nil {
			messageTitle = fmt.Sprintf("微服务-接口请求失败")
			log.Error(ctx).Msgf("微服务-接口请求失败：%+v, url:%s", err, url)
			return errorx.NewErrorMsg(res.Response().StatusCode, fmt.Sprintf("微服务-接口请求失败：%+v", err))
		}
		if res.Response().StatusCode != 200 {
			messageTitle = fmt.Sprintf("微服务-接口请求状态码异常：状态码:%d", res.Response().StatusCode)
			log.Error(ctx).Msgf("微服务-接口请求状态码异常：%s, url:%s", res.Response().Status, url)
			return errorx.NewErrorMsg(res.Response().StatusCode, fmt.Sprintf("微服务-接口请求状态码异常：%s", res.Response().Status))
		}

		rt, err := res.Read(clients.Text)
		if err != nil {
			messageTitle = fmt.Sprintf("微服务-接口返回数据获取失败")
			log.Error(ctx).Msgf("微服务-接口返回数据获取失败：%+v, url:%s", err, url)
			return errorx.NewErrorMsg(res.Response().StatusCode, fmt.Sprintf("微服务-接口返回数据获取失败：%+v, url:%s", err, url))
		}

		log.Info(ctx).Msgf("微服务, 地址：%s, 参数：%s, 返回结果：%s", url, body, rt.(string))

		respStr = rt.(string)

		err = c.parse.Parse([]byte(respStr), response)
		if err != nil {
			messageTitle = fmt.Sprintf("微服务-接口返回数据解析失败")
			log.Error(ctx).Msgf("微服务-接口返回数据解析失败：%+v, url:%s", err, url)
			return errorx.NewErrorMsg(res.Response().StatusCode, fmt.Sprintf("微服务-接口返回数据解析失败：%+v, url:%s", err, url))
		}

		return nil
	})

	if err != nil {
		resErr := err.(*errorx.Error)
		log.Error(ctx).Msgf("%s：%+v", resErr.Error(), resErr)
		go SendMessage(ctx, messageTitle, url, body, respStr, resErr)
		return err
	}

	return nil
}
