package middleware

import (
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/pkg/common"
)

func Access() gin.HandlerFunc {
	return func(context *gin.Context) {
		// 处理前
		// TODO 前端固定传递企业hash
		context.Set(common.ContextKeyEnum.EnterpriseHash, context.GetHeader(common.HeaderKeyEnum.EnterpriseHash))

		requestID := context.GetHeader(common.HeaderKeyEnum.XRequestID)
		if requestID == "" {
			requestID = uuid.NewV4().String()
		}
		context.Set(string(ctx.RequestIDKey), requestID)

		context.Next()
		// 处理后
	}
}
