package internal

import (
	"github.com/gin-gonic/gin"
	"micro-scrm/internal/app/http/handle/plan"
	"micro-scrm/internal/app/http/middleware/plan_task"
)

func Plan(planGroup *gin.RouterGroup) {
	planHandle := plan.NewPlansHandle()
	planGroup.POST("/import/fileDetection", planHandle.FileDetection)
	planGroup.POST("/import/saveInfo", planHandle.SaveInfo)
	planGroup.POST("/consumer/list", planHandle.ConsumerList)
	// 任务日历的预约单列表
	planGroup.POST("/month-lists", planHandle.MonthLists)
	// 任务中心的主任务列表
	planGroup.POST("/lists", planHandle.Lists)
	// 主任务详情
	planGroup.POST("/detail", planHandle.Detail)
	// 任务中心的未完成提示数
	planGroup.POST("/tips", planHandle.Tips)
	// 计划任务的退群顾客列表
	planGroup.POST("/leave-group-chat", planHandle.LeaveGroupChat)
	// 计划任务的sacode&&unionid|customerid的顾客列表
	planGroup.POST("/sa-customer-list", planHandle.SaCustomerLists)
	// coach升级的计划任务顾客列表
	planGroup.POST("/search-customer-list", planHandle.SearchCustomerLists)
	// coach升级的计划任务顾客列表
	planGroup.POST("/search-external-list", planHandle.SearchExternalLists)
	// coach升级的计划任务顾客列表,全选获取openid
	planGroup.POST("/search-external-all-openid", planHandle.SearchExternalAllOpenId)
	// coach顾客数据移动，群发移入单独沟通，单独沟通移出
	planGroup.POST("/customer-move", planHandle.CustomerMove)
	// 上报任务的发送情况
	planGroup.POST("/send-finish", planHandle.SendFinish)
	// 任务计划后台管理，子任务详情的数据明细导出接口
	planGroup.POST("/task-detail-export-base", planHandle.TaskDetailExportBase)
	planGroup.POST("/task-detail-export-list", planHandle.TaskDetailExportList)
	// 任务计划后台管理，子任务详情的数据明细导出接口--下载中心
	planGroup.POST("/task-detail-export-list-new", planHandle.TaskDetailExportListNew)
	// 计划任务点击转发给全部顾客的按钮上报
	planGroup.POST("/set-forward", planHandle.SetForward)
	// 主任务任务数据
	planGroup.POST("/task/data/new-stat", planHandle.NewStat)
	// 主任务的新数据明细
	planGroup.POST("/task/data/stat/new-detail", planHandle.NewStatDetail)
	// 已完成的主任务未触达提示关闭
	planGroup.POST("/close_touch_tips", planHandle.NewCloseTouchTips)
	// 已完成的子任务的未触达提示关闭
	planGroup.POST("/task/close_touch_tips", planHandle.NewCloseTaskTouchTips)
	// 点击添加好友的按钮上报
	planGroup.POST("/is-show", planHandle.IsShow)

	planGroup.POST("task/plan-name-list", planHandle.TaskPlanNameList)

	// 顾客的详情页的任务列表
	planGroup.POST("/customer-task-list", planHandle.CustomerTaskList)

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
	// 下面的接口要验证任务，有中间件验证 验证导购身份信息
	checkPlanGroup := planGroup.Use(plan_task.CheckConditionAuth())
	// 计划任务确认用户拨打电话
	checkPlanGroup.POST("/confirm-call", planHandle.ConfirmCall)
	// 计划任务确认联系好友
	checkPlanGroup.POST("/confirm-contact", planHandle.ConfirmContact)
	// 计划任务确认标记任务完成
	checkPlanGroup.POST("/confirm-finish", planHandle.ConfirmFinish)
	// 计划任务填写联系说明
	checkPlanGroup.POST("/set-remark", planHandle.SetRemark)
	// 点击添加好友的按钮上报
	planGroup.POST("/confirm-friend", planHandle.ConfirmFriend)

}
