package qyWork

import (
	"errors"
	"fmt"
	"micro-scrm/internal/app"
	microBossDao "micro-scrm/internal/app/dao/micro_boss"
	"time"
)

const (
	getToken = "cgi-bin/gettoken?corpid=%s&corpsecret=%s"
)

// GetToken
// identifier 应用标识
func (w *QyWork) getToken(c app.Ctx, identifier string) (string, error) {
	var qyTokenKey = "qy_token:%d:%s"
	qyTokenKey = fmt.Sprintf(qyTokenKey, w.enterpriseId, identifier)
	tokenCache := app.Redis().Get(c, qyTokenKey).Val()
	if tokenCache != "" {
		return tokenCache, nil
	}

	pluginInfo, err := microBossDao.NewEnterpriseWorkPlugin().GetOne(c, map[string]interface{}{
		"enterprise_id": w.enterpriseId,
		"identifier":    identifier,
	})

	if err != nil || pluginInfo.Secret == "" {
		return "", errors.New("获取企业微信自建应用失败")
	}

	var tokenRes struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int64  `json:"expires_in"`
	}

	err = w.requestPostJson(c, fmt.Sprintf(w.getBaseUrl(getToken), pluginInfo.CorpID, pluginInfo.Secret), nil, &tokenRes)
	if err != nil {
		return "", err
	}

	app.Redis().Set(c, qyTokenKey, tokenRes.AccessToken, time.Duration(tokenRes.ExpiresIn-100)*time.Second)

	return tokenRes.AccessToken, nil
}
