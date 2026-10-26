package internal

import (
	"github.com/gin-gonic/gin"
	"micro-scrm/internal/app/http/handle/sso"
)

func SSO(ssoGroup *gin.RouterGroup) {
	newSSOHandle := ssoHandle.NewSSOHandle()
	{
		// SSO登录
		ssoGroup.POST("/login", newSSOHandle.Login)
	}
}
