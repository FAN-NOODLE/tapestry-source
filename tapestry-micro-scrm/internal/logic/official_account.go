package logic

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	http2 "micro-scrm/internal/pkg/http"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/clients"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/wx/wxoffiaccount/cache"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/wx/wxoffiaccount/common"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/wx/wxoffiaccount/jssdk"
	"micro-scrm/internal/config"
	"micro-scrm/internal/logic/entity"
	"micro-scrm/internal/modules/marketingPromotion/storage"
	"micro-scrm/internal/modules/marketingPromotion/storage/blob"
	"micro-scrm/internal/orm"
	"micro-scrm/internal/orm/model/micro_scrm/enterprise"
	"micro-scrm/internal/pkg/errs"
	"micro-scrm/internal/pkg/redis"
	"micro-scrm/internal/pkg/sdk"
)

// 公众号功能公共逻辑逻辑

const (
	URL_GetMedia = "https://api.weixin.qq.com/cgi-bin/media/get?access_token=%s&media_id=%s"
)

// 重构登陆授权这块，客户可能是两种配置
// 1、微信公众号
// 2、第三方平台
type configValue struct {
	AppID        string                     `json:"app_id"`
	ClientID     string                     `json:"client_id"`
	ClientSecret string                     `json:"client_secret"`
	Secret       string                     `json:"secret"`
	Encrypt      map[string]string          `json:"encrypt"`
	URLS         map[string]string          `json:"urls"`
	ShuValue     configValueByComponentNoJs `json:"shu_value"`
}

type configValueByComponentNoJs struct {
	Id   string `json:"id"`
	Salt string `json:"salt"`
	Url  string `json:"url"`
}

type zxCache struct {
	cache *cache.RedisCache
	key   string
}

type officialAccount struct {
	ctx          ctx.Context
	enterpriseID int64
	model        enterprise.EnterpriseConfigModel
	config       configValue
	http         *sdk.DefaultClient
	zxCache
}

// 参与该功能的enterprise_config.key值集合
var (
	enterpriseConfigKeyMap = map[string]bool{
		"official_info":   true,
		"component":       true,
		"component_no_js": true,
	}
)

const (
	OfficialInfo  = "official_info"
	Component     = "component"
	ComponentNoJs = "component_no_js"
	CYFlag        = "cy_flag"
)

// ContentType
const (
	JSON     = "application/json"
	FORMDATA = "application/x-www-form-urlencoded"
)

func GetRandomString(l int) string {
	str := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	var result []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func NewOfficialAccount(ctx ctx.Context, enterpriseID int64) (*officialAccount, error) {
	// 根据企业ID获取认证的配置
	configList, queryErr := orm.Service.Enterprise.GetConfigListOneEnterprise(enterpriseID)
	if queryErr != nil {
		return nil, errs.SettingsInfoErr
	}

	for _, model := range configList {
		if "" == model.Value {
			continue
		}

		if _, ok := enterpriseConfigKeyMap[model.Key]; ok {
			conf := configValue{}
			jsonErr := json.Unmarshal([]byte(model.Value), &conf)
			if jsonErr != nil {
				return nil, jsonErr
			}

			if 1 == model.IsUse {
				return &officialAccount{
					ctx:          ctx,
					enterpriseID: enterpriseID,
					model:        model,
					config:       conf,
					http:         sdk.NewSdkClient(clients.Timeout(15 * time.Second)),
					zxCache: zxCache{
						cache: cache.NewRedisCache(redis.Client()),
						key:   fmt.Sprintf("%s::%s::%d", config.CfgData.App, "OfficialAccount", enterpriseID),
					},
				}, nil
			}
		}
	}

	return nil, errs.SettingsInfoErr
}

// GetAccessToken 获取access_token
func (o officialAccount) GetAccessToken() (accessToken string, err error) {
	switch o.model.Key {
	case OfficialInfo:
		tokenServer := common.NewAccessTokenCacheService(cache.NewRedisCache(redis.Client()), o.config.AppID, o.config.Secret, o.http.Client)
		return tokenServer.GetToken(o.ctx)
	case Component:
		return o.getThirdToken()
	}

	return "", errors.New("品牌没有公众号没有AccessToken")
}

// 通过三方获取微信公众号的access_token
func (o officialAccount) getWxAccessToken(ctx ctx.Context, token string) (wwToken string, err error) {
	key := fmt.Sprintf("%s::WxToken", o.key)
	if err = o.cache.Get(o.ctx, key, &wwToken); err == nil {
		return
	}

	url := fmt.Sprintf("https://api-wechat.lorealchina.com/api/weixin/get_mp_access_token?access_token=%s", token)

	client := http2.NewHTTPClient()
	res, err := client.Post(ctx, url)
	if err != nil {
		return "", err
	}
	if res.Response().StatusCode != 200 {
		return "", errors.New("接口请求状态码异常")
	}
	rt, rErr := res.Read(clients.Text)
	if rErr != nil {
		return "", rErr
	}

	respStr := rt.(string)
	var result struct {
		ErrorCode int `json:"errcode"`
		Data      struct {
			AccessToken string `json:"access_token"`
			TTL         int    `json:"ttl"`
		} `json:"data"`
		AppID string `json:"appid"`
	}
	err = json.Unmarshal([]byte(respStr), &result)
	if err != nil {
		return
	}

	// 写入缓存
	if result.ErrorCode == 0 {
		cacheErr := o.cache.Set(o.ctx, key, result.Data.AccessToken, 10)
		if cacheErr != nil {
			log.Info(o.ctx).Msgf("缓存%s设置错误accessToken:%+v", key, cacheErr)
		}
		wwToken = result.Data.AccessToken
		return
	} else {
		err = errors.New("请求微信公众号access_token失败")
		return
	}

	return
}

// GetMedia 获取素材
func (o officialAccount) GetMedia(ctx ctx.Context, mediaID string) (mediaUrl string, err error) {
	accessToken, accessTokenErr := o.GetAccessToken()
	if accessTokenErr != nil {
		err = accessTokenErr
		return
	}

	url := ""
	switch o.model.Key {
	case OfficialInfo:
		url = fmt.Sprintf(URL_GetMedia, accessToken, mediaID)
	case Component:
		token, wxTokenErr := o.getWxAccessToken(ctx, accessToken)
		if wxTokenErr != nil {
			err = wxTokenErr
			return
		}
		url = fmt.Sprintf(URL_GetMedia, token, mediaID)
	}

	response, responseErr := o.http.Get(ctx, url)
	if responseErr != nil {
		err = responseErr
		log.Info(ctx).Msgf("公众号::企业ID：%d::请求素材内容失败:%v", o.enterpriseID, responseErr)
		return
	}

	contentType := response.Response().Header.Get("Content-Type")
	if strings.Contains(contentType, "application/json") {
		err = errors.New("微信上传图片失败")
		return
	}

	mediaType := "png"
	switch contentType {
	case "image/jpeg", "image/jpg", "image/pjpeg":
		mediaType = "jpg"
	case "image/png", "image/x-png":
		mediaType = "png"
	}

	log.Info(ctx).Msgf("公众号::企业ID：%d::获取素材内容接口:%v", o.enterpriseID, response)

	mediaBytes, err := ioutil.ReadAll(response.Response().Body)
	if err != nil {
		log.Info(ctx).Msgf("公众号::企业ID：%d::获取素材内容失败:%v", o.enterpriseID, err)
		return
	}

	filename := fmt.Sprintf("%d.%s", time.Now().UnixNano(), mediaType)
	newFile, _ := os.Create(filename)
	newFile.Write(mediaBytes)
	defer func() {
		newFile.Close()
		os.Remove(filename)
	}()
	cs := storage.Init(blob.BLOB)
	path, err := cs.Upload(ctx, filename, filename)
	if err != nil {
		log.Info(ctx).Msgf("uploadCloud Error: %+v", err)
		return
	}
	mediaUrl = cs.GetURL(path)

	return
}

// 获取js-config
func (o officialAccount) GetJsConfig(urlStr string) (config entity.JsConfig, err error) {
	if o.model.Key == ComponentNoJs {
		shuConfig, err := o.getShuJsConfig(urlStr)
		if err != nil {
			return config, err
		}
		t, _ := strconv.ParseInt(shuConfig.Data.Timestamp, 10, 64)
		config = entity.JsConfig{
			AppID:     shuConfig.Data.AppID,
			NonceStr:  shuConfig.Data.NonceStr,
			Timestamp: t,
			Signature: shuConfig.Data.Signature,
		}
		return config, err
	}
	ticket := ""
	switch o.model.Key {
	case OfficialInfo:
		// copy之前的代码
		tokenServer := common.NewAccessTokenCacheService(cache.NewRedisCache(redis.Client()), o.config.AppID, o.config.Secret, clients.NewHTTPClient())
		accessToken, err := tokenServer.GetToken(o.ctx)
		if err != nil {
			return config, err
		}
		log.Info(o.ctx).Msgf("accessToken:%s", accessToken)

		ticketService := jssdk.NewTicketCacheService(cache.NewRedisCache(redis.Client()), o.config.AppID, accessToken, clients.NewHTTPClient())
		ticket, err = ticketService.GetTicket(o.ctx)
		log.Info(o.ctx).Msgf("ticket:%s", ticket)

		if err != nil {
			log.Info(o.ctx).Msgf("获取jsTicket失败, err:%s", err.Error())
			return config, err
		}
	case Component:
		token, err := o.getThirdToken()
		if err != nil {
			return config, err
		}
		log.Info(o.ctx).Msgf("accessToken:%s", token)

		ticket, err = o.getThirdJsConfig(token)
		if err != nil {
			return config, err
		}
		log.Info(o.ctx).Msgf("accessToken:%s", ticket)
	}

	nonceStr := GetRandomString(32)
	timestamp := time.Now().Unix()
	config = entity.JsConfig{
		AppID:     o.config.AppID,
		NonceStr:  nonceStr,
		Timestamp: timestamp,
		Signature: jssdk.WXConfigSign(ticket, nonceStr, strconv.FormatInt(timestamp, 10), urlStr),
	}

	return config, err
}

func (o officialAccount) getShuJsConfig(urlStr string) (config entity.ShuJsConfig, err error) {
	timeUnix := strconv.FormatInt(time.Now().Unix(), 10)
	signTex := o.config.ShuValue.Salt + timeUnix
	sign := substr(md5V(signTex), 0, 8)
	postData := make(map[string]string)
	postData["id"] = o.config.ShuValue.Id
	postData["slat"] = o.config.ShuValue.Salt
	postData["timestamp"] = timeUnix
	postData["sign"] = sign
	postData["url"] = urlStr
	postData["siteurl"] = urlStr
	body, err := postWithFormData("POST", o.config.URLS["ticket"], &postData)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &config)
	return
}

// 获取第三方的token
func (o officialAccount) getThirdToken() (token string, err error) {
	key := fmt.Sprintf("%s::ThirdToken", o.key)
	if err = o.cache.Get(o.ctx, key, &token); err == nil {
		return "", err
	}
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	reqParams := entity.ThirdTokenRequest{
		GrantType:    "client_credentials",
		ClientID:     o.config.ClientID,
		ClientSecret: o.config.ClientSecret,
	}

	jsonStr, jsonErr := json.Marshal(reqParams)
	if jsonErr != nil {
		return "", jsonErr
	}

	resp, err := client.Post(o.config.URLS["token"], "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respByte, ioErr := ioutil.ReadAll(resp.Body)
	if ioErr != nil {
		return "", ioErr
	}

	var result entity.ThirdTokenResponse
	err = json.Unmarshal(respByte, &result)
	if err != nil {
		return
	}

	if result.ErrorCode == 0 {
		cacheErr := o.cache.Set(o.ctx, key, result.AccessToken, result.ExpireIn)
		if cacheErr != nil {
			log.Info(o.ctx).Msgf("缓存%s设置错误accessToken:%+v", key, cacheErr)
		}
		token = result.AccessToken
		return
	} else {
		err = errors.New(result.ErrorMsg)
		return
	}
}

// 获取第三方的js-config
func (o officialAccount) getThirdJsConfig(token string) (ticket string, err error) {
	key := fmt.Sprintf("%s::ThirdTicket", o.key)
	if err = o.cache.Get(o.ctx, key, &ticket); err == nil {
		return "", err
	}
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	url := fmt.Sprintf("%s?access_token=%s", o.config.URLS["ticket"], token)

	resp, err := client.Post(url, "application/json", bytes.NewBuffer([]byte{}))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result entity.ThirdJsTicketResponse
	jsonErr := json.Unmarshal(respBytes, &result)
	if jsonErr != nil {
		return "", jsonErr
	}

	if result.ErrorCode == 0 {
		cacheErr := o.cache.Set(o.ctx, key, result.Data.Ticket, result.Data.ExpiresIn)
		if cacheErr != nil {
			log.Info(o.ctx).Msgf("缓存%s设置错误accessToken:%+v", key, cacheErr)
		}

		return result.Data.Ticket, nil
	} else {
		return "", errors.New("token获取失败")
	}
}

func md5V(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func substr(str string, start, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0
	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length
	if start > end {
		start, end = end, start
	}
	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}
	return string(rs[start:end])
}

func postWithFormData(method, url string, postData *map[string]string) ([]byte, error) {
	body := new(bytes.Buffer)
	w := multipart.NewWriter(body)
	for k, v := range *postData {
		w.WriteField(k, v)
	}
	w.Close()
	req, _ := http.NewRequest(method, url, body)
	req.Header.Set("Content-Type", w.FormDataContentType())
	resp, _ := http.DefaultClient.Do(req)
	data, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	return data, err
}
