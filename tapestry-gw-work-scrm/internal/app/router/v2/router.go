package v2

import (
	"gw-work-scrm/internal/app/http/handle/activity_qrcode"
	"gw-work-scrm/internal/app/http/handle/customer"
	"gw-work-scrm/internal/app/http/handle/display"
	"gw-work-scrm/internal/app/http/handle/invite"
	"gw-work-scrm/internal/app/http/handle/member"
	"gw-work-scrm/internal/app/http/handle/notification"
	"gw-work-scrm/internal/app/http/handle/online_warehouse"
	"gw-work-scrm/internal/app/http/handle/reserve"
	"gw-work-scrm/internal/app/http/handle/source"
	storeVisit "gw-work-scrm/internal/app/http/handle/store_visit"

	"github.com/gin-gonic/gin"
)

func InitV2Router(router *gin.RouterGroup) {
	// h5素材库
	sourceGroupRouter := router.Group("source")
	{
		// 产品类别列表
		categoryRouter := source.NewCategoryHandle()
		sourceGroupRouter.POST("/category/list", categoryRouter.List)

		// 素材分组归属
		groupRelationHandle := source.NewGroupRelationHandle()
		// h5素材分组归属门店列表
		sourceGroupRouter.POST("/group-relation/list", groupRelationHandle.List)
		// h5根据门店获取素材分组列表
		sourceGroupRouter.POST("/group-relation/get-list-for-store", groupRelationHandle.GetListForStore)
		// 素材分组归属门店保存
		// sourceGroupRouter.POST("/group-relation/save", groupRelationHandle.Save)
		// h5端门店员工根据查看分组权限查询各素材列表
		sourceGroupRouter.POST("/select-list", groupRelationHandle.SelectList)
		// h5各素材上新数量
		sourceGroupRouter.POST("/get-new-num", groupRelationHandle.GetNewSourceCount)
		// 我的收藏分类素材列表
		sourceGroupRouter.POST("/collect-select-list", groupRelationHandle.CollectSelectList)

		sourceHandle := source.NewSourceHandle()
		// 获取上传cos的凭证
		sourceGroupRouter.POST("/credential", sourceHandle.Credential)
		// H5端导购上传素材
		sourceGroupRouter.POST("/upload-media", sourceHandle.UploadMedia)
		// 验证导购是否有上传素材库的权限
		sourceGroupRouter.POST("/upload-auth", sourceHandle.CheckAuth)
		// 获取压缩结果
		sourceGroupRouter.POST("/compress", sourceHandle.Compress)
		// 验证是否用户首次打开素材交互提示
		sourceGroupRouter.POST("/check-first", sourceHandle.CheckFirst)
		// 上报计次
		sourceGroupRouter.POST("/report-times", sourceHandle.ReportTimes)
		// 素材收藏
		sourceGroupRouter.POST("/collect", sourceHandle.CollectEdit)

		brandHandle := source.NewBrandHandle()
		// 品牌素材库查询
		sourceGroupRouter.POST("/brand/search-list", brandHandle.BrandSearchList)
		// 品牌素材库获取图片缓存
		sourceGroupRouter.POST("/brand/get-cache-info", brandHandle.BrandGetCacheInfo)

		sourceGroupRouter.POST("/brand/collect/search-list", brandHandle.BrandCollectSearchList)
	}

	// 外部联系人路由组-
	customerGroup := router.Group("customer")
	{
		// 标签路由组
		customerTagsGroup := customerGroup.Group("tags")
		{
			customerTagsHandle := customer.NewTagsHandle()
			// h5-操作顾客的印象标签
			customerTagsGroup.POST("/impression/operate", customerTagsHandle.Operate)
		}

		// 客户
		customerRouter := customer.NewCustomerHandle()
		// 客户画像详情
		customerGroup.POST("/detail", customerRouter.Detail)
		// 修改客户备注
		customerGroup.POST("/update-remark", customerRouter.UpdateRemark)
		// 客户订单
		customerGroup.POST("/crm-order-list", customerRouter.CrmOrderList)
		// 客户足迹
		customerGroup.POST("/footprint-list", customerRouter.FootPrint)
		// 客户积分
		customerGroup.POST("/crm-points-list", customerRouter.CrmPointsList)
		// 客户卡劵
		customerGroup.POST("/crm-coupons-list", customerRouter.CrmCouponsList)
		// SW客户卡劵
		customerGroup.POST("/sw-user-crm-coupons-list", customerRouter.SwUserCrmCouponsList)

		// 客户加入的社群
		customerGroup.POST("/group-list", customerRouter.GroupChatList)
		// 编辑客户生日、地址、鞋码
		customerGroup.POST("/save-info-add", customer.NewInfoHandle().SaveInfoAdd)
		// 顾客的简单信息
		customerGroup.POST("/user-simple-info", customer.NewInfoHandle().SimpleInfo)

		// 播种(自设任务)
		customerSowTaskGroup := customerGroup.Group("sow-task")
		{
			sowTaskHandle := customer.NewSowTaskHandle()
			// 列表
			customerSowTaskGroup.POST("list", sowTaskHandle.List)
			// 创建
			customerSowTaskGroup.POST("create", sowTaskHandle.Create)
			// 完成
			customerSowTaskGroup.POST("finished", sowTaskHandle.Finished)
			// 删除
			customerSowTaskGroup.POST("delete", sowTaskHandle.Delete)
		}

		// h5-客户画像
		customerGroup.POST("/external-user-detail", customerRouter.CrmUserDetail)
		// h5-获取顾客邀约和播种数量
		customerGroup.POST("/user-detail-tips-num", customerRouter.UserDetailTipsNum)
		// h5-顾客和导购的绑定关系
		customerGroup.POST("/user-bind-relation", customerRouter.UserBindRelation)

	}

	reserveGroup := router.Group("reserve")
	{
		reserveHandle := reserve.NewOrderHandle()
		// 导购创建邀约单
		reserveGroup.POST("create-order", reserveHandle.CreateOrder)
		// 导购查看当前客户的邀约单列表
		reserveGroup.POST("order-list", reserveHandle.OrderList)
		// 判断时间段内的邀约数
		reserveGroup.POST("time-list", reserveHandle.TimeList)
		// 邀约服务列表
		reserveGroup.POST("service-list", reserveHandle.ServiceList)
		// 导购的所在门店信息
		reserveGroup.POST("member-info", reserveHandle.MemberDetail)
		// 导购的分享上报
		reserveGroup.POST("order-share", reserveHandle.OrderShare)
		// 导购取消小程序卡片发送，删除该邀约单
		reserveGroup.POST("order-delete", reserveHandle.OrderDelete)
		// SW&&ks验证同一时间段的邀约限制
		reserveGroup.POST("check-reserve-time", reserveHandle.CheckReserveTime)
		// 验证导购是否可以分享邀约单
		reserveGroup.POST("verify-share", reserveHandle.VerifyShare)
		// 导购确认订单是否到店购买,设置新客老客
		reserveGroup.POST("change-status", reserveHandle.ChangeStatus)
		// 邀约单分享了小程序卡片，更改为创建成功
		reserveGroup.POST("order-show", reserveHandle.OrderShow)
		// 顾客详情的邀约列表,展示顾客的所有邀约单列表信息
		reserveGroup.POST("customer-reserve-lists", reserveHandle.CustomerReserveLists)
	}

	inviteGroup := router.Group("invite")
	{
		inviteHandle := invite.NewInviteHandle()
		// 邀请函配置页列表
		inviteGroup.POST("config-list", inviteHandle.ConfigList)
		// 邀请函模板配置详情
		inviteGroup.POST("config-detail", inviteHandle.ConfigDetail)
		// 发送邀约单记录
		inviteGroup.POST("create-record", inviteHandle.CreateRecord)
		// 用户查看邀请函信息
		inviteGroup.POST("record-info", inviteHandle.RecordInfo)
		// 导购取消小程序卡片发送，删除该邀请函记录
		inviteGroup.POST("record-delete", inviteHandle.RecordDelete)

	}

	memberGroup := router.Group("member")
	{
		memberHandle := member.NewMemberHandle()
		// 导购身份和权限类型判定
		memberGroup.POST("guide-position", memberHandle.GuidePosition)
		// 导购的基本信息
		memberGroup.POST("base-info", memberHandle.BaseInfo)
		// 第一次进入某页面给的前端引导提示
		memberGroup.POST("first-open", memberHandle.FirstOpen)
		// 店长获取管理范围下的全部员工
		memberGroup.POST("guide-list", memberHandle.GuideList)
		// 获取店长管理的所有部门
		memberGroup.POST("store-or-department-list", memberHandle.StoreOrDepartmentList)
		// 获取财月财年
		memberGroup.POST("fiscal-month", memberHandle.GetFiscalMonth)

	}

	displayGroup := router.Group("display")
	{
		apiGroup := displayGroup.Group("api")
		{
			apiHandle := display.NewApiHandle()
			// 获取任务详情
			apiGroup.POST("task-info", apiHandle.TaskInfo)
			// 获取门店的任务详情内容
			apiGroup.POST("store-task-info", apiHandle.StoreTaskInfo)
			// 获取任务列表
			apiGroup.POST("task-list", apiHandle.TaskList)
			// 保存草稿箱
			apiGroup.POST("save", apiHandle.Save)
			// 导购查看VM备注
			apiGroup.POST("look", apiHandle.Look)
			// VM任务身份信息
			apiGroup.POST("auth", apiHandle.Auth)
			// 获取VM列表用户
			apiGroup.POST("vm-list", apiHandle.VmList)
			// 获取VM任务，审核成功的门店列表
			apiGroup.POST("store-list", apiHandle.StoreList)
			// 获取列表的tips数
			apiGroup.POST("list-tips", apiHandle.ListTips)
			// VM获取审核任务列表
			apiGroup.POST("check-list", apiHandle.CheckList)
			// VM获取任务详情的所有门店任务列表
			apiGroup.POST("check-store-list", apiHandle.CheckStoreList)
			// VM端 审核门店任务
			apiGroup.POST("vm-check", apiHandle.VmCheck)
			// VM端点击提醒
			apiGroup.POST("send-notice", apiHandle.SendNotice)
			// 获取welearningtoken
			apiGroup.POST("welearn-token", apiHandle.WeLearnToken)
			// VM端自动保存图片
			apiGroup.POST("save-image", apiHandle.SaveImage)
			// VM端自动删除图片
			apiGroup.POST("del-image", apiHandle.DelImage)
			// VM端自动保存备注
			apiGroup.POST("save-comment", apiHandle.SaveComment)
			// VM端自动删除备注
			apiGroup.POST("del-comment", apiHandle.DelComment)
			// 获取所有的任务列表（VM看板管理人员）
			apiGroup.POST("all-list", apiHandle.AllList)
			// VM看板数据的统计dashboard
			apiGroup.POST("vm-dashboard", apiHandle.DashBoard)
			// VM看板获取管理范围的门店最新任务未通过的数据
			apiGroup.POST("fail-list", apiHandle.FailList)
			// 下钻页面数据统计
			apiGroup.POST("dashboard-list", apiHandle.DashBoardList)
			// 获取门店下的任务列表
			apiGroup.POST("store-task-list", apiHandle.StoreTaskList)
			// pecial 保存草稿/提交
			apiGroup.POST("special-save", apiHandle.SpecialSave)
		}
		dashboardGroup := displayGroup.Group("dashboard")
		{
			dashboardHandle := display.NewDashboardHandle()
			// 组织架构树
			dashboardGroup.POST("tree", dashboardHandle.Tree)
			// 导购权限
			dashboardGroup.POST("member-auth", dashboardHandle.MemberAuth)
			// 搜索门店
			dashboardGroup.POST("search-store", dashboardHandle.SearchStore)
		}
	}

	notificationGroup := router.Group("notification")
	{
		notificationHandle := notification.NewInfoHandle()
		// 传讯列表
		notificationGroup.POST("info-list", notificationHandle.InfoList)
		// 组合传讯列表
		notificationGroup.POST("merge-list", notificationHandle.MergeList)
		// 组合传讯详情
		notificationGroup.POST("merge-detail", notificationHandle.MergeDetail)
		// 传讯详情
		notificationGroup.POST("info", notificationHandle.Info)
		// 撤回
		notificationGroup.POST("recall", notificationHandle.Recall)
		// 审批
		notificationGroup.POST("check", notificationHandle.Check)
		// 作废
		notificationGroup.POST("cancel", notificationHandle.Cancel)
		// 获取身份
		notificationGroup.POST("manager", notificationHandle.Manager)
		// 获取身份
		notificationGroup.POST("send-file-by-email", notificationHandle.SendFileByEmail)
		// 传讯未读数量
		notificationGroup.POST("list-read-num", notificationHandle.ListReadNum)
		// 保存pdf
		notificationGroup.POST("save-pdf", notificationHandle.SavePdf)

	}

	onlineWarehouseGroup := router.Group("online_warehouse")
	{
		onlineWarehouseHandle := onlineWarehouse.NewOrderHandle()
		// 订单列表
		onlineWarehouseGroup.POST("list", onlineWarehouseHandle.List)
		// 订单的各状态统计
		onlineWarehouseGroup.POST("tips_num", onlineWarehouseHandle.TipsNum)
		// 订单详情
		onlineWarehouseGroup.POST("info", onlineWarehouseHandle.Info)
		// 订单创建
		onlineWarehouseGroup.POST("create", onlineWarehouseHandle.Create)
		// 根据ross获取到的商品upCode,查询大仓的库存
		onlineWarehouseGroup.POST("/inventory/dc-list", onlineWarehouseHandle.DcList)
		// 根据ross获取到的商品upCode,查询大仓的库存
		onlineWarehouseGroup.POST("/inventory/list", onlineWarehouseHandle.InventoryList)
		// 门店的排行榜
		onlineWarehouseGroup.POST("/store-rank", onlineWarehouseHandle.StoreRank)
		// 导购的排行榜
		onlineWarehouseGroup.POST("/sa-rank", onlineWarehouseHandle.SaRank)
		// 导购导出文件发送邮件
		onlineWarehouseGroup.POST("/export-email", onlineWarehouseHandle.ExportMail)
		// 顾客端查看订单
		onlineWarehouseGroup.POST("/customer/order", onlineWarehouseHandle.CustomerInfo)
		// 顾客端提交订单
		onlineWarehouseGroup.POST("/customer/edit", onlineWarehouseHandle.CustomerEdit)

		// 离职导购的memberIds
		onlineWarehouseGroup.POST("/leave-sa", onlineWarehouseHandle.LeaveSa)
		// 离职导购的memberIds
		onlineWarehouseGroup.POST("/dep-list", onlineWarehouseHandle.DepList)
		// 网仓管理的所有导购
		onlineWarehouseGroup.POST("/sa-list", onlineWarehouseHandle.SaList)
		// 网仓订单取消为已废弃
		onlineWarehouseGroup.POST("/cancel", onlineWarehouseHandle.Cancel)
		// 获取网仓订单的分享海报media_id
		onlineWarehouseGroup.POST("/poster-media", onlineWarehouseHandle.PosterMedia)
		// 查询网仓订单的订单号是否存在
		onlineWarehouseGroup.POST("/check-order", onlineWarehouseHandle.CheckOrder)
		// 网仓订单看板
		onlineWarehouseGroup.POST("/dashboard", onlineWarehouseHandle.Dashboard)
		// 网仓订单商品排行
		onlineWarehouseGroup.POST("/goods-rank", onlineWarehouseHandle.GoodsRank)
	}

	storeVisitGroup := router.Group("store-visit")
	{
		storeVisitHandle := storeVisit.NewStoreVisitHandle()
		// 预约巡店
		storeVisitGroup.POST("/save-plan", storeVisitHandle.SavePlan)

		// 获取门店店长
		storeVisitGroup.POST("/get-store-manager", storeVisitHandle.GetStoreManager)

		// 获取门店店长
		storeVisitGroup.POST("/get-store-mm", storeVisitHandle.GetStoreMM)

		// 巡店列表
		storeVisitGroup.POST("/list", storeVisitHandle.StoreVisitList)

		// 预设部门及人员
		storeVisitGroup.POST("/dept-member-list", storeVisitHandle.DeptMemberList)

		// 个人邮箱通讯录
		storeVisitGroup.POST("/email-contact", storeVisitHandle.GetEmailContact)

		// 保存邮箱
		storeVisitGroup.POST("/save-email", storeVisitHandle.SaveEmail)

		// 新增/编辑邮箱
		storeVisitGroup.POST("/del-email", storeVisitHandle.DelEmail)

		// 搜索人员
		storeVisitGroup.POST("/search-member", storeVisitHandle.SearchMembers)

		// 同步指定人员
		storeVisitGroup.POST("/add-copy-member", storeVisitHandle.AddCopyMember)

		// 当前登录用户信息
		storeVisitGroup.POST("/login-user-info", storeVisitHandle.LoginUserInfo)

		// 开始巡店/巡店打卡
		storeVisitGroup.POST("/start-visit", storeVisitHandle.StartVisit)

		// 新增/修改反馈记录
		storeVisitGroup.POST("/save-show-record-content", storeVisitHandle.SaveShowRecordContent)

		// 新增反馈记录素材
		storeVisitGroup.POST("/save-show-record-media", storeVisitHandle.SaveShowRecordMedia)

		// 删除反馈记录
		storeVisitGroup.POST("/del-show-record", storeVisitHandle.DelShowRecord)

		// 删除反馈记录素材
		storeVisitGroup.POST("/del-show-record-media", storeVisitHandle.DelShowRecordMedia)

		// 反馈记录列表
		storeVisitGroup.POST("/show-record-list", storeVisitHandle.ShowRecordList)

		// 详情
		storeVisitGroup.POST("/info", storeVisitHandle.StoreVisitInfo)

		// 提交巡店记录（完成巡店）
		storeVisitGroup.POST("/complete-visit", storeVisitHandle.CompleteVisit)

		// 发送巡店报告通知
		storeVisitGroup.POST("/send-complete-report", storeVisitHandle.SendCompleteReport)

		// 保存更新 概览/整体表现/巡店报告pdf
		storeVisitGroup.POST("/save-visit-info", storeVisitHandle.SaveVisitInfo)

		// 保存更新 巡店报告查看权限及新增巡店人
		storeVisitGroup.POST("/save-report-permission", storeVisitHandle.SaveReportPermission)

		// 新增巡店留言
		storeVisitGroup.POST("/save-comment", storeVisitHandle.SaveComment)

		// 留言列表分页
		storeVisitGroup.POST("/comment-list", storeVisitHandle.CommentList)

		// 删除留言
		storeVisitGroup.POST("/del-comment", storeVisitHandle.DelComment)

		// 看板筛选树
		storeVisitGroup.POST("/dashboard-tree", storeVisitHandle.DashboardTree)

		// 看板导出
		storeVisitGroup.POST("/export-dashboard-list", storeVisitHandle.ExportDashboardList)

		// 发送巡店报告
		storeVisitGroup.POST("/send-pdf-to-email", storeVisitHandle.SendPdfToEmail)

		// 巡店配置获取
		storeVisitGroup.POST("/get-config", storeVisitHandle.GetConfig)

		// 保存用户语言配置
		storeVisitGroup.POST("/save-language-config", storeVisitHandle.SaveLanguageConfig)
	}

	activityQrCodeGroup := router.Group("activity-qrcode")
	{
		storeVisitHandle := activityQrCodeHandle.NewActivityQrCodeHandle()
		// 根据unionid获取有关系的导购二维码
		activityQrCodeGroup.POST("/get-qrcode-by-unionid", storeVisitHandle.GetQrCodeByUnionid)

		// 根据unionid获取有关系的导购二维码
		activityQrCodeGroup.POST("/get-all-store-list", storeVisitHandle.GetAllStoreList)

		// 活动信息获取
		activityQrCodeGroup.POST("/activity-detail", storeVisitHandle.ActivityDetail)
	}
}
