package jwt

import (
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"gw-work-scrm/internal/config"
	"gw-work-scrm/internal/pkg/errs"
)

const TokenHeadName = "Bearer"

// noJwtRoutes 无需走JWT的路由集合
var noJwtRoutes = map[string]bool{
	"/auth/login":                            false,
	"/auth/mock":                             false,
	"/auth/config":                           false,
	"/auth/agent-config":                     false,
	"/auth/redirect_url":                     false,
	"/heart-beat":                            false,
	"/":                                      false,
	"/swagger":                               false,
	"/online_warehouse/customer/order":       false,
	"/online_warehouse/customer/edit":        false,
	"/activity-qrcode/get-qrcode-by-unionid": false,
	"/activity-qrcode/get-all-store-list":    false,
	"/activity-qrcode/activity-detail":       false,
}

// AuthMw 中间件，检查token
func AuthMw() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 未匹配到相关路由，不做校验
		if c.FullPath() == "" {
			return
		}

		path := strings.Replace(c.FullPath(), config.CfgData.Restful.BasePath, "", 1)
		if _, ok := noJwtRoutes[path]; ok {
			return
		}

		// 验证 JWT
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			ginplus.ResError(c, errs.TokenNotFind)
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == TokenHeadName) {
			ginplus.ResError(c, errs.AuthTokenIllegal)
			c.Abort()
			return
		}
		token := parts[1]

		log.Info(ctx.Wrap(c)).Msgf("登录用户获取token: %s", token)

		j := New()
		// parseToken 解析token包含的信息
		claims, err := j.ParseToken(token)
		if err != nil {
			if err == TokenExpired {
				ginplus.ResError(c, errs.TokenIsExp)
				c.Abort()
				return
			}
			ginplus.ResError(c, err)

			c.Abort()
			return
		}

		// 验证信息
		memberInfo, err := ViableMembers(ctx.Wrap(c), c.FullPath(), claims)
		if err != nil {
			ginplus.ResError(c, err)
			c.Abort()
			return
		}

		if memberInfo == nil {
			ginplus.ResError(c, errs.UserIdNotFind)
			c.Abort()
			return
		}

		// 继续交由下一个路由处理,并将解析出的信息传递下去
		c.Set("enterpriseId", claims.EnterpriseID)
		c.Set("memberId", memberInfo.ID)
		c.Set("memberInfo", memberInfo)
		c.Set("openId", claims.OpenId)
		c.Set("userId", memberInfo.Userid)
	}
}
