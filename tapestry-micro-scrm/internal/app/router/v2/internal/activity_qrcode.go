package internal

import (
	"github.com/gin-gonic/gin"
	activityQRCode "micro-scrm/internal/app/http/handle/activity_qrcode"
)

// ActivityQRCode  活动二维码
func ActivityQRCode(groupRouter *gin.RouterGroup) {
	activityQRCodeHandle := activityQRCode.NewHandle()
	// 活码活动保存保存
	groupRouter.POST("/activity-save", activityQRCodeHandle.ActivitySave)

	// 活码活动详情
	groupRouter.POST("/activity-list", activityQRCodeHandle.ActivityList)

	// 活码活动列表
	groupRouter.POST("/activity-detail", activityQRCodeHandle.ActivityDetail)

	// 活码活动删除
	groupRouter.POST("/activity-delete", activityQRCodeHandle.ActivityDelete)

	// 投放渠道列表
	groupRouter.POST("/activity-category-list", activityQRCodeHandle.ActivityCategoryList)

	// 根据unionid获取有关系的导购二维码
	groupRouter.POST("/get-qrcode-by-unionid", activityQRCodeHandle.GetQrCodeByUnionid)

	// 根据unionid获取有关系的导购二维码
	groupRouter.POST("/get-all-store-list", activityQRCodeHandle.GetAllStoreList)

	// 投放渠道欢迎语保存
	groupRouter.POST("/welcome-save", activityQRCodeHandle.WelcomeSave)

	// 投放渠道欢迎语详情
	groupRouter.POST("/welcome-list", activityQRCodeHandle.WelcomeList)

	// 投放渠道欢迎语列表
	groupRouter.POST("/welcome-detail", activityQRCodeHandle.WelcomeDetail)

	// 投放渠道欢迎语删除
	groupRouter.POST("/welcome-delete", activityQRCodeHandle.WelcomeDelete)

	// 管理后台统计
	groupRouter.POST("/statistics", activityQRCodeHandle.Statistics)
}
