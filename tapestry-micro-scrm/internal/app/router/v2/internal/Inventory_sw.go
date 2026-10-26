package internal

import (
	"github.com/gin-gonic/gin"
	"micro-scrm/internal/app/http/handle/inventory"
)

// InventorySw   库存查询
func InventorySw(InventorySwGroup *gin.RouterGroup) {

	Handle := inventory.NewInvSwHandle()
	{
		// 根据skuItem查询库存
		InventorySwGroup.POST("/list", Handle.List)
		// 获取库存的查询搜索关键字
		InventorySwGroup.POST("/search-keywords", Handle.SearchKeywords)
		// 删除库存查询记录
		InventorySwGroup.POST("/del-keywords", Handle.DelKeywords)
	}

}
