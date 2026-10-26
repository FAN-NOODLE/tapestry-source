package middleware

import (
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/pkg/common"
	"strconv"
)

// Access 访问
func Access() gin.HandlerFunc {
	return func(context *gin.Context) {
		// 处理前
		enterpriseID, _ := strconv.ParseInt(context.GetHeader(common.HeaderKeyEnum.EnterpriseID), 10, 64)
		enterpriseHash := context.GetHeader(common.HeaderKeyEnum.EnterpriseHash)
		context.Set(common.ContextKeyEnum.EnterpriseID, enterpriseID)
		context.Set(common.ContextKeyEnum.EnterpriseHash, enterpriseHash)
		requestID := context.GetHeader(common.HeaderKeyEnum.XRequestID)
		if requestID == "" {
			requestID = uuid.NewV4().String()
		}
		context.Set(string(ctx.RequestIDKey), requestID)

		context.Next()
		// 处理后
	}
}
