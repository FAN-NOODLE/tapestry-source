package internal

import (
	"github.com/gin-gonic/gin"
	dataCompass "micro-scrm/internal/app/http/handle/datacompass"
)

func DataCompass(group *gin.RouterGroup) {
	dataCompassHandle := dataCompass.NewChatGroupHandle()
	group.POST("chat-group/list", dataCompassHandle.List)
}
