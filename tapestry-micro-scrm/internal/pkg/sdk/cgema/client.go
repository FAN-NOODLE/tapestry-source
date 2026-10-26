package cgema

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/clients"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"micro-scrm/internal/config"
	"micro-scrm/internal/pkg/sdk"
)

const (
	sendResultPath         = "/api/ma/task/sendResult"
	sendFinishPath         = "/api/ma/saCompletedTask"
	sendMomentResultPath   = "/api/coach-wecom/callback/moment_send_result"
	sendMomentCommentsPath = "/api/coach-wecom/callback/moment_comments"
	relationDataPath       = "/expand/api/v1/projects/wecom/data-collect"
)

type Client struct {
	domain    string
	appKey    string
	appSecret string
	timeout   time.Duration
}

func NewClient() *Client {
	conf := config.CfgData.CgeMa
	return &Client{
		domain:    strings.TrimRight(conf.Domain, "/"),
		appKey:    conf.AppKey,
		appSecret: conf.AppSecret,
		timeout:   50 * time.Second,
	}
}

func (c *Client) SendResult(context ctx.Context, request SendResultRequest) (Response, error) {
	return c.post(context, sendResultPath, request)
}

func (c *Client) SendFinish(context ctx.Context, request SendFinishRequest) (Response, error) {
	return c.post(context, sendFinishPath, request)
}

func (c *Client) SendMomentResult(context ctx.Context, request SendMomentResultRequest) (Response, error) {
	return c.post(context, sendMomentResultPath, request)
}

func (c *Client) SendMomentComments(context ctx.Context, request SendMomentCommentsRequest) (Response, error) {
	return c.post(context, sendMomentCommentsPath, request)
}

func (c *Client) PushRelationEvent(context ctx.Context, request RelationEventRequest) (Response, error) {
	return c.postURL(context, c.relationURL(), request)
}

func (c *Client) post(context ctx.Context, path string, request interface{}) (Response, error) {
	url := c.domain + path
	return c.postURL(context, url, request)
}

func (c *Client) postURL(context ctx.Context, url string, request interface{}) (Response, error) {
	var response Response
	if c.domain == "" || c.appKey == "" || c.appSecret == "" {
		err := errors.New("深演MA配置不完整")
		c.alert(context, "深演MA-配置异常", url, nil, request, response, err)
		return response, err
	}

	body, err := json.Marshal(request)
	if err != nil {
		c.alert(context, "深演MA-请求参数序列化失败", url, nil, request, response, err)
		return response, err
	}

	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
	nonce, err := newNonce()
	if err != nil {
		c.alert(context, "深演MA-Nonce生成失败", url, nil, string(body), response, err)
		return response, err
	}
	signature := sign(c.appSecret, c.appKey+timestamp+nonce+string(body))
	idempotency := buildIdempotencyKey(body, timestamp, nonce)
	headers := buildRequestHeaders(c.appKey, timestamp, nonce, signature, idempotency)

	client := clients.NewHTTPClient(clients.Timeout(c.timeout))
	res, err := client.Post(context, url,
		clients.Body(bytes.NewReader(body)),
		clients.ContentTypeJSON,
		clients.Header("X-App-Key", headers["X-App-Key"]),
		clients.Header("X-Timestamp", headers["X-Timestamp"]),
		clients.Header("X-Nonce", headers["X-Nonce"]),
		clients.Header("X-Signature", headers["X-Signature"]),
		clients.Header("X-Idempotency-Key", headers["X-Idempotency-Key"]),
	)
	if err != nil {
		c.alert(context, "深演MA-接口请求失败", url, headers, string(body), response, err)
		return response, err
	}

	responseText, err := res.Read(clients.Text)
	if err != nil {
		c.alert(context, "深演MA-响应读取失败", url, headers, string(body), response, err)
		return response, err
	}
	respBody := responseText.(string)

	if res.Response().StatusCode != 200 {
		err = fmt.Errorf("HTTP状态码异常: %d", res.Response().StatusCode)
		c.alert(context, "深演MA-HTTP状态码异常", url, headers, string(body), respBody, err)
		return response, err
	}

	if err = json.Unmarshal([]byte(respBody), &response); err != nil {
		c.alert(context, "深演MA-响应解析失败", url, headers, string(body), respBody, err)
		return response, err
	}
	if !isSuccessCode(response.Code) {
		err = fmt.Errorf("业务状态码异常: code=%d, message=%s", response.Code, response.Message)
		c.alert(context, "深演MA-业务状态码异常", url, headers, string(body), respBody, err)
		return response, err
	}

	log.Info(context).Msgf("深演MA接口请求成功, url:%s, headers:%s, request:%s, response:%s, code:%d, idempotency_key:%s", url, marshalLog(headers), string(body), respBody, response.Code, idempotency)
	return response, nil
}

func (c *Client) relationURL() string {
	return strings.TrimSuffix(c.domain, "/dcs") + relationDataPath
}

func (c *Client) alert(context ctx.Context, title, url string, headers map[string]string, body, response interface{}, err error) {
	log.Error(context).Err(err).Msgf("%s, url:%s, headers:%s", title, url, marshalLog(headers))
	go sdk.SendMessage(context, title, url, body, response, err)
}

func buildRequestHeaders(appKey, timestamp, nonce, signature, idempotency string) map[string]string {
	return map[string]string{
		"X-App-Key":         appKey,
		"X-Timestamp":       timestamp,
		"X-Nonce":           nonce,
		"X-Signature":       signature,
		"X-Idempotency-Key": idempotency,
	}
}

func marshalLog(value interface{}) string {
	body, err := json.Marshal(value)
	if err != nil {
		return fmt.Sprintf("%v", value)
	}
	return string(body)
}

func newNonce() (string, error) {
	value := make([]byte, 16)
	if _, err := rand.Read(value); err != nil {
		return "", err
	}
	return hex.EncodeToString(value), nil
}

func sign(secret, content string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(content))
	return hex.EncodeToString(mac.Sum(nil))
}

func buildIdempotencyKey(body []byte, timestamp, nonce string) string {
	value := string(body) + "_" + timestamp + "_" + nonce
	sum := md5.Sum([]byte(value))
	return hex.EncodeToString(sum[:])
}

func isSuccessCode(code int) bool {
	return code == 0 || code == 1006
}
