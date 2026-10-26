package common

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/clients"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/config"
	"net/http"
	"strconv"
	"time"
)

type ClientMode struct {
}

type GetClientUrl struct {
	Domain string
}

// 调用内部服务接口Client
type SdkClient struct {
	Client *clients.HTTPClient
}

/**
  返回调用系统微服务接口地址
*/
func (a ClientMode) GetClientUrl(domain GetClientUrl) string {
	switch domain.Domain {
	case "work":
		return config.CfgData.Service.Work.Domain
	case "shop":
		return config.CfgData.Service.Shop.Domain
	case "boss":
		return config.CfgData.Service.Boss.Domain
	case "console":
		return config.CfgData.Service.Console.Domain
	case "micro-pay":
		return config.CfgData.Service.MicroPay.Domain
	case "micro-activity":
		return config.CfgData.Service.MicroActivity.Domain
	case "micro-goods":
		return config.CfgData.Service.MicroGoods.Domain
	case "micro-authority":
		return config.CfgData.Service.MicroAuthority.Domain
	case "micro-user":
		return config.CfgData.Service.MicroUser.Domain
	}
	return ""
}

type GetUcUrl struct {
}

/**
  返回调用UC服务接口地址
*/
func (a ClientMode) GetUcUrl() string {

	mode := config.CfgData.Mode
	ucHost := ""

	if mode == "debug" {
		ucHost = "http://t-rst.vchangyi.com"
	}
	if mode == "test" {
		ucHost = "http://t-rst.vchangyi.com"
	}
	if mode == "release" {
		ucHost = "http://restapi.99hr.com"
	}
	if mode == "uat" {
		ucHost = "http://restapi.99hr.com"
	}
	if mode == "prod" {
		ucHost = "http://restapi.99hr.com"
	}

	return ucHost
}

/**
  返回智慧导购安装地址
*/
func (a ClientMode) GetAppSuiteId() string {

	mode := config.CfgData.Mode
	suiteId := ""

	if mode == "debug" {
		suiteId = "wx4e4c1f741ed9e582"
	}
	if mode == "test" {
		suiteId = "wxc614e8e34e3bd443"
	}
	if mode == "prod" {
		suiteId = "ww1e1ae4d3739061dd"
	}
	if mode == "release" {
		suiteId = "ww1e1ae4d3739061dd"
	}
	if mode == "uat" {
		suiteId = "ww1e1ae4d3739061dd"
	}

	return suiteId
}

func (a ClientMode) GetThisClient() string {
	client := "console"
	return client
}

// 设置Header签名
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

// HeaderSign 设置请求sdk接口的Header签名
var HeaderSign = HeaderAuthSign()

func NewHttpClient(options ...clients.HTTPClientOption) *SdkClient {
	var client SdkClient
	client.Client = clients.NewHTTPClient(options...)
	return &client
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
