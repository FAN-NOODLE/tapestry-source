package internal

import (
	"github.com/gin-gonic/gin"
	"micro-scrm/internal/app/http/handle/member"
)

func Member(memberGroup *gin.RouterGroup) {

	memberHandle := member.NewMembersHandle()
	{
		// 导购部门详情
		memberGroup.POST("/get-detail", memberHandle.GetDetail)
		// 导购信息
		memberGroup.POST("/info", memberHandle.Info)
		// 获取导购的门店列表
		memberGroup.POST("/store-lists", memberHandle.StoreLists)
		// 获取员工列表
		memberGroup.POST("/staff-lists", memberHandle.StaffLists)
		// 根据部门id获取员工列表
		memberGroup.POST("/list-by-third-ids", memberHandle.MemberListByThirdIds)
		// 数据上报
		memberGroup.POST("/report-data", memberHandle.ReportData)
		// 下载中心白名单用户
		memberGroup.POST("/download-white-list", memberHandle.DownloadWhiteList)
		// 第一次进入某页面给的前端引导提示
		memberGroup.POST("/first-open", memberHandle.FirstOpen)
		// 根据门店编码获取员工信息
		memberGroup.POST("/get-member-list-by-store-codes", memberHandle.GetMemberListByStoreCodes)
	}

	PerformanceListHandle := member.NewPerformanceListHandle()
	{
		// 导购活跃榜顾客榜单(coach 和 ks共用)
		memberGroup.POST("performance-list/active-customers", PerformanceListHandle.ActiveCustomers)
		// 邀约排行(coach)
		memberGroup.POST("performance-list/invitation", PerformanceListHandle.Invitation)
		// 明星榜-活跃榜(coach)
		memberGroup.POST("performance-list/star-active-customers", PerformanceListHandle.StarActiveCustomers)
		// 明星榜-邀约排行(coach)
		memberGroup.POST("performance-list/star-invitation", PerformanceListHandle.StarInvitation)

		// sw顾客榜单
		memberGroup.POST("performance-list/sw-customers", PerformanceListHandle.SwCustomers)
		// sw活跃榜单
		memberGroup.POST("performance-list/sw-active", PerformanceListHandle.SwActive)
		// 业绩榜单(ks 和 sw)
		memberGroup.POST("performance-list/achievement", PerformanceListHandle.Achievement)
	}

}
