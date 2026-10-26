package internal

import (
	"github.com/gin-gonic/gin"
	"micro-scrm/internal/app/http/handle/marketing"
)

func MarketingStrategies(marketingGroup *gin.RouterGroup) {
	customCard := marketing.NewCustomCard()
	marketingGroup.POST("/template/list", customCard.Lists)
	marketingGroup.POST("/template/save", customCard.Create)
	marketingGroup.POST("/template/delete", customCard.Delete)
	marketingGroup.POST("/template/check", customCard.Check)
	marketingGroup.POST("/template/detail", customCard.Detail)
}
