package v2

import (
	"gw-scrm/internal/app/http/handle/activity_qrcode"
	"gw-scrm/internal/app/http/handle/customer"
	"gw-scrm/internal/app/http/handle/dashboard"
	"gw-scrm/internal/app/http/handle/department"
	"gw-scrm/internal/app/http/handle/display"
	"gw-scrm/internal/app/http/handle/drainage_contacts"
	"gw-scrm/internal/app/http/handle/enterprise"
	"gw-scrm/internal/app/http/handle/excitation"
	"gw-scrm/internal/app/http/handle/invite"
	"gw-scrm/internal/app/http/handle/notification"
	"gw-scrm/internal/app/http/handle/online_warehouse"
	"gw-scrm/internal/app/http/handle/plan"
	"gw-scrm/internal/app/http/handle/rbac"
	"gw-scrm/internal/app/http/handle/reserve"
	"gw-scrm/internal/app/http/handle/source"
	"gw-scrm/internal/app/http/handle/store_operation"
	"gw-scrm/internal/app/http/handle/store_visit"
	"gw-scrm/internal/app/http/handle/tags"
	"gw-scrm/internal/app/http/handle/user"

	"github.com/gin-gonic/gin"
)

func InitV2Router(router *gin.RouterGroup) {
	registerStoreOperationRoutes(router.Group("store-operation"))

	// 素材库
	sourceGroupRouter := router.Group("source")
	{
		categoryHandle := source.NewCategoryHandle()
		{
			sourceGroupRouter.POST("/category/list", categoryHandle.List)
			sourceGroupRouter.POST("/category/create", categoryHandle.Create)
			sourceGroupRouter.POST("/category/update", categoryHandle.Update)
			sourceGroupRouter.POST("/category/delete", categoryHandle.Delete)
			sourceGroupRouter.POST("/category/source_relation", categoryHandle.SourceRelation)
		}

		// 素材-图片编辑
		imageHandle := source.NewImageHandle()
		sourceGroupRouter.POST("/image/update", imageHandle.Update)
		sourceGroupRouter.POST("/image/detail", imageHandle.Detail)
		// 素材-视频编辑
		videoHandle := source.NewVideoHandle()
		sourceGroupRouter.POST("/video/update", videoHandle.Update)
		sourceGroupRouter.POST("/video/detail", videoHandle.Detail)

		// 素材-朋友圈
		friendHandle := source.NewFriendHandle()
		sourceGroupRouter.POST("/friend/create", friendHandle.Create)
		sourceGroupRouter.POST("/friend/update", friendHandle.Update)
		sourceGroupRouter.POST("/friend/detail", friendHandle.Detail)

		// 素材分组归属
		groupRelationHandle := source.NewGroupRelationHandle()
		{
			// 素材分组归属门店列表
			sourceGroupRouter.POST("/group-relation/list", groupRelationHandle.List)
			// 根据门店获取素材分组列表
			sourceGroupRouter.POST("/group-relation/get-list-for-store", groupRelationHandle.GetListForStore)
			// 素材分组归属门店保存
			sourceGroupRouter.POST("/group-relation/save", groupRelationHandle.Save)
		}

		// 素材类型
		sourceHandle := source.NewHandle()
		{
			// 后台管理查询各素材列表
			sourceGroupRouter.POST("/console-list", sourceHandle.ConsoleList)

			// 文本素材列表
			sourceGroupRouter.POST("/text-list", sourceHandle.TextList)
			// 图片素材列表
			sourceGroupRouter.POST("/image-list", sourceHandle.ImageList)
			// 链接素材列表
			sourceGroupRouter.POST("/link-list", sourceHandle.LinkList)
			// 小程序素材列表
			sourceGroupRouter.POST("/mini-program-list", sourceHandle.MiniProgramList)
			// 视频素材列表
			sourceGroupRouter.POST("/video-list", sourceHandle.VideoList)
			// 文本素材列表
			sourceGroupRouter.POST("/file-list", sourceHandle.FileList)
			// 各类素材详情
			sourceGroupRouter.POST("/source-detail", sourceHandle.SourceDetail)
			// 各类素材编辑保存
			sourceGroupRouter.POST("/source-update", sourceHandle.SourceUpdate)
			// 移除产品类别
			sourceGroupRouter.POST("/remove-category", sourceHandle.RemoveCategory)
			// 批量启用/禁用素材
			sourceGroupRouter.POST("/update-status", sourceHandle.UpdateStatus)
			// 验证素材在企微能否获取到media-id
			sourceGroupRouter.POST("/get-media-info", sourceHandle.GetMediaInfo)
			// 压缩文件
			sourceGroupRouter.POST("/compress-file", sourceHandle.CompressFile)
			// 置顶/取消置顶
			sourceGroupRouter.POST("/set-top", sourceHandle.SetTop)
		}
		// 品牌素材
		brandHandle := source.NewBrandHandle()
		{
			//获取商品图片
			sourceGroupRouter.POST("/brand/search-list", brandHandle.BrandSearchList)
			//图片AI检索
			sourceGroupRouter.POST("/brand/talkie-image-search", brandHandle.TalkieImageSearch)
			//sourceGroupRouter.POST("/brand/link-save", brandHandle.BrandLinkSave)
		}
	}

	customerGroupRouter := router.Group("customer")
	{
		customerHandle := customer.NewCustomerHandle()
		// 后台用户列表
		customerGroupRouter.POST("user-list", customerHandle.UserList)
		// 后台用户详情
		customerGroupRouter.POST("user-detail", customerHandle.UserDetail)

		// 离职继承
		customerInheritGroup := customerGroupRouter.Group("inherit")
		{
			customerInheritHandle := customer.NewInheritHandle()
			// 离职成员列表
			customerInheritGroup.POST("list", customerInheritHandle.List)
			// 分配离职成员的客户
			customerInheritGroup.POST("transfer_customer", customerInheritHandle.TransferCustomer)
			// 分配离职成员的客户群
			customerInheritGroup.POST("transfer_customer_group", customerInheritHandle.TransferCustomerGroup)
			// 获取待分配人员列表
			customerInheritGroup.POST("allot_external_list", customerInheritHandle.GetExternalAllotList)
			// 获取待分配群列表
			customerInheritGroup.POST("allot_group_list", customerInheritHandle.GetGroupChatAllotList)
			// 已分配导购列表
			customerInheritGroup.POST("takeover_list", customerInheritHandle.GetTakeOverList)
			// 接替外部联系人列表
			customerInheritGroup.POST("takeover_external_list", customerInheritHandle.GetTakeOverExternalList)
			// 接替群列表
			customerInheritGroup.POST("takeover_group_list", customerInheritHandle.GetTakeOverGroupList)
			// 同步离职待分配数据
			customerInheritGroup.POST("sync_customer_data", customerInheritHandle.SyncCustomerData)

		}
		// 播种任务模版
		sowTemplateGroup := customerGroupRouter.Group("sow-template")
		{
			sowTemplateHandle := customer.NewSowTaskHandle()
			sowTemplateGroup.POST("list", sowTemplateHandle.List)
			sowTemplateGroup.POST("save", sowTemplateHandle.Save)
			sowTemplateGroup.POST("detail", sowTemplateHandle.Detail)
			sowTemplateGroup.POST("update-status", sowTemplateHandle.UpdateStatus)
			sowTemplateGroup.POST("delete", sowTemplateHandle.Delete)
			sowTemplateGroup.POST("all-task-list", sowTemplateHandle.AllTaskList)
			sowTemplateGroup.POST("download-sow-task-record", sowTemplateHandle.DownloadSowTaskRecord)
			sowTemplateGroup.POST("template-list-export", sowTemplateHandle.TemplateListExport)
		}
	}

	// 标签路由组
	tagGroup := router.Group("tag")
	{
		// 印象标签
		impressionGroup := tagGroup.Group("impression")
		{
			impressionTagHandle := tags.NewImpressionHandle()
			impressionGroup.POST("/group/list", impressionTagHandle.GroupList)
			impressionGroup.POST("/tag/list", impressionTagHandle.TagList)
		}
	}

	// 到店邀约
	reserveGroup := router.Group("reserve")
	{
		// 预约到店订单
		orderGroup := reserveGroup.Group("order")
		{
			reserveHandle := reserve.NewReserveHandle()
			// 预约列表
			orderGroup.POST("/lists", reserveHandle.Lists)
			// 预约详情
			orderGroup.POST("/detail", reserveHandle.Detail)
			// 预约列表导出
			orderGroup.POST("/export-lists", reserveHandle.OrderListsExport)
			// 下载中心预约列表导出
			orderGroup.POST("/export-lists-new", reserveHandle.OrderListsExportNew)

			// 顾客受邀单列表
			orderGroup.POST("/reserve-order-list", reserveHandle.ReserveOrderList)
			// 顾客受邀单列表导出
			orderGroup.POST("/export-reserve-order-list", reserveHandle.ReserveOrderListExport)
			// 下载中心顾客受邀单列表导出
			orderGroup.POST("/export-reserve-order-list-new", reserveHandle.ReserveOrderListExportNew)
			// 管理后台受邀统计的订单明细导出
			orderGroup.POST("/export-order-record", reserveHandle.ExportOrderRecord)

		}
		// 预约服务
		serviceGroup := reserveGroup.Group("service")
		{
			reserveHandle := reserve.NewReserveHandle()
			// 新增服务
			serviceGroup.POST("/create", reserveHandle.ServiceCreate)
			// 修改
			serviceGroup.POST("/update", reserveHandle.ServiceUpdate)
			// 删除
			serviceGroup.POST("/delete", reserveHandle.ServiceDelete)
			// 列表
			serviceGroup.POST("/list", reserveHandle.ServiceList)
			// coach 检查是否可以删除预约服务
			serviceGroup.POST("/check-delete", reserveHandle.ServiceCheckDelete)
			// 服务详情接口
			serviceGroup.POST("/service-info", reserveHandle.ServiceInfo)
		}

		scoreGroup := reserveGroup.Group("service-score")
		{
			reserveHandle := reserve.NewReserveHandle()
			// 用户评分列表
			scoreGroup.POST("/list", reserveHandle.ScoreList)
			// 导出评价列表
			scoreGroup.POST("/export-list", reserveHandle.ScoreExport)
		}

	}

	// 企业配置
	enterpriseGroup := router.Group("enterprise")
	{
		enterpriseHandle := enterprise.NewEnterpriseHandle()
		enterpriseGroup.POST("/config", enterpriseHandle.EnterpriseConfig)
	}

	// 邀请函
	inviteGroup := router.Group("invite")
	{
		inviteHandle := invite.NewInviteHandle()
		// 创建或保存邀请函模板
		inviteGroup.POST("/template/save", inviteHandle.TemplateSave)
		// 邀请函模板列表
		inviteGroup.POST("/template/list", inviteHandle.TemplateList)
		// 邀请函模板详情信息
		inviteGroup.POST("/template/info", inviteHandle.TemplateInfo)
		// 删除邀请函模板
		inviteGroup.POST("/template/delete", inviteHandle.TemplateDelete)
		// 邀请函-操作上下架
		inviteGroup.POST("/template/update-status", inviteHandle.UpdateStatus)

		// 后台管理 邀请数据列表
		inviteGroup.POST("/record-order-list", inviteHandle.RecordOrderList)
		// 后台管理 邀请数据列表导出
		inviteGroup.POST("/record-order-list-export", inviteHandle.RecordOrderListExport)

	}

	departmentGroup := router.Group("departments")
	{
		departmentHandle := department.NewHandle()
		// 获取coach下的门店，筛选奥莱和零售
		departmentGroup.POST("/sa-store/lists", departmentHandle.StoreLists)
	}

	// 奖章路由
	excitationGroup := router.Group("excitation")
	{
		medalGroup := excitationGroup.Group("medal")
		{
			medalHandle := excitation.NewMedalHandle()
			medalGroup.POST("list", medalHandle.List)
			medalGroup.POST("create", medalHandle.Create)
			medalGroup.POST("detail", medalHandle.Detail)
			medalGroup.POST("update", medalHandle.Update)
			medalGroup.POST("remove", medalHandle.Remove)
			medalGroup.POST("status", medalHandle.Status)

			medalRecordHandle := excitation.NewMedalRecordHandle()
			// 领取记录列表
			medalGroup.POST("record-list", medalRecordHandle.List)
			// 勋章补发
			medalGroup.POST("record-reissue", medalRecordHandle.Reissue)
			// 勋章删除
			medalGroup.POST("record-delete", medalRecordHandle.Remove)

			// 奖章规则列表
			medalGroup.POST("rule-list", excitation.NewMedalRuleHandle().List)
		}

		// 积分规则
		integralRuleGroup := excitationGroup.Group("integral-rule")
		{
			integralRuleHandle := excitation.NewIntegralRuleHandle()
			// 积分规则列表
			integralRuleGroup.POST("list", integralRuleHandle.List)
			// 积分规则配置信息
			integralRuleGroup.POST("configuration", integralRuleHandle.ParameterConfiguration)
			// 积分规则添加
			integralRuleGroup.POST("create", integralRuleHandle.Create)
			// 积分规则详情
			integralRuleGroup.POST("detail", integralRuleHandle.Detail)
			// 积分规则删除
			integralRuleGroup.POST("destroy", integralRuleHandle.Destroy)
			// 积分规则编辑
			integralRuleGroup.POST("update", integralRuleHandle.Update)

			// 积分规则说明
			// 积分规则说明修改添加
			integralRuleGroup.POST("description-save", integralRuleHandle.DescriptionSave)
			// 积分规则说明信息
			integralRuleGroup.POST("description-detail", integralRuleHandle.DescriptionDetail)
		}

		// 成长激励后台维度配置
		incentiveGroup := excitationGroup.Group("incentive")
		{
			incentiveHandle := excitation.NewIncentiveHandle()
			// 积分规则维度配置
			incentiveGroup.POST("/config", incentiveHandle.IncentiveConfig)
			// 关系能力图谱-新增或修改指标配置
			incentiveGroup.POST("/relation/save", incentiveHandle.RelationSave)
			// 关系能力图谱-新增或修改指标配置
			incentiveGroup.POST("/relation/detail", incentiveHandle.RelationDetail)
			// 关系能力图谱-关联图谱规则详情列表
			incentiveGroup.POST("/relation/detail-list", incentiveHandle.RelationDetailList)
		}

		// 导购等级配置
		levelRuleGroup := excitationGroup.Group("level-rule")
		{
			levelRuleHandle := excitation.NewLevelRuleHandle()
			// 新增或修改等级规则
			levelRuleGroup.POST("/save", levelRuleHandle.Save)
			// 等级规则详情
			levelRuleGroup.POST("/detail-list", levelRuleHandle.DetailList)
			// 等级规则列表
			levelRuleGroup.POST("/list", levelRuleHandle.List)

		}
		pointsGroup := excitationGroup.Group("points")
		{
			PointsHandle := excitation.NewPointsHandle()
			// 积分获取记录列表
			pointsGroup.POST("/record-list", PointsHandle.RecordList)

			pointsGroup.POST("/download-record", PointsHandle.DownloadRecord)
		}

		// 导购等级数据
		memberDataGroup := excitationGroup.Group("member-data")
		{
			memberDataHandle := excitation.NewMemberDataHandle()
			// 后台列表
			memberDataGroup.POST("/list", memberDataHandle.List)
			// 后台列表 数据导出
			memberDataGroup.POST("/list-export", memberDataHandle.ListExport)
		}
		// 金币相关
		coinsGroup := excitationGroup.Group("coins")
		{
			coinsHandle := excitation.NewCoinsHandle()
			// 金币获取记录列表
			coinsGroup.POST("/record-list", coinsHandle.RecordList)
			// 金币配置获取
			coinsGroup.POST("/coins-settings-info", coinsHandle.CoinsSettingsInfo)
			// 金币配置保存
			coinsGroup.POST("/save-coins-settings", coinsHandle.SaveCoinsSettings)
			// 金币记录下载中心
			coinsGroup.POST("/download-record", coinsHandle.DownloadRecord)
		}
	}

	planGroup := router.Group("plan")
	{
		planHandle := plan.NewTaskHandle()
		planGroup.POST("/task-detail-export", planHandle.TaskDetailExport)

		// 后台管理子任务详情，数据明细导出---下载中心
		planGroup.POST("/task-detail-export-new", planHandle.TaskDetailExportNew)
		// 获取子任务的任务名称列表
		planGroup.POST("/task/plan-name-list", planHandle.TaskPlanNameList)
		// 计划任务的模板任务
		planTemplateGroup := planGroup.Group("template")
		{
			planTemplateHandle := plan.NewTemplateHandle()
			// 任务模板列表
			planTemplateGroup.POST("/lists", planTemplateHandle.List)
			// 任务模板创建
			planTemplateGroup.POST("/create", planTemplateHandle.Create)
			// 任务模板详情
			planTemplateGroup.POST("/details", planTemplateHandle.Details)
			planTemplateGroup.POST("/update", planTemplateHandle.Update)
			planTemplateGroup.POST("/status", planTemplateHandle.Status)
		}
		// 主任务任务数据
		planGroup.POST("/task/data/new-stat", planHandle.NewStat)
		// 主任务的新数据明细
		planGroup.POST("/task/data/stat/new-detail", planHandle.NewStatDetail)
		// 主任务明细数据的导出
		planGroup.POST("/task/data/stat/new-detail/export", planHandle.NewStatDetailExport)

		// 后台管理，主任务复制是判断子任务中是否有非1v1沟通类型的子任务
		planGroup.POST("/is-show", planHandle.IsShow)
	}

	displayGroup := router.Group("display")
	{
		taskGroup := displayGroup.Group("task")
		{
			taskHandle := display.NewTaskHandle()
			// VM陈列任务创建
			taskGroup.POST("add", taskHandle.AddTask)
			// VM陈列任务列表
			taskGroup.POST("list", taskHandle.List)
			// VM陈列任务查看
			taskGroup.POST("detail", taskHandle.Detail)
			// VM陈列任务复制
			taskGroup.POST("copy", taskHandle.Copy)
			// VM陈列任务上下架
			taskGroup.POST("update-show", taskHandle.UpdateShow)
			// VM陈列任务编辑
			taskGroup.POST("edit", taskHandle.UpdateTask)
			// VM陈列数据列表
			taskGroup.POST("allot-list", taskHandle.AllotList)
			// VM陈列数据的统计数
			taskGroup.POST("allot-num", taskHandle.AllotNum)
			// VM陈列数据的导出
			taskGroup.POST("allot-list-export", taskHandle.AllotListExport)
			// VM陈列消息提醒配置
			taskGroup.POST("notice-setting", taskHandle.NoticeSetting)
			// VM车裂消息的配置内容
			taskGroup.POST("notice-setting-info", taskHandle.NoticeSettingInfo)
			// 查看VM陈列数据
			taskGroup.POST("allot-info", taskHandle.AllotInfo)
			// VM陈列数据审核
			taskGroup.POST("allot-check", taskHandle.AllotCheck)
			// VM任务陈列数据发送模板消息 一键提醒
			taskGroup.POST("send-notice", taskHandle.AllotSendNotice)
			// VM陈列数据获取一键提醒的门店数
			taskGroup.POST("notice-num", taskHandle.AllotNoticeNum)
			// VM一键保存图片
			taskGroup.POST("save-image", taskHandle.SaveImage)
			// 获取VM列表用户
			taskGroup.POST("vm-list", taskHandle.VmList)
			// 导购头像上传至cos
			taskGroup.POST("/cos-avatar", taskHandle.CosAvatar)
			// 上传分配名单校验
			taskGroup.POST("/review-file-check", taskHandle.ReviewFileCheck)
			// 下载分配详情数据
			taskGroup.POST("/review-export", taskHandle.ReviewExport)
		}
	}

	notificationGroup := router.Group("notification")
	{
		handle := notification.NewInfoHandle()
		// FT&RO新增|编辑传讯内容库
		notificationGroup.POST("save", handle.Save)
		// FT&RO展示传讯内容库列表
		notificationGroup.POST("list", handle.List)
		// FT&RO的传讯内容详情
		notificationGroup.POST("info", handle.Info)
		// 部门列表
		notificationGroup.POST("departments", handle.Departments)
		// FT&RO的传讯内容复制
		notificationGroup.POST("copy", handle.Copy)
		// FT&RO传讯内容撤回
		notificationGroup.POST("recall", handle.Recall)
		// FT&RO传讯内容作废
		notificationGroup.POST("cancel", handle.Cancel)
		// FT&RO审核
		notificationGroup.POST("check", handle.Check)
		// 获取管理后台身份角色是RO还是FT
		notificationGroup.POST("manager", handle.Manager)

		// 白名单
		crowdGroup := notificationGroup.Group("crowd")
		{
			crowdHandle := notification.NewCrowdHandle()
			// 新增编辑
			crowdGroup.POST("save", crowdHandle.Save)
			// 详情
			crowdGroup.POST("detail", crowdHandle.Detail)
			// 删除
			crowdGroup.POST("delete", crowdHandle.Delete)
			// 列表
			crowdGroup.POST("list", crowdHandle.List)
			// 下载导入模板
			crowdGroup.GET("template", crowdHandle.Template)
			// 导入门店编码并解析员工
			crowdGroup.POST("import", crowdHandle.Import)
			// 导出未匹配门店编码
			crowdGroup.POST("export-failures", crowdHandle.ExportFailures)
		}

		// RO组合传讯
		infoMergeGroup := notificationGroup.Group("merge")
		{
			infoMergeHandle := notification.NewInfoMergeHandle()
			// 新增编辑
			infoMergeGroup.POST("save", infoMergeHandle.Save)
			// 详情
			infoMergeGroup.POST("detail", infoMergeHandle.Detail)
			// 列表
			infoMergeGroup.POST("list", infoMergeHandle.List)
			// 发布
			infoMergeGroup.POST("release", infoMergeHandle.Release)
			// 保存更新PDF
			infoMergeGroup.POST("save-pdf", infoMergeHandle.SavePdf)
		}

		// 传讯提醒配置
		noticeConfigGroup := notificationGroup.Group("notice-config")
		{
			noticeConfigHandle := notification.NewNoticeConfigHandle()
			// 传讯提醒配置编辑
			noticeConfigGroup.POST("/save", noticeConfigHandle.Save)
			// 传讯提醒配置查看
			noticeConfigGroup.POST("/detail", noticeConfigHandle.Detail)
		}
	}

	dashboardGroup := router.Group("dashboard")
	{
		areaGroup := dashboardGroup.Group("area")
		{
			areaHandle := dashboard.NewAreaHandle()
			// 组织架构树
			areaGroup.POST("tree", areaHandle.AreaTree)
			// 插入/更新节点信息
			areaGroup.POST("save", areaHandle.Save)
			// 删除节点
			areaGroup.POST("remove", areaHandle.Remove)
			// 待分配门店列表
			areaGroup.POST("not-relation-store-list", areaHandle.NotRelationStoreList)
			// //组织架构门店列表
			// areaGroup.POST("store-list", areaHandle.StoreList)
			// 区域列表（区域是到type=3的分区数据列表）
			areaGroup.POST("list", areaHandle.AreaList)
		}

		memberGroup := dashboardGroup.Group("manager")
		{
			memberHandle := dashboard.NewManagerHandle()
			// 更新导购管理权限
			memberGroup.POST("auth-save", memberHandle.AuthSave)
			// 删除导购管理权限
			memberGroup.POST("auth-delete", memberHandle.AuthDelete)
			// 获取管理者权限列表
			memberGroup.POST("auth-list", memberHandle.AuthList)
			// 更新导购管理权限
			memberGroup.POST("auth-detail", memberHandle.AuthDetail)
			// 所有已配置管理者id
			memberGroup.POST("all-member-ids", memberHandle.AuthManagerMemberIds)
		}

		hkManagerGroup := dashboardGroup.Group("hkmo-manager")
		{
			hkManagerHandle := dashboard.NewHkManagerHandle()
			hkManagerGroup.POST("store-list", hkManagerHandle.StoreList)
			hkManagerGroup.POST("auth-list", hkManagerHandle.AuthList)
			hkManagerGroup.POST("auth-save", hkManagerHandle.AuthSave)
			hkManagerGroup.POST("auth-detail", hkManagerHandle.AuthDetail)
			hkManagerGroup.POST("auth-delete", hkManagerHandle.AuthDelete)
			hkManagerGroup.POST("all-member-ids", hkManagerHandle.AllMemberIds)
		}

		hkStoreGroup := dashboardGroup.Group("hkmo-store")
		{
			hkStoreHandle := dashboard.NewHkStoreHandle()
			hkStoreGroup.POST("sync", hkStoreHandle.Sync)
			hkStoreGroup.POST("last-sync-time", hkStoreHandle.LastSyncTime)
			hkStoreGroup.POST("list", hkStoreHandle.List)
			hkStoreGroup.POST("detail", hkStoreHandle.Detail)
			hkStoreGroup.POST("save", hkStoreHandle.Save)
			hkStoreGroup.POST("close", hkStoreHandle.Close)
		}

		hkAreaGroup := dashboardGroup.Group("hkmo-area")
		{
			hkAreaHandle := dashboard.NewHkAreaHandle()
			hkAreaGroup.POST("list", hkAreaHandle.List)
		}

	}

	onlineWarehouseGroup := router.Group("online_warehouse")
	{
		onlineWareHandle := online_warehouse.NewOrderHandle()
		// 获取订单列表
		onlineWarehouseGroup.POST("/list", onlineWareHandle.List)
		// 获取订单详情
		onlineWarehouseGroup.POST("/info", onlineWareHandle.Info)
		// 门店订单列表
		onlineWarehouseGroup.POST("/store_list", onlineWareHandle.StoreList)
		// 订单确认提交
		onlineWarehouseGroup.POST("/submit", onlineWareHandle.Submit)
		// 订单编辑
		onlineWarehouseGroup.POST("/edit", onlineWareHandle.Edit)
		// 待发货的网仓订单导出
		onlineWarehouseGroup.POST("/export", onlineWareHandle.Export)
		// 待发货模板上传
		onlineWarehouseGroup.POST("/import_order", onlineWareHandle.ImportOrder)
		// 门店发货日期订单列表
		onlineWarehouseGroup.POST("/store_order_list", onlineWareHandle.StoreOrderList)
		// 门店批量上传转货单
		onlineWarehouseGroup.POST("/batch_transfer", onlineWareHandle.BatchTransfer)
		// 门店据单
		onlineWarehouseGroup.POST("/refund", onlineWareHandle.Refund)
		// 获取网仓订单的门店列表
		onlineWarehouseGroup.POST("/dep_list", onlineWareHandle.DepList)
		// 网仓订单门店的已发货待上传转货单的列表
		onlineWarehouseGroup.POST("/transfer_export", onlineWareHandle.TransferExport)
		// 网仓订单门店已发货的批量更新上传转货单
		onlineWarehouseGroup.POST("/transfer_import", onlineWareHandle.TransferImport)

	}

	storeVisitGroup := router.Group("store-visit")
	{
		storeVisitHandle := store_visit.NewStoreVisitMainHandle()

		// 巡店记录列表_list
		storeVisitGroup.POST("/list", storeVisitHandle.List)
		// 巡店记录关闭
		storeVisitGroup.POST("/close", storeVisitHandle.Close)
		// 获取所有部门列表
		storeVisitGroup.POST("/dept_name", storeVisitHandle.DeptName)
		// 巡店记录导出
		storeVisitGroup.POST("/export", storeVisitHandle.Export)
		// 巡店记录报表
		storeVisitGroup.POST("/dashboard", storeVisitHandle.Dashboard)
		// 巡店记录报表导出
		storeVisitGroup.POST("/dashboard/export", storeVisitHandle.DashboardExport)

	}
	activityQRCodeGroup := router.Group("activity-qrcode")
	{
		activityQRCodeHandle := activity_qrcode.NewActivityQRCodeMainHandle()
		// 保存活动二维码
		activityQRCodeGroup.POST("/activity-save", activityQRCodeHandle.ActivitySave)
		// 活动二维码详情
		activityQRCodeGroup.POST("/activity-detail", activityQRCodeHandle.ActivityDetail)
		// 活动二维码删除
		activityQRCodeGroup.POST("/activity-delete", activityQRCodeHandle.ActivityDelete)
		// 活动二维码列表
		activityQRCodeGroup.POST("/activity-list", activityQRCodeHandle.ActivityList)
		// 投放渠道列表
		activityQRCodeGroup.POST("/activity-category-list", activityQRCodeHandle.ActivityCategoryList)

		// 投放渠道欢迎语保存
		activityQRCodeGroup.POST("/welcome-save", activityQRCodeHandle.WelcomeSave)
		// 投放渠道欢迎语详情
		activityQRCodeGroup.POST("/welcome-list", activityQRCodeHandle.WelcomeList)
		// 投放渠道欢迎语列表
		activityQRCodeGroup.POST("/welcome-detail", activityQRCodeHandle.WelcomeDetail)
		// 投放渠道欢迎语删除
		activityQRCodeGroup.POST("/welcome-delete", activityQRCodeHandle.WelcomeDelete)

		// 管理后台统计
		activityQRCodeGroup.POST("/statistics", activityQRCodeHandle.Statistics)
	}

	userGroup := router.Group("sso")
	{
		userHandle := user.NewSSOHandle()
		// SSO登录
		userGroup.POST("/login", userHandle.SSOLogin)
	}

	rbacGroup := router.Group("rbac")
	{
		roleHandle := rbac.NewRoleHandle()
		rbacGroup.POST("/role/configurable-permissions", roleHandle.ConfigurablePermissions)
	}

	drainageContactCoachGroup := router.Group("drainage-contact-coach")
	{
		drainageContactCoachHandle := drainage_contacts.NewDrainageContactsCoachHandle()
		// 保存活动二维码
		drainageContactCoachGroup.POST("/save", drainageContactCoachHandle.Save)
		// 活动二维码详情
		drainageContactCoachGroup.POST("/detail", drainageContactCoachHandle.Detail)
		// 活动二维码删除
		drainageContactCoachGroup.POST("/delete", drainageContactCoachHandle.Delete)
		// 活动二维码列表
		drainageContactCoachGroup.POST("/list", drainageContactCoachHandle.List)
		// 投放渠道列表
		drainageContactCoachGroup.POST("/detail-list", drainageContactCoachHandle.DetailList)
		// 投放渠道列表
		drainageContactCoachGroup.POST("/zip-url", drainageContactCoachHandle.ZipUrl)

	}

}

func registerStoreOperationRoutes(group *gin.RouterGroup) {
	storeHandle := store_operation.NewStoreHandle()
	group.POST("/store/list", storeHandle.List)
	group.POST("/store/detail", storeHandle.Detail)
	group.POST("/store/update", storeHandle.Update)
	organizationHandle := store_operation.NewOrganizationHandle()
	group.POST("/organization/tree", organizationHandle.Tree)
	group.POST("/organization/store/list", organizationHandle.StoreList)
}
