package v2

import (
	"github.com/gin-gonic/gin"
	"micro-scrm/internal/app/router/v2/internal"
)

func InitRouter(router *gin.RouterGroup) {
	// 素材库
	internal.Source(router.Group("source"))

	// 企业微信回调路由组
	internal.Callback(router.Group("callback"))

	// 标签路由组
	internal.Tags(router.Group("tag"))

	// 外部联系人路由组
	internal.Customer(router.Group("customer"))

	// 计划任务
	internal.Plan(router.Group("plan"))

	// 营销策略
	internal.MarketingStrategies(router.Group("marketing-strategies"))

	// 预约
	internal.Reserve(router.Group("reserve"))

	// 导购信息
	internal.Member(router.Group("member"))

	// 三方提供接口
	internal.Open(router.Group("open"))

	// 邀请函
	internal.Invite(router.Group("invite"))
	// coach 门店的接口
	internal.Department(router.Group("department"))

	// 成长激励
	internal.Excitation(router.Group("excitation"))

	// 数据罗盘.群组统计
	internal.DataCompass(router.Group("data-compass"))

	// 数据罗盘.群组统计
	internal.Dashboard(router.Group("dashboard"))

	// 库存查询
	internal.Inventory(router.Group("inventory"))

	// Sw库存查询
	internal.InventorySw(router.Group("inventory-sw"))

	// VM陈列任务
	internal.Display(router.Group("display"))

	// 传讯任务
	internal.Notification(router.Group("notification"))
	// 网仓订单
	internal.OnlineWarehouse(router.Group("online_warehouse"))

	// 巡店
	internal.StoreVisit(router.Group("store-visit"))

	// 活动二维码
	internal.ActivityQRCode(router.Group("activity-qrcode"))

	// 联系我活动coach
	internal.DrainageContact(router.Group("drainage-contact-coach"))

	// SSO登录
	internal.SSO(router.Group("sso"))
}
