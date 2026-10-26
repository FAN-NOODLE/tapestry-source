package cdp

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/clients"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"io/ioutil"
	"micro-scrm/internal/config"
	"micro-scrm/internal/pkg/common"
	"micro-scrm/internal/pkg/errs"
	"micro-scrm/internal/pkg/http"
	nethttp "net/http"
	"strconv"
	"strings"
	"time"
)

// CommonResponse 基础返回参数
type CommonResponse struct {
	App       string      `json:"app"`
	Code      int         `json:"code"`
	Msg       string      `json:"msg"`
	RequestID string      `json:"request_id"`
	Data      interface{} `json:"data"`
}

// Response 返回值接口
type Response interface {
	ResponseCheck
}

// ResponseCheck 返回值校验接口
type ResponseCheck interface {
	CheckCode() error
}

// CheckCodeResponse 检查code的返回值结构体头
type CheckCodeResponse struct {
	App       string `json:"app"`
	Code      int    `json:"code"`
	Msg       string `json:"msg"`
	RequestID string `json:"request_id"`
}

// CheckCode 检查接口返回code
func (r *CheckCodeResponse) CheckCode() error {
	if r.Code == 0 {
		return nil
	}
	return errorx.NewErrorMsg(r.Code, r.Msg)
}

// CheckCodeResponse 检查code的返回值结构体头
type CheckErrorResponse struct {
	Error int    `json:"error"`
	Msg   string `json:"msg"`
}

// CheckCode 检查接口返回code
func (r *CheckErrorResponse) CheckCode() error {
	if r.Error == 0 {
		return nil
	}
	return errorx.NewErrorMsg(r.Error, r.Msg)
}

// cdp接口url配置
var urlEnum = struct {
	UserTagList                    string // 标签分组列表
	ManualGroupCreate              string // 创建手动标签分组
	ManualTagCreate                string // 创建手动标签
	AutoTagUpdate                  string // 更新自动标签
	CreateUserGroup                string // 创建用户分组
	EditUserGroup                  string // 编辑客户分组
	DeleteUserGroup                string // 删除客户分组
	UserGroupList                  string // 用户分组列表
	UpdateTagGroupSync             string // 更新标签组同步企业微信标识
	TagGroupInfo                   string // 标签组详情
	SortTagGroup                   string // 标签组排序
	SearchTag                      string // 标签搜索
	DeleteManualTagGroup           string // 删除标签组
	UserGroupPrecalculateCount     string // 客户分组预计算人数
	UpdateManualTag                string // 编辑更新手动标签
	DeleteManualTag                string // 删除手动标签
	TagEnterpriseInit              string // 标签初始化
	GetMarketingStrategyDetailsUrl string //策略数据
	UserGroupDetail                string //策略数据
	GetCouponStrategyDetailUrl     string //自动发券策略
	GetUserGroupsByUid             string //根据用户信息查询用户分组信息
}{
	UserTagList:                    "%s/scrm/user/tag/list",
	ManualGroupCreate:              "%s/scrm/user/tag/manual/group/create",
	ManualTagCreate:                "%s/scrm/user/tag/manual/tag/create",
	AutoTagUpdate:                  "%s/scrm/user/tag/update",
	CreateUserGroup:                "%s/scrm/user/group/create",
	EditUserGroup:                  "%s/scrm/user/group/edit",
	DeleteUserGroup:                "%s/scrm/user/group/delete",
	UserGroupList:                  "%s/scrm/user/group/list",
	UpdateTagGroupSync:             "%s/scrm/user/update/sync_qw",
	TagGroupInfo:                   "%s/scrm/user/group/get",
	SortTagGroup:                   "%s/scrm/user/tag/rank",
	SearchTag:                      "%s/scrm/user/tag/search",
	DeleteManualTagGroup:           "%s/scrm/user/tag/manual/group/delete",
	UserGroupPrecalculateCount:     "%s/scrm/user/group/precalculate/usercount",
	UpdateManualTag:                "%s/scrm/user/tag/manual/tag/edit",
	DeleteManualTag:                "%s/scrm/user/tag/manual/tag/delete",
	TagEnterpriseInit:              "%s/scrm/enterprise/init",
	GetMarketingStrategyDetailsUrl: "%s/apis/%d/query/market",
	UserGroupDetail:                "%s/scrm/user/group/details",
	GetCouponStrategyDetailUrl:     "%s/apis/%d/getCoupon",
	GetUserGroupsByUid:             "%s/scrm/user/query/groups",
}

// GetURL 获取cdp请求链接
func GetURL(uri string, a ...interface{}) string {
	var paramList = []interface{}{
		config.CfgData.Service.Cdp.Domain,
	}
	for _, v := range a {
		paramList = append(paramList, v)
	}
	return fmt.Sprintf(uri, paramList...)
}

// CheckResponse 检查cdp接口返回结果
func CheckResponse(response CommonResponse) error {
	if response.Code == 0 {
		return nil
	}
	return errorx.NewErrorMsg(response.Code, response.Msg)
}

// cdp封装发送POST请求接口
func sendPost(c ctx.Context, url string, body interface{}) (interface{}, error) {
	log.Info(c).Msgf("request cdp interface, url:%s, request:%+v", url, body)
	var options []clients.Option
	if body != nil {
		options = append(options, clients.JSONBody(body))
	}
	options = append(options, HeaderAuthSign(c))
	res, err := http.CdpClientObject.Post(c, url, options...)
	if err != nil {
		log.Info(c).Msgf("cdp接口请求失败：%s, url:%s", err.Error(), url)
		return nil, errorx.NewErrorMsg(errs.ErrorStatus, fmt.Sprintf("cdp接口请求失败：%s", err.Error()))
	}
	if res.Response().StatusCode != 200 {
		log.Info(c).Msgf("cdp接口请求错误：%s, url:%s", res.Response().Status, url)
		return nil, errorx.NewErrorMsg(errs.ErrorStatus, fmt.Sprintf("cdp接口请求错误：%s", res.Response().Status))
	}

	resStr, err := res.Read(clients.Text)
	if err != nil {
		log.Info(c).Msgf("cdp接口返回数据获取失败：%s, url:%s", err.Error(), url)
		return nil, errorx.NewErrorMsg(errs.ErrorStatus, fmt.Sprintf("cdp接口返回数据获取失败：%s, url:%s", err.Error(), url))
	}
	responseStr := resStr.(string)
	if len(responseStr) > 50000 {
		responseStr = responseStr[:50000]
	}
	log.Info(c).Msgf("request cdp interface, url:%s, response:%s", url, responseStr)
	var response CommonResponse
	err = json.Unmarshal([]byte(resStr.(string)), &response)
	if err != nil {
		log.Info(c).Msgf("cdp接口返回数据解析失败：%s, url:%s", err.Error(), url)
		return nil, errorx.NewErrorMsg(errs.ErrorStatus, fmt.Sprintf("cdp接口返回数据解析失败：%s, url:%s", err.Error(), url))
	}
	err = CheckResponse(response)
	if err != nil {
		log.Info(c).Msgf("cdp接口数据错误：%+v", err)
		return nil, err
	}
	return &response.Data, nil
}

var cdpClient = clients.NewHTTPClient()

// Client 请求代理结构体
type Client struct {
}

// NewHTTPClient 创建Client对象
func NewHTTPClient() *Client {
	var client = Client{}
	return &client
}

// HeaderAuthSign 设置Header签名
func HeaderAuthSign(c ctx.Context) clients.Option {
	return func(request *nethttp.Request) {
		enterpriseID := common.EnterpriseID(c)
		enterpriseHash := common.EnterpriseHash(c)

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
		log.Info(c).Msgf("cdp接口请求签名信息,签名参数:%s, 签名结果: %s, 企业id:%d", bodyString+curTime+enterpriseHash, signStr, enterpriseID)

		request.Header.Set("Enterprise-ID", strconv.FormatInt(enterpriseID, 10))
		request.Header.Set("Auth-Sign", signStr)
		request.Header.Set("Auth-Timestamp", curTime)
	}
}

// Post 发送Post请求
func (c *Client) Post(ctx ctx.Context, url string, body interface{}, response Response) error {
	log.Info(ctx).Msgf("request cdp interface, url:%s, request:%+v", url, body)
	var options []clients.Option
	if body != nil {
		options = append(options, clients.JSONBody(body))
	}
	options = append(options, HeaderAuthSign(ctx))
	res, err := cdpClient.Post(ctx, url, options...)
	// res, err := clients.HTTPClient.Post(ctx, url, options...)
	if err != nil {
		log.Info(ctx).Msgf("cdp接口请求失败：%s, url:%s", err.Error(), url)
		return errorx.NewErrorMsg(errs.ErrorStatus, fmt.Sprintf("cdp接口请求失败：%s", err.Error()))
	}
	if res.Response().StatusCode != 200 {
		log.Info(ctx).Msgf("cdp接口请求错误：%s, url:%s", res.Response().Status, url)
		return errorx.NewErrorMsg(errs.ErrorStatus, fmt.Sprintf("cdp接口请求错误：%s", res.Response().Status))
	}

	resStr, err := res.Read(clients.Text)
	if err != nil {
		log.Info(ctx).Msgf("cdp接口返回数据获取失败：%s, url:%s", err.Error(), url)
		return errorx.NewErrorMsg(errs.ErrorStatus, fmt.Sprintf("cdp接口返回数据获取失败：%s, url:%s", err.Error(), url))
	}
	log.Info(ctx).Msgf("request cdp interface, url:%s, response:%s", url, strings.Replace(resStr.(string), "\n", "", -1))
	log.Info(ctx).Msgf("1026-cdp接口返回数据：%s", resStr.(string))
	err = json.Unmarshal([]byte(resStr.(string)), response)
	if err != nil {
		log.Info(ctx).Msgf("1026-cdp接口返回数据解析失败：%s, url:%s", err.Error(), url)
		return errorx.NewErrorMsg(errs.ErrorStatus, fmt.Sprintf("cdp接口返回数据解析失败：%s, url:%s", err.Error(), url))
	}
	err = response.CheckCode()
	if err != nil {
		log.Info(ctx).Msgf("cdp接口数据错误：%+v", err)
		return err
	}
	return nil
}
