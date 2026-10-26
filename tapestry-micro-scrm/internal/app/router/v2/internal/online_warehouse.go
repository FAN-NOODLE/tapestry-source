package internal

import (
	"github.com/gin-gonic/gin"
	onlineWarehouse "micro-scrm/internal/app/http/handle/online_warehouse"
)

// OnlineWarehouse 网仓订单路由
func OnlineWarehouse(groupRouter *gin.RouterGroup) {
	orderHandle := onlineWarehouse.NewHandle()
	// 网仓订单列表
	groupRouter.POST("/list", orderHandle.List)
	// 网仓订单发货提交
	groupRouter.POST("/submit", orderHandle.Submit)
	// 网仓订单详情
	groupRouter.POST("/info", orderHandle.Info)
	// 编辑网仓订单
	groupRouter.POST("/edit", orderHandle.Edit)
	// 门店数据列表
	groupRouter.POST("/store_list", orderHandle.StoreList)
	// 门店订单列表
	groupRouter.POST("/store_order_list", orderHandle.StoreOrderList)
	// 订单导出数据
	groupRouter.POST("/export", orderHandle.Export)
	// 上传发货订单数据
	groupRouter.POST("/import_order", orderHandle.ImportOrder)
	// 批量上传转货单
	groupRouter.POST("/batch_transfer", orderHandle.BatchTransfer)
	// 拒单
	groupRouter.POST("/refund", orderHandle.Refund)
	// 获取状态的统计数
	groupRouter.POST("/tips_num", orderHandle.TipsNum)
	// 创建销售单
	groupRouter.POST("/create", orderHandle.Create)

	// 门店的排行榜
	groupRouter.POST("/store-rank", orderHandle.StoreRank)
	// 导购的排行榜
	groupRouter.POST("/sa-rank", orderHandle.SaRank)
	// 导出文件发送至邮箱
	groupRouter.POST("/send-file-by-email", orderHandle.SendFileByEmail)
	// 顾客端获取订单信息
	groupRouter.POST("/customer-order-info", orderHandle.CustomerOrderInfo)
	// 顾客端提交订单信息
	groupRouter.POST("/customer-order-edit", orderHandle.CustomerOrderEdit)
	// 离职导购的memberIds
	groupRouter.POST("/leave-sa", orderHandle.LeaveSa)
	// 离职导购的memberIds
	groupRouter.POST("/dep-list", orderHandle.DepList)
	// 网仓管理的所有导购
	groupRouter.POST("/sa-list", orderHandle.SaList)
	// 网仓订单取消变更为废弃
	groupRouter.POST("/cancel", orderHandle.Cancel)
	// 网仓订单已发货的门店的待上传转货单
	groupRouter.POST("/transfer-export", orderHandle.TransferExport)
	// 网仓订单已发货的门店批量上传转货单
	groupRouter.POST("/transfer-import", orderHandle.TransferImport)
	// 获取网仓订单的海报mediaId
	groupRouter.POST("/poster-media", orderHandle.PosterMedia)
	// 查询网仓订单的订单号是否存在
	groupRouter.POST("/check-order", orderHandle.CheckOrderSn)
	// 网仓订单看板
	groupRouter.POST("/dashboard", orderHandle.Dashboard)
	// 商品的排行榜
	groupRouter.POST("/goods-rank", orderHandle.GoodsRank)

}
