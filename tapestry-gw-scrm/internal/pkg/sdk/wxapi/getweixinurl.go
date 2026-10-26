package wxapi

import (
	"bytes"
	"encoding/json"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/modules/marketing"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// 获取微信跳转链接
func GetWxUrlQuery(c ctx.Context, pageUrl []string, accessToken interface{}) (int64, string, error) {
	WxUrlRes := marketing.WxUrlRes{}
	expireTime := time.Now().Unix() + 15552000 //半年的有效期

	wxUrlReq := map[string]interface{}{
		"jump_wxa": map[string]interface{}{
			"path":  pageUrl[0],
			"query": pageUrl[1],
		},
		"is_expire":   true,
		"expire_time": expireTime,
	}
	bytesData, _ := json.Marshal(wxUrlReq)
	content := strings.Replace(string(bytesData), "\\u0026", "&", -1)

	resp, err := http.Post(GetUrl(urlEnum.GetWxUrl, accessToken), "application/json", bytes.NewReader([]byte(content)))
	if err != nil {
		return 0, "", err
	}
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &WxUrlRes)
	if err != nil || WxUrlRes.ErrCode == 40001 {
		return 40001, "", err
	}

	return WxUrlRes.ErrCode, WxUrlRes.OpenLink, err

}
