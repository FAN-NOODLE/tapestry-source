package sdk

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/clients"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/helper"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"gw-work-scrm/internal/config"
	"gw-work-scrm/internal/pkg/errs"
	"net/http"
	"strconv"
	"strings"
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

// CheckResponseOld 校验返回数据结构
func CheckResponseOld(response CommonResponse) error {
	if response.Code == 200 {
		return nil
	}
	return errorx.NewErrorMsg(response.Code, response.Msg)
}

// Request sdk封装发送POST请求接口
func Request(c ctx.Context, method string, url string, body interface{}) (interface{}, error) {
	log.Info(c).Msgf("request sdk, url:%s, request:%+v", url, body)

	var enterpriseHash string

	var options []clients.Option
	options = append(options, clients.JSONBody(body))
	options = append(options, clients.Header(HeaderKeyEnum.XRequestID, c.ID()))

	if enterpriseHash != "" {
		options = append(options, clients.Header(HeaderKeyEnum.EnterpriseHash, enterpriseHash))
	}
	var res *clients.ResponseProxy
	var err error

	switch strings.ToUpper(method) {
	case "POST":
		res, err = ClientObject.Post(c, url, options...)
	case "GET":
		res, err = ClientObject.Get(c, url, options...)
	}
	if err != nil {
		return nil, err
	}

	if res.Response().StatusCode != 200 {
		return nil, errorx.New(fmt.Sprintf("接口请求错误：%s", res.Response().Status))
	}

	resStr, err := res.Read(clients.Text)
	helper.Must(err)
	log.Info(c).Msgf("sdk response, url:%s, response:%s", url, resStr.(string))

	var response CommonResponse
	helper.Must(json.Unmarshal([]byte(resStr.(string)), &response))
	err = CheckResponseOld(response)
	if err != nil {
		log.Info(c).Msgf("sdk response error, err:%+v", err)
		return nil, err
	}
	return response.Data, nil
}

// SendPost sdk封装发送POST请求接口
func SendPost(c ctx.Context, url string, body interface{}) (interface{}, error) {
	return Request(c, "POST", url, body)
}

// SendGet GET 请求
func SendGet(c ctx.Context, url string, body interface{}) (interface{}, error) {
	return Request(c, "GET", url, body)
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

// Post 发送Post请求
func (c *Client) Post(ctx ctx.Context, url string, body interface{}, response Response, options ...clients.Option) error {
	log.Info(ctx).Msgf("request interface, url:%s, request:%+v", url, body)
	if body != nil {
		options = append(options, clients.JSONBody(body))
	}
	res, err := ClientObject.Post(ctx, url, options...)
	if err != nil {
		log.Info(ctx).Msgf("接口请求失败：%s, url:%s", err.Error(), url)
		return errorx.NewErrorMsg(errs.ErrorStatus, fmt.Sprintf("接口请求失败：%s", err.Error()))
	}
	if res.Response().StatusCode != 200 {
		log.Info(ctx).Msgf("接口请求错误：%s, url:%s", res.Response().Status, url)
		return errorx.NewErrorMsg(errs.ErrorStatus, fmt.Sprintf("接口请求错误：%s", res.Response().Status))
	}

	resStr, err := res.Read(clients.Text)
	if err != nil {
		log.Info(ctx).Msgf("接口返回数据获取失败：%s, url:%s", err.Error(), url)
		return errorx.NewErrorMsg(errs.ErrorStatus, fmt.Sprintf("接口返回数据获取失败：%s, url:%s", err.Error(), url))
	}
	log.Info(ctx).Msgf("request cdp interface, url:%s, response:%s", url, resStr.(string))

	err = c.parse.Parse([]byte(resStr.(string)), response)
	if err != nil {
		log.Info(ctx).Msgf("接口返回数据解析失败：%s, url:%s", err.Error(), url)
		return errorx.NewErrorMsg(errs.ErrorStatus, fmt.Sprintf("接口返回数据解析失败：%s, url:%s", err.Error(), url))
	}
	if response != nil {
		err = response.CheckCode()
	}
	if err != nil {
		log.Info(ctx).Msgf("接口数据错误：%+v", err)
		return err
	}
	return nil
}
