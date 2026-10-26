package sdk

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/clients"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/helper"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"gw-scrm/internal/config"
	"gw-scrm/internal/pkg/common"
	scrmErrs "gw-scrm/internal/pkg/errs"
)

var ClientObject = NewSdkClient(clients.Timeout(60 * time.Second))

// CommonResponse 基础返回参数
type CommonResponse struct {
	App       string      `json:"app"`
	Code      int         `json:"code"`
	Error     int         `json:"error"`
	Msg       string      `json:"message"`
	RequestId string      `json:"request_id"`
	Data      interface{} `json:"data"`
}

// CommonResponseByWork 基础返回参数
type CommonResponseByWork struct {
	App       string      `json:"app"`
	Code      string      `json:"code"`
	Error     string      `json:"error"`
	Msg       string      `json:"message"`
	RequestId string      `json:"request_id"`
	Data      interface{} `json:"data"`
}

type TokenInfo struct {
	ExpireTime   int64 `json:"expire_time"`
	IsTmp        int   `json:"is_tmp"`
	EnterpriseId int64 `json:"enterprise_id"`
	AdminId      int64 `json:"admin_id"`
	ServeType    int   `json:"serve_type"`
}

type DefaultClient struct {
	Client *clients.HTTPClient
}

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

func CheckResponseOld(response CommonResponse) error {
	if response.Code == 200 {
		return nil
	}
	return errorx.NewErrorMsg(response.Code, response.Msg)
}

func Request(c ctx.Context, method string, url string, body interface{}, outLog bool) (interface{}, error) {
	if outLog {
		log.Info(c).Msgf("request sdk, url:%s, request:%+v", url, body)
	}

	var tokenInfo *common.TokenInfo

	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	var hashStr string
	hash := c.Value(common.ContextKeyEnum.EnterpriseHash)
	switch hash.(type) {
	case string:
		hashStr = hash.(string)
	}
	var options []clients.Option
	options = append(options, clients.JSONBody(body))
	options = append(options, clients.Header(common.HeaderKeyEnum.XRequestID, c.ID()))
	options = append(options, clients.Header(common.HeaderKeyEnum.EnterpriseHash, hashStr))

	if tokenInfo != nil {
		options = append(options, clients.Header(common.HeaderKeyEnum.EnterpriseID, strconv.FormatInt(tokenInfo.EnterpriseId, 10)))
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
		return nil, errorx.NewErrorMsg(500, fmt.Sprintf("接口请求错误：%s", res.Response().Status))
	}
	resStr, err := res.Read(clients.Text)
	helper.Must(err)

	if outLog {
		log.Info(c).Msgf("sdk response, url:%s, response:%s", url, resStr.(string))
	}

	var response CommonResponse
	response.Code = res.Response().StatusCode
	response.Msg = res.Response().Status

	helper.Must(json.Unmarshal([]byte(resStr.(string)), &response))

	if outLog {
		log.Info(c).Msgf("json unm is %v\n", response)
	}

	err = CheckResponseOld(response)
	if err != nil {
		log.Info(c).Msgf("sdk response error, err:%+v", err)
		return nil, err
	}
	return response.Data, nil
}

func SendFormData(c ctx.Context, url string, params map[string]string, upFileName string, files ...*multipart.FileHeader) (interface{}, error) {
	log.Info(c).Msgf("request sdk, url:%s, request:%+v, files: %+v", url, params, files)

	var tokenInfo *common.TokenInfo

	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	var hashStr string
	hash := c.Value(common.ContextKeyEnum.EnterpriseHash)
	switch hash.(type) {
	case string:
		hashStr = hash.(string)
	}
	var options []clients.Option
	options = append(options, func(request *http.Request) {
		var rb = &bytes.Buffer{}
		w := multipart.NewWriter(rb)
		for k, v := range params {
			if err := w.WriteField(k, v); err != nil {
				panic(err)
			}
		}
		for _, f := range files {
			fr, err := w.CreateFormFile(upFileName, f.Filename)
			if err != nil {
				panic(err)
			}
			fi, err := f.Open()
			if err != nil {
				panic(err)
			}
			if _, err := io.Copy(fr, fi); err != nil {
				panic(err)
			}
			_ = fi.Close()
		}
		_ = w.Close()
		clients.Body(rb)(request)
		request.Header.Set("Content-Type", w.FormDataContentType())
	})
	options = append(options, clients.Header(common.HeaderKeyEnum.XRequestID, c.ID()))
	options = append(options, clients.Header(common.HeaderKeyEnum.EnterpriseHash, hashStr))

	if tokenInfo != nil {
		options = append(options, clients.Header(common.HeaderKeyEnum.EnterpriseID, strconv.FormatInt(tokenInfo.EnterpriseId, 10)))
	}

	res, err := ClientObject.Client.Request(c, "POST", url, options...)
	if err != nil {
		return nil, err
	}

	if res.Response().StatusCode != 200 {
		return nil, errorx.NewErrorMsg(500, fmt.Sprintf("接口请求错误：%s", res.Response().Status))
	}
	resStr, err := res.Read(clients.Text)
	helper.Must(err)
	log.Info(c).Msgf("sdk response, url:%s, response:%s", url, resStr.(string))

	var response CommonResponse
	response.Code = res.Response().StatusCode
	response.Msg = res.Response().Status

	helper.Must(json.Unmarshal([]byte(resStr.(string)), &response))
	log.Info(c).Msgf("json unm is %v\n", response)

	err = CheckResponseOld(response)
	if err != nil {
		log.Info(c).Msgf("sdk response error, err:%+v", err)
		return nil, err
	}
	return response.Data, nil
}

func RequestByWork(c ctx.Context, method string, url string, body interface{}) (interface{}, error) {
	log.Info(c).Msgf("request sdk, url:%s, request:%+v", url, body)

	var tokenInfo *common.TokenInfo

	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	var hashStr string
	hash := c.Value(common.ContextKeyEnum.EnterpriseHash)
	switch hash.(type) {
	case string:
		hashStr = hash.(string)
	}
	var options []clients.Option
	options = append(options, clients.JSONBody(body))
	options = append(options, clients.Header(common.HeaderKeyEnum.XRequestID, c.ID()))
	options = append(options, clients.Header(common.HeaderKeyEnum.EnterpriseHash, hashStr))

	if tokenInfo != nil {
		options = append(options, clients.Header(common.HeaderKeyEnum.EnterpriseID, strconv.FormatInt(tokenInfo.EnterpriseId, 10)))
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
		return nil, errorx.NewErrorMsg(500, fmt.Sprintf("接口请求错误：%s", res.Response().Status))
	}
	resStr, err := res.Read(clients.Text)
	helper.Must(err)
	log.Info(c).Msgf("sdk response, url:%s, response:%s", url, resStr.(string))

	var response CommonResponseByWork
	response.Code = strconv.Itoa(res.Response().StatusCode)
	response.Msg = res.Response().Status

	helper.Must(json.Unmarshal([]byte(resStr.(string)), &response))
	log.Info(c).Msgf("json unm is %v\n", response)

	// err = CheckResponseOld(response)
	// if err != nil {
	//	log.Info(c).Msgf("sdk response error, err:%+v", err)
	//	return nil, err
	// }
	return response.Data, nil
}

func SendPost(c ctx.Context, url string, body interface{}) (interface{}, error) {
	return Request(c, "POST", url, body, true)
}

func SendPostNoLog(c ctx.Context, url string, body interface{}) (interface{}, error) {
	return Request(c, "POST", url, body, false)
}

func SendPostByWork(c ctx.Context, url string, body interface{}) (interface{}, error) {
	return RequestByWork(c, "POST", url, body)
}

func SendGet(c ctx.Context, url string, body interface{}) (interface{}, error) {
	return Request(c, "GET", url, body, true)
}

type ParseResponse interface {
	Parse([]byte, interface{}) error
}

type ParseJson struct {
}

func (p ParseJson) Parse(bytes []byte, response interface{}) error {
	return json.Unmarshal(bytes, response)
}

// Client 请求代理结构体
type Client struct {
	parse ParseResponse
}

// NewHttpClient 创建Client对象
func NewHttpClient(parse ParseResponse) *Client {
	var client = Client{
		parse: parse,
	}
	return &client
}

type Response interface {
	CheckResponse
}

type CheckResponse interface {
	CheckCode() error
}

type ResponseBase struct {
	App       string `json:"app"`
	Code      int    `json:"code"`
	Msg       string `json:"msg"`
	RequestId string `json:"request_id"`
}

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

	var tokenInfo *common.TokenInfo

	if ctx.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = ctx.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}

	var hashStr string
	hash := ctx.Value(common.ContextKeyEnum.EnterpriseHash)
	switch hash.(type) {
	case string:
		hashStr = hash.(string)
	}
	options = append(options, clients.Header(common.HeaderKeyEnum.XRequestID, ctx.ID()))
	options = append(options, clients.Header(common.HeaderKeyEnum.EnterpriseHash, hashStr))

	if tokenInfo != nil {
		options = append(options, clients.Header(common.HeaderKeyEnum.EnterpriseID, strconv.FormatInt(tokenInfo.EnterpriseId, 10)))
	}

	res, err := ClientObject.Post(ctx, url, options...)
	if err != nil {
		log.Info(ctx).Msgf("接口请求失败：%s, url:%s", err.Error(), url)
		return errorx.NewErrorMsg(scrmErrs.ErrorStatus, fmt.Sprintf("接口请求失败：%s", err.Error()))
	}
	if res.Response().StatusCode != 200 {
		log.Info(ctx).Msgf("接口请求错误：%s, url:%s", res.Response().Status, url)
		return errorx.NewErrorMsg(scrmErrs.ErrorStatus, fmt.Sprintf("接口请求错误：%s", res.Response().Status))
	}

	resStr, err := res.Read(clients.Text)
	if err != nil {
		log.Info(ctx).Msgf("接口返回数据获取失败：%s, url:%s", err.Error(), url)
		return errorx.NewErrorMsg(scrmErrs.ErrorStatus, fmt.Sprintf("接口返回数据获取失败：%s, url:%s", err.Error(), url))
	}
	log.Info(ctx).Msgf("request cdp interface, url:%s, response:%s", url, resStr.(string))

	err = c.parse.Parse([]byte(resStr.(string)), response)
	if err != nil {
		log.Info(ctx).Msgf("接口返回数据解析失败：%s, url:%s", err.Error(), url)
		return errorx.NewErrorMsg(scrmErrs.ErrorStatus, fmt.Sprintf("接口返回数据解析失败：%s, url:%s", err.Error(), url))
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

func (c *Client) Get(ctx ctx.Context, url string, body interface{}, response Response, options ...clients.Option) error {
	log.Info(ctx).Msgf("request interface, url:%s, request:%+v", url, body)
	if body != nil {
		options = append(options, clients.JSONBody(body))
	}

	var tokenInfo *common.TokenInfo

	if ctx.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = ctx.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}

	var hashStr string
	hash := ctx.Value(common.ContextKeyEnum.EnterpriseHash)
	switch hash.(type) {
	case string:
		hashStr = hash.(string)
	}
	options = append(options, clients.Header(common.HeaderKeyEnum.XRequestID, ctx.ID()))
	options = append(options, clients.Header(common.HeaderKeyEnum.EnterpriseHash, hashStr))

	if tokenInfo != nil {
		options = append(options, clients.Header(common.HeaderKeyEnum.EnterpriseID, strconv.FormatInt(tokenInfo.EnterpriseId, 10)))
	}

	res, err := ClientObject.Get(ctx, url, options...)
	if err != nil {
		log.Info(ctx).Msgf("接口请求失败：%s, url:%s", err.Error(), url)
		return errorx.NewErrorMsg(scrmErrs.ErrorStatus, fmt.Sprintf("接口请求失败：%s", err.Error()))
	}
	if res.Response().StatusCode != 200 {
		log.Info(ctx).Msgf("接口请求错误：%s, url:%s", res.Response().Status, url)
		return errorx.NewErrorMsg(scrmErrs.ErrorStatus, fmt.Sprintf("接口请求错误：%s", res.Response().Status))
	}

	resStr, err := res.Read(clients.Text)
	if err != nil {
		log.Info(ctx).Msgf("接口返回数据获取失败：%s, url:%s", err.Error(), url)
		return errorx.NewErrorMsg(scrmErrs.ErrorStatus, fmt.Sprintf("接口返回数据获取失败：%s, url:%s", err.Error(), url))
	}
	log.Info(ctx).Msgf("request cdp interface, url:%s, response:%s", url, resStr.(string))

	err = c.parse.Parse([]byte(resStr.(string)), response)
	if err != nil {
		log.Info(ctx).Msgf("接口返回数据解析失败：%s, url:%s", err.Error(), url)
		return errorx.NewErrorMsg(scrmErrs.ErrorStatus, fmt.Sprintf("接口返回数据解析失败：%s, url:%s", err.Error(), url))
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
