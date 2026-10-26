package microScrm

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/clients"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/json"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"gw-work-scrm/internal/config"
	"gw-work-scrm/internal/pkg/common"
	"net/http"
	"strconv"
	"time"
)

const StatusOK = 200

var urlHost string

const (
	ServiceId         = "micro-scrm"
	EnterpriseHashKey = "enterprise_hash"
	SecretKey         = "REDACTED_FOR_SOURCE_SHARING"
)

type BaseSDK struct {
	c *clients.HTTPClient
}

func (b BaseSDK) Post(c ctx.Context, url string, req interface{}) (*clients.ResponseProxy, error) {

	httpClient := clients.NewHTTPClient(clients.Timeout(59 * time.Second))
	val := c.Value(common.ContextKeyEnum.EnterpriseHash)
	hash, _ := val.(string)

	var option []clients.Option

	// 从授权信息中获取信息
	enterpriseID := ""
	if eid := c.Value("enterpriseId"); eid != nil {
		v, ok := eid.(int64)
		if ok {
			enterpriseID = strconv.FormatInt(v, 10)
		} else {
			log.Error(c).Msg("上下文中获取enterpriseId类型出错")
		}
	}

	option = append(option, Header(ServiceId, SecretKey, hash, enterpriseID))
	option = append(option, clients.JSONBody(req))

	start := time.Now()
	res, err := httpClient.Post(ctx.Wrap(c), url, option...)
	spend := time.Since(start)

	var headers interface{}
	if res != nil {
		headers = res.Request().Header
	}

	if err != nil {
		logInfo(c, url, headers, req, fmt.Sprintf("服务异常: %v", err), spend)
		return nil, err
	}

	if statusCode := res.Response().StatusCode; statusCode != 200 {
		logInfo(c, url, headers, req, fmt.Sprintf("状态异常: %d", statusCode), spend)
		if err == nil {
			errMsg := fmt.Sprintf("请求%s接口错误,HTTPcode=%d", ServiceId, statusCode)
			err = errorx.New(errMsg)
		}
		return nil, err
	}

	logInfo(c, url, headers, req, res, spend)

	return res, nil
}

func logInfo(c ctx.Context, url string, headers, req, resp interface{}, spend time.Duration) {

	type format struct {
		Url      string
		Headers  interface{}
		Params   interface{}
		Time     string
		Response interface{}
	}

	info := &format{
		Url:      url,
		Headers:  headers,
		Params:   req,
		Time:     fmt.Sprintf("%v", spend),
		Response: resp,
	}

	log.Info(c).Msgf("SDK-Micro-Scrm [Post] %s \n\n", json.MarshalToString(info))
}

func Header(serviceId, key, enterpriseHash, enterpriseID string) clients.Option {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	sign := sign(serviceId, key, timestamp)

	return func(request *http.Request) {
		request.Header.Set("Auth-Timestamp", timestamp)
		request.Header.Set("Auth-Sign", sign)
		request.Header.Set("Auth-Client", serviceId)
		request.Header.Set("Enterprise-Hash", enterpriseHash)
		request.Header.Set("Enterprise-ID", enterpriseID)
	}
}

func sign(serviceId, key string, timestamp string) string {
	m := md5.New()
	m.Write([]byte(serviceId))
	m.Write([]byte(key))
	m.Write([]byte(timestamp))

	return hex.EncodeToString(m.Sum(nil))
}

func (b BaseSDK) GetUrl(uri string) string {
	return fmt.Sprintf(uri, config.CfgData.Service.MicroScrm.Domain)
}
