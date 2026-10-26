package internal

import (
	"github.com/gin-gonic/gin"
	"micro-scrm/internal/app/http/handle/dashboard"
)

func Dashboard(group *gin.RouterGroup) {
	// dashboard区域
	areaGroup := group.Group("area")
	{
		areaHandle := dashboard.NewAreaHandle()
		// 组织架构树
		areaGroup.POST("tree", areaHandle.Tree)
		// 导购节点权限
		areaGroup.POST("manager-auth", areaHandle.ManagerAuth)
		// 插入/更新节点信息
		areaGroup.POST("save", areaHandle.Save)
		// 删除节点
		areaGroup.POST("remove", areaHandle.Remove)
		// 待分配门店列表
		areaGroup.POST("not-relation-store-list", areaHandle.NotRelationStoreList)
		// 闭店或者重新打开
		areaGroup.POST("store-close-or-open", areaHandle.StoreCloseOrOpen)
		// 闭店或者重新打开
		areaGroup.POST("search-store", areaHandle.SearchStore)
		// 区域列表（区域是到type=3的分区数据列表）
		areaGroup.POST("list", areaHandle.AreaList)
	}

	// dashboard管理者
	memberGroup := group.Group("manager")
	{
		memberHandle := dashboard.NewManagerHandle()
		// 获取管理者权限列表
		memberGroup.POST("auth-list", memberHandle.MemberAuthList)
		// 新增更新用户权限
		memberGroup.POST("auth-save", memberHandle.SaveManagerAuth)
		// 删除已配置管理者
		memberGroup.POST("auth-delete", memberHandle.DeleteManagerAuth)
		// 获取用户权限
		memberGroup.POST("auth-detail", memberHandle.MemberAuthDetail)
		// 所有已配置管理者id
		memberGroup.POST("all-member-ids", memberHandle.GetAllManagerIds)
	}

	// dashboard香港门店管理者
	hkMemberGroup := group.Group("hkmo-manager")
	{
		hkMemberHandle := dashboard.NewHkManagerHandle()
		// 获取香港门店列表
		hkMemberGroup.POST("store-list", hkMemberHandle.StoreList)
		// 获取香港管理者权限列表
		hkMemberGroup.POST("auth-list", hkMemberHandle.MemberAuthList)
		// 新增更新香港用户权限
		hkMemberGroup.POST("auth-save", hkMemberHandle.SaveManagerAuth)
		// 删除已配置香港管理者
		hkMemberGroup.POST("auth-delete", hkMemberHandle.DeleteManagerAuth)
		// 获取香港用户权限
		hkMemberGroup.POST("auth-detail", hkMemberHandle.MemberAuthDetail)
		// 所有已配置香港管理者id
		hkMemberGroup.POST("all-member-ids", hkMemberHandle.GetAllManagerIds)
	}

	// dashboard香港门店
	hkStoreGroup := group.Group("hkmo-store")
	{
		hkStoreHandle := dashboard.NewHkStoreHandle()
		// 同步企微香港门店
		hkStoreGroup.POST("sync", hkStoreHandle.Sync)
		// 香港门店最后同步时间
		hkStoreGroup.POST("last-sync-time", hkStoreHandle.LastSyncTime)
		// 香港门店列表
		hkStoreGroup.POST("list", hkStoreHandle.List)
		// 香港门店详情
		hkStoreGroup.POST("detail", hkStoreHandle.Detail)
		// 更新香港门店
		hkStoreGroup.POST("save", hkStoreHandle.Save)
		// 关闭香港门店
		hkStoreGroup.POST("close", hkStoreHandle.Close)
	}

	hkAreaGroup := group.Group("hkmo-area")
	{
		hkAreaHandle := dashboard.NewHkAreaHandle()
		hkAreaGroup.POST("list", hkAreaHandle.List)
	}

	// dashboard计划任务
	planGroup := group.Group("plan")
	{
		planHandle := dashboard.NewPlanHandle()
		// 门店数据
		planGroup.POST("store-data", planHandle.StoreData)
		// 任务数据
		planGroup.POST("plan-data", planHandle.PlanData)
		// 任务列表
		planGroup.POST("plan-list", planHandle.PlanList)
		// 任务详情
		// planGroup.POST("plan-detail", planHandle.PlanDetail)
		// 播种列表
		planGroup.POST("sow-task-list", planHandle.SowTaskList)

		// planGroup.POST("sow-task-detail", planHandle.SowTaskDetail)
	}

	reserveGroup := group.Group("reserve")
	{
		reserveHandle := dashboard.NewReserveHandle()
		// 数据看板
		reserveGroup.POST("data-kanban", reserveHandle.DataKanban)
		// 数据列表
		reserveGroup.POST("data-list", reserveHandle.DataList)
	}
}
