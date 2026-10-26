package http

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/clients"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
	"gw-work-scrm/internal/config"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

// ResponseBase 基础返回参数
type ResponseBase struct {
	App       string `json:"app"`
	Code      int    `json:"code"`
	Msg       string `json:"msg"`
	RequestID string `json:"request_id"`
}

// ClientObject http请求代理变量
var ClientObject = NewHTTPClient()

// CdpClientObject cdp请求代理变量
var CdpClientObject = NewCdpClient()

// SdkClient 调用内部服务接口Client
type SdkClient struct {
	Client *clients.HTTPClient
}

// CdpClient cdp接口调用Client
type CdpClient struct {
	Client *clients.HTTPClient
}

// HeaderAuthSign 设置Header签名
func HeaderAuthSign() clients.Option {
	return func(request *http.Request) {
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

// CdpHeaderAuthSign 设置Header签名
func CdpHeaderAuthSign(c ctx.Context) clients.Option {
	return func(request *http.Request) {
		var enterpriseID int64
		if c.Value("enterprise_id") != nil {
			enterpriseID = c.Value("enterprise_id").(int64)
		}
		var enterpriseHash string
		if c.Value("enterprise_hash") != nil {
			enterpriseHash = c.Value("enterprise_hash").(string)
		}
		enterpriseHash = "ecd85280aa135bbd0108dd6aa424565a"
		enterpriseID = 2

		var bodyBytes []byte
		if request.Body != nil {
			bodyBytes, _ = ioutil.ReadAll(request.Body)
			request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		var bodyString string
		if bodyBytes != nil {
			bodyString = string(bodyBytes)
		}

		curTime := strconv.FormatInt(time.Now().Unix(), 10)

		m := md5.New()
		m.Write([]byte(bodyString + curTime + enterpriseHash))
		signStr := fmt.Sprintf("%s", hex.EncodeToString(m.Sum(nil)))

		fmt.Println(bodyString + curTime + enterpriseHash)
		fmt.Println(signStr)
		request.Header.Set("Enterprise-ID", strconv.FormatInt(enterpriseID, 10))
		request.Header.Set("Auth-Sign", signStr)
		request.Header.Set("Auth-Timestamp", curTime)
	}
}

// HeaderSign 设置请求sdk接口的Header签名
var HeaderSign = HeaderAuthSign()

// NewHTTPClient 获取http请求代理
func NewHTTPClient(options ...clients.HTTPClientOption) *SdkClient {
	var client SdkClient
	client.Client = clients.NewHTTPClient(options...)
	return &client
}

// NewCdpClient 获取cdp请求代理
func NewCdpClient(options ...clients.HTTPClientOption) *CdpClient {
	var client CdpClient
	client.Client = clients.NewHTTPClient(options...)
	return &client
}

// Post 发送 Post 请求
func (cc *CdpClient) Post(c ctx.Context, url string, options ...clients.Option) (*clients.ResponseProxy, error) {
	return cc.Client.Post(c, url, options...)
}

// Get 发送 Get 请求
func (cc *CdpClient) Get(c ctx.Context, url string, options ...clients.Option) (*clients.ResponseProxy, error) {
	return cc.Client.Get(c, url, options...)
}

// Post 发送 Post 请求
func (s *SdkClient) Post(c ctx.Context, url string, options ...clients.Option) (*clients.ResponseProxy, error) {
	options = append(options, HeaderSign)
	return s.Client.Post(c, url, options...)
}

// Get 发送 Get 请求
func (s *SdkClient) Get(c ctx.Context, url string, options ...clients.Option) (*clients.ResponseProxy, error) {
	options = append(options, HeaderSign)
	return s.Client.Get(c, url, options...)
}

// CheckResponse 校验返回数据结构
func CheckResponse(response ResponseBase) error {
	if response.Code != 0 {
		return nil
	}
	return errorx.NewErrorMsg(response.Code, response.Msg)
}
