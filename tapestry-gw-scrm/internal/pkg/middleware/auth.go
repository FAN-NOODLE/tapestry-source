package middleware

import (
	"encoding/json"
	"fmt"
	"gw-scrm/internal/modules/right_rule"
	"gw-scrm/internal/pkg/sdk/microboss"
	"net/http"
	"time"

	redis2 "github.com/go-redis/redis/v8"

	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	middleware "gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus/mw"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/helper"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"gw-scrm/internal/config"
	"gw-scrm/internal/pkg/common"
	"gw-scrm/internal/pkg/redis"
)

const baseLimitNum = 100

var (
	fcUri = map[string]int{
		"gw-scrm/open/order": baseLimitNum,
	}
)

// Auth 崩溃恢复中间件
func AuthCheck(skippers ...middleware.SkipperFunc) gin.HandlerFunc {
	return middleware.Register("AuthCheck", func(c *gin.Context) {
		if middleware.SkipHandler(c, skippers...) {
			log.Info(ctx.Wrap(c)).Msgf("AuthCheck, path:%s do not check", c.Request.URL.Path)
			c.Next()
			return
		}
		log.Info(ctx.Wrap(c)).Msgf("AuthCheck, path: %s", c.Request.URL.Path)
		token := c.GetHeader(common.HeaderKeyEnum.AuthToken)
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":       451004,
				"message":    "身份信息异常",
				"request_id": ctx.Wrap(c).ID(),
				"data":       nil,
			})
			c.Abort()
			return
		}

		// 验证token
		key := fmt.Sprintf("%s:token:%s", config.CfgData.App, token)
		reRet := redis.Client().Get(c, key).Val()
		if reRet == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":       451004,
				"message":    "身份过期",
				"request_id": ctx.Wrap(c).ID(),
				"data":       nil,
			})
			c.Abort()
			return
		}

		var tokenInfo common.TokenInfo
		if err := json.Unmarshal([]byte(reRet), &tokenInfo); err != nil {
			ginplus.ResError(c, errorx.ErrInvalidToken)
			c.Abort()
			return
		}
		// 验证过期时间
		if tokenInfo.ExpireTime < time.Now().Unix() {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":       451004,
				"message":    "身份过期",
				"request_id": ctx.Wrap(c).ID(),
				"data":       nil,
			})
			c.Abort()
			return
		}
		// 如果操作 则每次刷新token的过期时间
		m, _ := time.ParseDuration("15m")
		expireTime := time.Now().Add(m).Unix()
		tokenInfo.ExpireTime = expireTime
		info, err := json.Marshal(tokenInfo)
		helper.Must(err)
		redis.Client().SetEX(c, key, string(info), time.Second*60*60*12)

		path := c.FullPath()
		if path == "" {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "路由不存在"})
			c.Abort()
			return
		}
		var AdminDetailRequest = right_rule.AdminDetailRequest{
			AdminId: tokenInfo.AdminId,
		}
		AdminDetailRes, err := microboss.GetAdminDetail(ctx.Wrap(c), AdminDetailRequest)
		if err != nil {
			ginplus.ResError(c, err)
			c.Abort()
			return
		}
		if AdminDetailRes.Data.Status == 2 {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 451000, "message": "您的管理员身份已被禁用"})
			c.Abort()
			return
		}
		tokenInfo.Mobile = AdminDetailRes.Data.Mobile
		tokenInfo.Username = AdminDetailRes.Data.Username
		// 不需要进行签名检查的请求
		// if PathNotCheckPower(c) == false {
		//	access, err := microauthority.AccessControl(ctx.Wrap(c), microauthority.AccessControlReq{
		//		AdminID:      tokenInfo.AdminId,
		//		EnterpriseID: tokenInfo.EnterpriseId,
		//		URL:          string([]byte(path)[1:]),
		//	})
		//	if err != nil {
		//		c.JSON(http.StatusUnauthorized, gin.H{"code": 451005, "message": err})
		//		c.Abort()
		//		return
		//	}
		//	if !access {
		//		c.JSON(http.StatusUnauthorized, gin.H{"code": 451006, "message": "无权限访问"})
		//		c.Abort()
		//		return
		//	}
		// }
		c.Set(common.ContextKeyEnum.TokenInfo, &tokenInfo)

		if !frequencyControl(c, string([]byte(c.FullPath())[1:])) {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 451008, "message": "频次超限"})
			c.Abort()
			return
		}
		log.Info(ctx.Wrap(c)).Msgf("请求地址:%s, 请求参数: %v, AdminID:%d", c.Request.URL.Path, c.Request.Body, tokenInfo.AdminId)
		c.Next()
	})
}

// PathNotCheckAuth 不需要进行签名检查的请求
func PathNotCheckAuth() middleware.SkipperFunc {
	return func(c *gin.Context) bool {
		path := c.Request.URL.Path

		// todo appname&模糊匹配&公共配置
		doNotCheckPaths := []string{
			"/",
			"/gw-scrm/login",
			"/gw-scrm/sso/login",
			"/gw-scrm/white-login",
			"/gw-scrm/open/auth",
			"/gw-scrm/current-version-info",
			"/gw-scrm/validate-code",
			"/gw-scrm/get-token",
			"/gw-scrm/get-industry-list",
			"/gw-scrm/work-wechat-login",
			"/gw-scrm/auto-access",
			"/gw-scrm/check-email",
			"/gw-scrm/get-bind-code",
			"/gw-scrm/get-code",
			"/gw-scrm/install-third-app",
			"/gw-scrm/check-email",
			"/gw-scrm/auto-access",
			"/gw-scrm/qy-code",
			"/gw-scrm/heart-beat",
			"/gw-scrm/swagger",
			"/gw-scrm/file/upload",
			"/gw-scrm/moment/js-config",
			"/gw-scrm/moment/appearance-get",
			"/gw-scrm/moment/member-topic",
			"/gw-scrm/source/link-detail",
			"/gw-scrm/chat/get-file-content",
			"/gw-scrm/file/private-by-key",
		}

		for _, p := range doNotCheckPaths {
			if path == p {
				return true
			}
		}
		return false
	}
}

// PathNotCheckAuth 不需要进行签名检查的请求
func PathNotCheckPower(c *gin.Context) bool {
	path := c.FullPath()
	doNotCheckPaths := []string{
		"/",
		"/gw-scrm/logout",
		"/gw-scrm/change-enterprise-list",
	}
	for _, p := range doNotCheckPaths {
		if path == p {
			return true
		}
	}

	return false
}

func frequencyControl(c *gin.Context, uri string) bool {
	if _, ok := fcUri[uri]; !ok {
		return true
	}
	baseLimit := fcUri[uri]
	baseKey := fmt.Sprintf("%s:fc_set:%s", config.CfgData.App, uri)
	baseLimit, err := redis.Client().Get(c, baseKey).Int()
	if err != nil && err != redis2.Nil {
		return true
	}
	if baseLimit == 0 {
		baseLimit = baseLimitNum
	}

	key := fmt.Sprintf("%s:fc_res:%s", config.CfgData.App, uri)
	reRet, err := redis.Client().Get(c, key).Int()
	if err != nil && err != redis2.Nil {
		return true
	}
	if reRet < baseLimit {
		redis.Client().Incr(c, key)
		redis.Client().Expire(c, key, time.Minute)
		return true
	}
	return false
}
