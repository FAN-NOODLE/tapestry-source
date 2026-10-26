package internal

import (
	"github.com/gin-gonic/gin"
	storeVisit "micro-scrm/internal/app/http/handle/store_visit"
)

// StoreVisit 巡店
func StoreVisit(groupRouter *gin.RouterGroup) {
	storeVisitHandle := storeVisit.NewHandle()
	// 获取门店店长
	groupRouter.POST("/get-store-manager", storeVisitHandle.GetStoreManager)

	// 获取门店MM
	groupRouter.POST("/get-store-mm", storeVisitHandle.GetStoreMM)

	// 预约巡店
	groupRouter.POST("/save-plan", storeVisitHandle.SavePlan)

	// 巡店列表
	groupRouter.POST("/list", storeVisitHandle.StoreVisitList)

	// 巡店详情
	groupRouter.POST("/info", storeVisitHandle.StoreVisitInfo)

	// 预设部门及人员
	groupRouter.POST("/dept-member-list", storeVisitHandle.DeptMemberList)

	// 个人邮箱通讯录
	groupRouter.POST("/email-contact", storeVisitHandle.GetEmailContact)

	// 保存邮箱
	groupRouter.POST("/save-email", storeVisitHandle.SaveEmail)

	// 新增/编辑邮箱
	groupRouter.POST("/del-email", storeVisitHandle.DelEmail)

	// 搜索人员
	groupRouter.POST("/search-member", storeVisitHandle.SearchMembers)

	// 同步指定人员
	groupRouter.POST("/add-copy-member", storeVisitHandle.AddCopyMember)

	// 当前登录用户信息
	groupRouter.POST("/login-user-info", storeVisitHandle.LoginUserInfo)

	// 开始巡店/巡店打卡
	groupRouter.POST("/start-visit", storeVisitHandle.StartVisit)

	// 新增/修改反馈记录
	groupRouter.POST("/save-show-record-content", storeVisitHandle.SaveShowRecordContent)

	// 新增反馈记录素材
	groupRouter.POST("/save-show-record-media", storeVisitHandle.SaveShowRecordMedia)

	// 删除反馈记录
	groupRouter.POST("/del-show-record", storeVisitHandle.DelShowRecord)

	// 删除反馈记录素材
	groupRouter.POST("/del-show-record-media", storeVisitHandle.DelShowRecordMedia)

	// 反馈记录列表
	groupRouter.POST("/show-record-list", storeVisitHandle.ShowRecordList)

	// 提交巡店记录（完成巡店）
	groupRouter.POST("/complete-visit", storeVisitHandle.CompleteVisit)

	// 发送巡店报告通知
	groupRouter.POST("/send-complete-report", storeVisitHandle.SendCompleteReport)

	// 保存更新 概览/整体表现/巡店报告pdf
	groupRouter.POST("/save-visit-info", storeVisitHandle.SaveVisitInfo)

	// 保存更新 巡店报告查看权限及新增巡店人
	groupRouter.POST("/save-report-permission", storeVisitHandle.SaveReportPermission)

	// 获取部门列表
	groupRouter.POST("/dept-name-list", storeVisitHandle.DeptNameList)

	// 新增巡店留言
	groupRouter.POST("/save-comment", storeVisitHandle.SaveComment)

	// 留言列表分页
	groupRouter.POST("/comment-list", storeVisitHandle.CommentList)

	// 删除留言
	groupRouter.POST("/del-comment", storeVisitHandle.DelComment)

	// 看板筛选树
	groupRouter.POST("/dashboard-tree", storeVisitHandle.DashboardTree)

	// 看板导出
	groupRouter.POST("/export-dashboard-list", storeVisitHandle.ExportDashboardList)

	// 发送巡店报告pdf
	groupRouter.POST("/send-pdf-to-email", storeVisitHandle.SendPdfToEmail)

	// 巡店配置获取
	groupRouter.POST("/get-config", storeVisitHandle.GetConfig)

	// 保存用户语言切换
	groupRouter.POST("/save-language-config", storeVisitHandle.SaveLanguageConfig)

	adminGroupRouter := groupRouter.Group("admin")
	{
		// 巡店列表——list
		adminGroupRouter.POST("list", storeVisitHandle.AdminStoreVisitList)
		// 巡店关闭
		adminGroupRouter.POST("close", storeVisitHandle.AdminStoreVisitClose)

		adminGroupRouter.POST("export", storeVisitHandle.AdminStoreVisitExport)
		// 巡店统计
		adminGroupRouter.POST("dashboard", storeVisitHandle.AdminDashboard)
		// 巡店统计导出
		adminGroupRouter.POST("dashboard/export", storeVisitHandle.AdminDashboardExport)

	}

}
