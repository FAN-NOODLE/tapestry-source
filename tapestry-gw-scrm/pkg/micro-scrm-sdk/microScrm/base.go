package microScrm

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/clients"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/json"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
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

	httpClient := clients.NewHTTPClient()
	val := c.Value(EnterpriseHashKey)
	hash, _ := val.(string)

	var option []clients.Option
	option = append(option, Header(ServiceId, SecretKey, hash))
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

func Header(serviceId, key, enterpriseHash string) clients.Option {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	sign := sign(serviceId, key, timestamp)

	return func(request *http.Request) {
		request.Header.Set("Auth-Timestamp", timestamp)
		request.Header.Set("Auth-Sign", sign)
		request.Header.Set("Auth-Client", serviceId)
		request.Header.Set("Enterprise-Hash", enterpriseHash)

	}
}

func sign(serviceId, key string, timestamp string) string {
	m := md5.New()
	m.Write([]byte(serviceId))
	m.Write([]byte(key))
	m.Write([]byte(timestamp))

	return hex.EncodeToString(m.Sum(nil))
}
