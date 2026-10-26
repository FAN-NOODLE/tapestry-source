package sendcloud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/clients"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"micro-scrm/internal/config"
	"micro-scrm/internal/pkg/errs"
	"micro-scrm/internal/pkg/http"
	"net/url"
	"reflect"
	"strconv"
)

var urlEnum = struct {
	CreateSmsSign        string
	CreateSmsTemplate    string
	UpdateTemplate       string
	SubmitsmsSmsTemplate string
	TemplateList         string
	TemplateSignList     string
	Send                 string
}{
	// 第三方短信平台接口地址
	CreateSmsSign:        "%s/smsapi/sign/save", // 添加短信签名
	CreateSmsTemplate:    "%s/smsapi/addsms",    // 添加短信模板
	UpdateTemplate:       "%s/smsapi/updatesms", // 编辑短信模板
	SubmitsmsSmsTemplate: "%s/smsapi/submitsms", // 提审短信模板
	TemplateList:         "%s/smsapi/list",      // 短信模板列表
	TemplateSignList:     "%s/smsapi/sign/list", // 短信签名列表
	Send:                 "%s/smsapi/send",      // 发送短信
}

var SendCloudSmsKey = "6AyhdQbGfaVq0EbOOuB2vzT0hTU6rwcU" // 短信签名Key
var SendCloudSmsUser = "chedianai_sms"                   // 短信签名User

// CommonResponse 基础返回参数
type CommonResponse struct {
	Result     bool   `json:"result"`
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

func (c CommonResponse) checkResponse() error {
	if c.StatusCode == 200 {
		return nil
	}
	return errorx.NewErrorMsg(c.StatusCode, c.Message)
}

type SendCloudResponse interface {
	checkResponse() error
}

// GetURL 获取cdp请求链接
func GetURL(uri string) string {
	return fmt.Sprintf(uri, config.CfgData.Service.SendCloud.Domain)
}

// CheckResponse 检查sendCloud接口返回结果
func CheckResponse(response CommonResponse) error {
	if response.StatusCode == 200 {
		return nil
	}
	return errorx.NewErrorMsg(response.StatusCode, response.Message)
}

// paramsUrl 组装url
func paramsUrl(urlStr string, body interface{}) string {
	t := reflect.TypeOf(body)
	v := reflect.ValueOf(body)
	keys := make([]string, 0, t.NumField())

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Tag.Get("json")] = v.Field(i).Interface()
		keys = append(keys, t.Field(i).Tag.Get("json"))
	}

	var buf bytes.Buffer
	for _, k := range keys {
		if data[k] == "" {
			continue
		}

		switch value := data[k].(type) {
		case string:
			if value == "" {
				continue
			}
		case struct{}:
			bytes, _ := json.Marshal(value)
			if string(bytes) == "" {
				continue
			}
		case interface{}:
			bytes, _ := json.Marshal(value)
			if string(bytes) == "" {
				continue
			}
		}

		if buf.Len() > 0 {
			buf.WriteByte('&')
		}

		buf.WriteString(k)
		buf.WriteByte('=')
		switch vv := data[k].(type) {
		case string:
			var value = url.QueryEscape(vv)
			buf.WriteString(value)
		case int:
			var value = url.QueryEscape(strconv.FormatInt(int64(vv), 10))
			buf.WriteString(value)
		case int8:
		case int16:
		case int32:
		case int64:
			var value = url.QueryEscape(strconv.FormatInt(int64(vv), 10))
			buf.WriteString(value)
		case struct{}:
			bytes, _ := json.Marshal(vv)
			var value = url.QueryEscape(string(bytes))
			buf.WriteString(value)
		case interface{}:
			bytes, _ := json.Marshal(vv)
			var value = url.QueryEscape(string(bytes))
			buf.WriteString(value)
		default:
			continue
		}
	}

	var paramsStr = buf.String()

	return urlStr + "?" + paramsStr
}

// cdp封装发送POST请求接口
func sendGet(c ctx.Context, urlStr string, body interface{}, response SendCloudResponse) error {

	log.Info(c).Msgf("request sendCloud interface, url:%s, request:%+v", urlStr, body)
	var options []clients.Option
	if body != nil {
		options = append(options, clients.JSONBody(body))
	}

	// 获取请求地址
	urlStr = paramsUrl(urlStr, body)

	res, err := http.CdpClientObject.Get(c, urlStr, options...)
	if err != nil {
		log.Info(c).Msgf("sendCloud接口请求失败：%s, url:%s", err.Error(), urlStr)
		return errorx.NewErrorMsg(errs.ErrorStatus, fmt.Sprintf("sendCloud接口请求失败：%s", err.Error()))
	}
	resStr, err := res.Read(clients.Text)
	if err != nil {
		log.Info(c).Msgf("sendCloud接口返回数据获取失败：%s, url:%s", err.Error(), urlStr)
		return errorx.NewErrorMsg(errs.ErrorStatus, fmt.Sprintf("sendCloud接口返回数据获取失败：%s, url:%s", err.Error(), urlStr))
	}
	log.Info(c).Msgf("request sendCloud interface, url:%s, response:%s", urlStr, resStr.(string))
	err = json.Unmarshal([]byte(resStr.(string)), &response)
	if err != nil {
		log.Info(c).Msgf("sendCloud接口返回数据解析失败：%s, url:%s", err.Error(), urlStr)
		return errorx.NewErrorMsg(errs.ErrorStatus, fmt.Sprintf("sendCloud接口返回数据解析失败：%s, url:%s", err.Error(), urlStr))
	}
	err = response.checkResponse()
	if err != nil {
		log.Info(c).Msgf("sendCloud接口数据错误：%+v", err)
		return err
	}
	return nil
}
