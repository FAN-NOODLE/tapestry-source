package internal

import (
	"micro-scrm/internal/app/http/handle/display"

	"github.com/gin-gonic/gin"
)

// Display VM陈列任务路由
func Display(displayGroupRouter *gin.RouterGroup) {

	taskHandle := display.NewTaskHandle()
	displayTask := displayGroupRouter.Group("task")
	{
		// 新增VM陈列任务
		displayTask.POST("add", taskHandle.AddTask)
		// VM陈列任务列表
		displayTask.POST("list", taskHandle.List)
		// VM陈列任务查看
		displayTask.POST("detail", taskHandle.Detail)
		// VM陈列任务复制
		displayTask.POST("copy", taskHandle.Copy)
		// VM陈列任务更新上下架
		displayTask.POST("status", taskHandle.UpdateShow)
		// VM陈列任务编辑
		displayTask.POST("update", taskHandle.UpdateTask)
		// VM陈列数据
		displayTask.POST("allot-list", taskHandle.AllotList)
		// VM陈列数据的统计数
		displayTask.POST("allot-num", taskHandle.AllotNum)
		// VM陈列数据导出
		displayTask.POST("allot-list-export", taskHandle.AllotListExport)
		// VM陈列消息提醒配置
		displayTask.POST("notice-setting", taskHandle.NoticeSetting)
		// VM陈列消息的配置内容
		displayTask.POST("notice-setting-info", taskHandle.NoticeSettingInfo)
		// VM任务门店的预览，详情
		displayTask.POST("allot-info", taskHandle.AllotInfo)
		// VM任务门店管理后台审核
		displayTask.POST("allot-check", taskHandle.AllotCheck)
		// VM任务陈列数据发送模板消息 一键提醒
		displayTask.POST("send-notice", taskHandle.AllotSendNotice)
		// VM陈列任务数据一键提醒门店数
		displayTask.POST("notice-num", taskHandle.AllotNoticeNum)
		// VM imageList列表
		displayTask.POST("image-list", taskHandle.ImageList)
		// 获取VM成员列表
		displayTask.POST("vm-list", taskHandle.VmList)
		// 导购头像跨域，上传cos
		displayTask.POST("cos-avatar", taskHandle.CosAvatar)
		// 上传分配名单校验
		displayTask.POST("review-file-check", taskHandle.ReviewFileCheck)
		// 下载分配名单数据
		displayTask.POST("review-export", taskHandle.ReviewExport)
	}

	displayApi := displayGroupRouter.Group("api")
	{
		// 获取任务详情
		displayApi.POST("task-info", taskHandle.TaskInfo)
		// 获取门店任务列表
		displayApi.POST("task-list", taskHandle.TaskList)
		// 保存草稿/提交
		displayApi.POST("save", taskHandle.Save)
		// 导购查看图片VM的备注上报
		displayApi.POST("look", taskHandle.Look)
		// 获取VM的权限信息
		displayApi.POST("auth", taskHandle.Auth)
		// 获取审核成功的门店列表
		displayApi.POST("store-list", taskHandle.StoreList)
		// H5端列表获取任务状态下的数量
		displayApi.POST("list-tips", taskHandle.ListTips)
		// VM的审核列表
		displayApi.POST("check-list", taskHandle.CheckList)
		// VM任务详情中的门店列表
		displayApi.POST("check-store-list", taskHandle.CheckStoreList)
		// VM任务门店的预览，审核，详情
		displayApi.POST("store-task-info", taskHandle.StoreTaskInfo)
		// VM端，审核门店的任务
		displayApi.POST("vm-check", taskHandle.VmCheck)
		// VM端自动保存图片
		displayApi.POST("save-image", taskHandle.SaveImage)
		// VM端自动删除图片
		displayApi.POST("del-image", taskHandle.DelImage)
		// VM端自动保存备注
		displayApi.POST("save-comment", taskHandle.SaveComment)
		// VM端自动删除备注
		displayApi.POST("del-comment", taskHandle.DelComment)
		// VM端点击提醒
		displayApi.POST("send-notice", taskHandle.SendNotice)
		// VM看板获取所有已上架的任务列表
		displayApi.POST("all-list", taskHandle.AllList)
		// VM看板数据的统计dashboard
		displayApi.POST("vm-dashboard", taskHandle.DashBoard)
		// VM看板获取管理范围的门店最新任务未通过的数据
		displayApi.POST("fail-list", taskHandle.FailList)
		// 下钻页面数据统计
		displayApi.POST("dashboard-list", taskHandle.DashBoardList)
		// 获取门店下的任务列表
		displayApi.POST("store-task-list", taskHandle.StoreTaskList)
		// special 保存草稿/提交
		displayApi.POST("special-save", taskHandle.SpecialSave)
	}
}
