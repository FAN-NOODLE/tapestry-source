package internal

import (
	"github.com/gin-gonic/gin"
	"micro-scrm/internal/app/http/handle/inventory"
)

// Inventory 库存查询
func Inventory(InventoryGroup *gin.RouterGroup) {

	Handle := inventory.NewInvHandle()
	{
		// 根据skuItem查询库存
		InventoryGroup.POST("/list", Handle.List)
		// 根据skuItem查询库存（只查大仓库存）
		InventoryGroup.POST("/dc-list", Handle.DcList)
		// 获取库存的查询搜索关键字
		InventoryGroup.POST("/search-keywords", Handle.SearchKeywords)
		// 删除库存查询记录
		InventoryGroup.POST("/del-keywords", Handle.DelKeywords)
	}

}
