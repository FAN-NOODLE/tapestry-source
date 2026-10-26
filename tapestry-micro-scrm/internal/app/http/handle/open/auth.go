package open

import (
	"errors"
	"github.com/gin-gonic/gin"
	"micro-scrm/internal/app"
	openDao "micro-scrm/internal/app/dao/open"
	"micro-scrm/internal/app/http/handle"
	"micro-scrm/internal/app/http/middleware"
	openReq "micro-scrm/internal/app/http/request/open"
	"micro-scrm/internal/config"
	"micro-scrm/internal/pkg/jwt"
	"strings"
	"time"
)

type Auth struct {
	handle.BaseHandle
}

func NewAuth() Auth {
	return Auth{}
}

// Login 登陆授权
func (h Auth) Login(c *gin.Context) {
	cc := h.WarpContext(c)

	var request openReq.LoginRequest

	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	openUserInfo, err := openDao.NewOpenUser().GetInfo(h.WarpContext(c),
		map[string]interface{}{
			"secret_id": request.SecretId,
		})

	if err != nil || openUserInfo.ID == 0 {
		h.ResError(c, errors.New("账户错误"))
		if err != nil {
			app.Logger().Error(cc).Msgf("查询账号信息失败,错误:%v", err)
		}
		return
	}

	// 验证密码
	if request.SecretKey != openUserInfo.SecretKey {
		h.ResError(c, errors.New("账户或密码错误"))
		return
	}

	// 生成 token
	token, err := jwt.New().CreateToken(jwt.NewCustomClaims(cc, int64(openUserInfo.SecretID), openUserInfo.EnterpriseID))
	if err != nil {
		app.Logger().Error(cc).Msgf("生成token失败,错误:%v", err)
		h.ResError(c, errors.New("生成token失败"))
		return
	}

	h.ResSuccess(c, gin.H{
		"secret_id":  request.SecretId,
		"token_type": jwt.TokenHeadName,
		"token":      token,
		"expires_in": time.Now().Add(time.Duration(config.CfgData.JWT.Ttl) * time.Second).Unix(),
	})
	return
}

// TokenInfo token解密数据
func (h Auth) TokenInfo(c *gin.Context) {
	authHeader := c.Request.Header.Get(middleware.Authorization)
	parts := strings.SplitN(authHeader, " ", 2)
	parseToken, err := jwt.New().ParseToken(parts[1])
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, parseToken)
	return
}

// RefreshToken 刷新token
func (h Auth) RefreshToken(c *gin.Context) {
	authHeader := c.Request.Header.Get(middleware.Authorization)
	parts := strings.SplitN(authHeader, " ", 2)

	h.Log().Info(h.WarpContext(c)).Msgf("登录用户获取token: %s", parts[1])

	token, err := jwt.New().RefreshToken(parts[1])
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, gin.H{
		"token":      token,
		"expires_in": time.Now().Add(2 * time.Hour).Unix(),
	})
	return
}
