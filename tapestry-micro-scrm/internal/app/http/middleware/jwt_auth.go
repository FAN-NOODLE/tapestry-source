package middleware

import (
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"micro-scrm/internal/app"
	openDao "micro-scrm/internal/app/dao/open"
	openModel "micro-scrm/internal/app/models/open"
	"micro-scrm/internal/config"
	"micro-scrm/internal/pkg/errs"
	"micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/jwt"
	"strings"
)

const (
	TokenHeadName = "Bearer"
	Authorization = "Authorization"
)

const (
	EnterpriseId = "enterpriseId"
	UserId       = "userId"
)

// whiteRoutes 无需走JWT的路由集合
var whiteRoutes = map[string]bool{
	"/":           false,
	"/open/login": false,
}

// AuthMw 中间件，检查token
func AuthMw() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 根据Authorization判断是否需要走授权
		if !runAuth(c) {
			return
		}

		claims, err := ParseToken(c)
		if err != nil {
			ginplus.ResError(c, err)
			c.Abort()
			return
		}
		app.Logger().Info(ctx.Wrap(c)).Msgf(
			"open用户%v,访问url%s,clientIP:%s,remoteAddr:%s,xForwardedFor:%s,xRealIP:%s",
			claims,
			c.Request.URL,
			c.ClientIP(),
			c.Request.RemoteAddr,
			c.GetHeader("X-Forwarded-For"),
			c.GetHeader("X-Real-IP"),
		)
		// 继续交由下一个路由处理,并将解析出的信息传递下去
		c.Set(EnterpriseId, claims.EnterpriseID)
		c.Set(UserId, claims.ID)
		// 检查api权限
		check, err := apiCheck(c, claims.ID)
		if !check {
			ginplus.ResError(c, err)
			c.Abort()
			return
		}
	}
}

// runAuth 判断是否需要走授权
func runAuth(c *gin.Context) bool {
	// 未匹配到相关路由，不做校验
	if c.FullPath() == "" {
		return false
	}

	path := strings.Replace(c.FullPath(), config.CfgData.Restful.BasePath, "", 1)
	if _, ok := whiteRoutes[path]; ok {
		return false
	}

	return true
}

// apiCheck 判断token是否有api权限
func apiCheck(c *gin.Context, id int64) (bool, error) {
	info, err := openDao.NewOpenUser().GetInfo(ctx.Wrap(c), map[string]interface{}{
		"secret_id": id,
	})
	if err != nil {
		return false, err
	}
	if info.AuthType == openModel.SuperUser {
		return true, nil
	}
	url := strings.TrimPrefix(c.Request.URL.Path, "/micro-scrm/v1/open/")
	apiListSlice := strings.Split(info.ApiList, ",")
	if helper.InArray(url, apiListSlice) {
		return true, nil
	}
	return false, errs.PermissionDenied
}

// ParseToken 解析token
func ParseToken(c *gin.Context) (*jwt.CustomClaims, error) {
	authHeader := c.Request.Header.Get(Authorization)
	if authHeader == "" {
		return nil, errs.TokenNotFind
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != TokenHeadName {
		return nil, errs.AuthTokenIllegal
	}

	token := parts[1]
	app.Logger().Info(ctx.Wrap(c)).Msgf("登录用户获取token: %s", token)

	claims, err := jwt.New().ParseToken(token)
	if err != nil {
		if err == jwt.TokenExpired {
			return nil, errs.TokenIsExp
		}
		return nil, errs.AuthTokenIllegal
	}

	return claims, nil
}
