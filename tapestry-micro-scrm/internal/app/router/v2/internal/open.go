package internal

import (
	"github.com/gin-gonic/gin"
	"micro-scrm/internal/app/http/handle/open"
	"micro-scrm/internal/app/http/middleware"
)

func Open(routerGroup *gin.RouterGroup) {
	oauthHandle := open.NewAuth()
	// 登陆授权
	routerGroup.POST("/login", oauthHandle.Login)

	limit := middleware.SetupAccountRateLimiter(20, 30)
	// 全局中间件（以下路由全部会走中间件）
	routerGroup.Use(middleware.AuthMw(), middleware.NewAccountRateMiddleware(limit))

	// 刷新token
	routerGroup.POST("/refresh_token", oauthHandle.RefreshToken)
	// token解密数据
	routerGroup.POST("/token_info", oauthHandle.TokenInfo)

	commonGroup := routerGroup.Group("common")
	{
		commonHandle := open.NewCommonHandle()
		// 批量解密 OpenAPI 返回的 PI 字段
		commonGroup.POST("decrypt", commonHandle.Decrypt)
	}

	// 三方门店预约.
	reserveGroup := routerGroup.Group("reserve")
	{
		reserveHandle := open.NewReserveHandle()
		// 分配预约导购
		reserveGroup.POST("assign_guide", reserveHandle.AssignReserveGuide)
		// 获取预约信息
		reserveGroup.POST("get_reserve", reserveHandle.GetReserve)
		// 用户获取预约列表
		reserveGroup.POST("user-lists", reserveHandle.UserLists)
	}

	assign := routerGroup.Group("assign")
	{
		reserveHandle := open.NewReserveHandle()
		// 分配导购
		assign.POST("store_worker", reserveHandle.AssignGuide)
	}

	// 导购、门店
	storeWorkerGroup := routerGroup.Group("store_worker")
	{
		storeWorkerHandle := open.NewStoreWorkerHandle()
		storeOperationHandle := open.NewStoreOperationHandle()
		// 获取门店主数据列表
		storeWorkerGroup.POST("store_master_list", storeOperationHandle.StoreMasterList)
		// 获取员工主数据列表
		storeWorkerGroup.POST("employee_master_list", storeOperationHandle.EmployeeMasterList)
		// 获取顾客的导购
		storeWorkerGroup.POST("sa_list_by_customer", storeWorkerHandle.SaListByCustomer)
		// 获取门店列表
		storeWorkerGroup.POST("store_list", storeWorkerHandle.StoreList)
		// 获取导购列表
		storeWorkerGroup.POST("sa_list", storeWorkerHandle.SaList)
		// 获取导购信息列表
		storeWorkerGroup.POST("store_work_list", storeWorkerHandle.GetStoreWorkList)
		// 获取客户添加的导购列表
		storeWorkerGroup.POST("customer_store_work_list", storeWorkerHandle.GetCustomerStoreWorkList)
		// 生成导购联系我二维码
		storeWorkerGroup.POST("gen_qrcode", storeWorkerHandle.GenQrCode)
		// dashboard组织架构树获取
		storeWorkerGroup.POST("area_tree", storeWorkerHandle.AreaTree)
		// 管理者权限获取
		storeWorkerGroup.POST("area_manager", storeWorkerHandle.AreaManager)
	}

	// 香港导购、门店
	hkStoreWorkerGroup := routerGroup.Group("hkmo_store_worker")
	{
		hkStoreWorkerHandle := open.NewHkStoreWorkerHandle()
		// 获取香港门店列表
		hkStoreWorkerGroup.POST("store_list", hkStoreWorkerHandle.StoreList)
		// 获取香港导购信息列表
		hkStoreWorkerGroup.POST("store_work_list", hkStoreWorkerHandle.GetStoreWorkList)
		// 香港管理者权限获取
		hkStoreWorkerGroup.POST("area_manager", hkStoreWorkerHandle.AreaManager)
	}

	// 发送待付款消息提醒
	messageSendGroup := routerGroup.Group("message-send")
	{
		messageSendHandle := open.NewMessageSendHandle()
		// 待付款消息提醒
		messageSendGroup.POST("wait-payment", messageSendHandle.SendWaitPayment)
	}
	// MA任务
	maTaskGroup := routerGroup.Group("ma-task")
	{
		maTaskHandle := open.NewMaTaskHandle()
		// MA获取任务模板列表
		maTaskGroup.POST("plan-template-lists", maTaskHandle.PlanTemplateLists)
		// MA三方创建任务
		maTaskGroup.POST("add-plan", maTaskHandle.AddPlan)
		// MA的task人群包
		maTaskGroup.POST("new-task-crowd", maTaskHandle.NewTaskCrowd)
	}

	cgeTaskGroup := routerGroup.Group("cge-task")
	{
		maTaskHandle := open.NewMaTaskHandle()
		cgeTaskGroup.POST("create-plan", maTaskHandle.CreatePlanByContent)
		cgeTaskGroup.POST("create-task", maTaskHandle.CreateTask)
		cgeTaskGroup.POST("push-task-crowd", maTaskHandle.PushTaskCrowd)
	}

	crmGroup := routerGroup.Group("crm")
	{
		crmHandle := open.NewCrm()
		crmGroup.POST("customer-clear", crmHandle.CustomerClear)
	}

	customerGroup := routerGroup.Group("customer")
	{
		customerHandle := open.NewCustomer()
		//批量更新优惠券状态
		customerGroup.POST("batch_update_coupon_status", customerHandle.BatchUpdateCouponStatus)
		//批量更新顾客触达时间
		customerGroup.POST("batch_update_interact_info", customerHandle.BatchUpdateActiveStatus)
		//批量更新顾客活跃沟通时间
		customerGroup.POST("batch_update_active_info", customerHandle.BatchUpdateActiveStatus)
	}
	customerSwGroup := routerGroup.Group("customer-sw")
	{
		customerSwHandle := open.NewCustomerSw()
		//批量更新SW优惠券状态
		customerSwGroup.POST("batch-update-coupon-status", customerSwHandle.BatchUpdateCouponStatus)
	}
	excitationGroup := routerGroup.Group("wejourney")
	{
		excitationHandle := open.NewExcitationHandle()
		excitationGroup.POST("member-info", excitationHandle.MemberInfo)

		excitationGroup.POST("add-order", excitationHandle.AddOrder)

		excitationGroup.POST("cancel-order", excitationHandle.CancelOrder)

		excitationGroup.POST("order-info", excitationHandle.OrderInfo)
	}

	sourceGroup := routerGroup.Group("source")
	{
		sourceHandle := open.NewSource()
		//这是给业务用的
		sourceGroup.POST("link-list", sourceHandle.LinkList)
		sourceGroup.POST("text-list", sourceHandle.TextList)
		sourceGroup.POST("image-list", sourceHandle.ImageList)
		sourceGroup.POST("video-list", sourceHandle.VideoList)
		sourceGroup.POST("friend-list", sourceHandle.FriendList)
		sourceGroup.POST("search-all", sourceHandle.SearchAll)

		//这是传数据用的
		sourceGroup.POST("list", sourceHandle.SourceList)
		sourceGroup.POST("group-tree", sourceHandle.SourceTree)
	}

	wecomCommonGroup := routerGroup.Group("wecom-common")
	{
		wecomCommonHandle := open.NewWecomCommon()
		wecomCommonGroup.POST("open-to-userid", wecomCommonHandle.ConvertOpenUserIdToUserId)
		wecomCommonGroup.POST("open-to-external-userid", wecomCommonHandle.ConvertOpenExternalUseridToExternalUserid)
		wecomCommonGroup.POST("get-info-by-source-external", wecomCommonHandle.GetInfoBySourceExternal)

	}
}
