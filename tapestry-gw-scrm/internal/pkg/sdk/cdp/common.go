package cdp

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
	"gw-scrm/internal/config"
	"gw-scrm/internal/pkg/common"
	scrm_errs "gw-scrm/internal/pkg/errs"
	"strconv"
	"strings"
	"time"
)

var urlEnum = struct {
	FansTargetQuery   string
	ShopperListTarget string
	ShopFeeTarget     string
	ChanListQuery     string
}{
	// 公众号指标查询
	FansTargetQuery:   "%s/apis/%s/%s/gzh",
	ShopperListTarget: "%s/apis/%s/query/dg/%s",
	ShopFeeTarget:     "%s/apis/%s/query/index/%s/%s",
	ChanListQuery:     "%s/apis/%s/chanList/%s",
}

var ClientObject = NewSdkClient(clients.Timeout(30 * time.Second))

type TokenInfo struct {
	ExpireTime   int64 `json:"expire_time"`
	IsTmp        int   `json:"is_tmp"`
	EnterpriseId int64 `json:"enterprise_id"`
	AdminId      int64 `json:"admin_id"`
	ServeType    int   `json:"serve_type"`
}

// 调用内部服务接口Client
type DefaultClient struct {
	Client *clients.HTTPClient
}

// 获取cdp请求链接
func GetUrl(uri string) string {
	return fmt.Sprintf(uri, config.CfgData.Service.Cdp.Domain)
}

func GetSign(data map[string]interface{}) (string, error) {
	res, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	secretTools := md5.New()
	md5Str := string(res) + config.CfgData.Service.Cdp.Secret
	secretTools.Write([]byte(md5Str))
	return strings.ToUpper(hex.EncodeToString(secretTools.Sum(nil))), nil
}

func NewSdkClient(options ...clients.HTTPClientOption) *DefaultClient {
	var client DefaultClient
	client.Client = clients.NewHTTPClient(options...)
	return &client
}

// sdk封装发送POST请求接口
func Request(c ctx.Context, method string, url string, body interface{}) (interface{}, error) {
	log.Info(c).Msgf("request sdk, url:%s, request:%+v", url, body)

	var tokenInfo *common.TokenInfo

	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}

	var options []clients.Option
	options = append(options, clients.JSONBody(body))
	options = append(options, clients.Header(common.HeaderKeyEnum.XRequestID, c.ID()))

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
		// todo 避免使用panic
		return nil, errorx.NewErrorMsg(int(scrm_errs.DefaultError.Code), fmt.Sprintf("接口请求错误：%s", res.Response().Status))
	}
	resStr, err := res.Read(clients.Text)
	helper.Must(err)
	log.Info(c).Msgf("sdk response, url:%s, response:%s", url, resStr.(string))
	//
	//var response interface{}
	//helper.Must(json.Unmarshal([]byte(resStr.(string)), &response))
	return resStr, nil
}

// sdk封装发送POST请求接口
func SendPost(c ctx.Context, url string, body interface{}) (interface{}, error) {
	return Request(c, "POST", url, body)
}

func SendGet(c ctx.Context, url string, body interface{}) (interface{}, error) {
	return Request(c, "GET", url, body)
}

// Post 发送 Post 请求
func (s *DefaultClient) Post(c ctx.Context, url string, options ...clients.Option) (*clients.ResponseProxy, error) {
	return s.Client.Post(c, url, options...)
}

// Get 发送 Get 请求
func (s *DefaultClient) Get(c ctx.Context, url string, options ...clients.Option) (*clients.ResponseProxy, error) {
	return s.Client.Get(c, url, options...)
}
