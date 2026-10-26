package internal

import (
	"github.com/gin-gonic/gin"
	"micro-scrm/internal/app/http/handle/notification"
)

// Notification 传讯任务路由
func Notification(groupRouter *gin.RouterGroup) {
	notificationHandle := notification.NewHandle()
	// FT&RO的传讯任务创建和编辑
	groupRouter.POST("/save", notificationHandle.Save)
	// FT&RO的传讯任务列表
	groupRouter.POST("/list", notificationHandle.List)
	// FT&RO的传讯任务查看
	groupRouter.POST("/info", notificationHandle.Info)
	// 获取部门列表
	groupRouter.POST("/departments", notificationHandle.Departments)
	// FT&RO的传讯复制
	groupRouter.POST("/copy", notificationHandle.Copy)
	// FT&RO的传讯撤回
	groupRouter.POST("/recall", notificationHandle.Recall)
	// FT&RO的传讯作废
	groupRouter.POST("/cancel", notificationHandle.Cancel)
	// FT&RO的传讯审核
	groupRouter.POST("/check", notificationHandle.Check)
	// 获取当前身份信息是FT还是RO
	groupRouter.POST("/manager", notificationHandle.Manager)
	//使用邮件发送附件
	groupRouter.POST("/send-file-by-email", notificationHandle.SendFileByEmail)
	// 未读标识数量
	groupRouter.POST("/list-read-num", notificationHandle.ListReadNum)
	// 白名单
	crowdGroup := groupRouter.Group("crowd")
	{
		crowdHandle := notification.NewCrowdHandle()
		// 白名单创建编辑
		crowdGroup.POST("/save", crowdHandle.Save)
		// 白名单列表
		crowdGroup.POST("/list", crowdHandle.List)
		// 白名单查看
		crowdGroup.POST("/detail", crowdHandle.Detail)
		// 白名单删除
		crowdGroup.POST("/delete", crowdHandle.DeleteCrowd)
	}

	// RO组合传讯
	infoMergeGroup := groupRouter.Group("merge")
	{
		infoMergeHandle := notification.NewInfoMergeHandle()
		// RO组合传讯创建编辑
		infoMergeGroup.POST("/save", infoMergeHandle.Save)
		// RO组合传讯查看
		infoMergeGroup.POST("/detail", infoMergeHandle.Detail)
		// RO组合传讯列表
		infoMergeGroup.POST("/list", infoMergeHandle.List)
		// RO组合传讯发布
		infoMergeGroup.POST("/release", infoMergeHandle.Release)
		// RO组合传讯H5列表
		infoMergeGroup.POST("/list-show", infoMergeHandle.ListShow)
		// 保存PDF
		infoMergeGroup.POST("/save-pdf", infoMergeHandle.SavePdf)

	}

	// 传讯提醒配置
	noticeConfigGroup := groupRouter.Group("notice-config")
	{
		noticeConfigHandle := notification.NewNoticeConfigHandle()
		// 传讯提醒配置编辑
		noticeConfigGroup.POST("/save", noticeConfigHandle.Save)
		// 传讯提醒配置查看
		noticeConfigGroup.POST("/detail", noticeConfigHandle.Detail)
	}

}
