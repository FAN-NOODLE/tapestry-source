package internal

import (
	"github.com/gin-gonic/gin"
	"micro-scrm/internal/app/http/handle/reserve"
)

func Reserve(groupRouter *gin.RouterGroup) {
	// ---------------------到店预约---------------------
	reservesHandle := reserve.NewOrderHandle()
	{
		// 导购的预约单创建
		groupRouter.POST("/order/create", reservesHandle.CreateOrder)
		// 导购的预约单修改
		groupRouter.POST("/order/edit", reservesHandle.EditOrder)
		// 导购(店长)取消预约
		groupRouter.POST("/order/cancel", reservesHandle.CancelOrder)
		// 预约单列表
		groupRouter.POST("/order/lists", reservesHandle.OrderList)
		// 预约单列表(coach)
		groupRouter.POST("/order/lists-by-coach", reservesHandle.OrderListByCoach)
		// 预约单列表详情
		groupRouter.POST("/order/detail", reservesHandle.OrderDetail)
		// 微信端用户获取预约单信息
		groupRouter.POST("/order/user-order-info", reservesHandle.UserOrderInfo)
		// 微信端用户取消预约
		groupRouter.POST("/order/user-cancel-order", reservesHandle.UserCancelOrder)
		// (KS/SW)微信端用户编辑预约单信息
		groupRouter.POST("/order/user-edit-order", reservesHandle.UserEditOrder)
		// (KS/SW)微信端用户给服务打分
		groupRouter.POST("/order/user-make-score", reservesHandle.UserMakeScore)
		// (KS/SW)微信端用户的评分信息
		groupRouter.POST("/order/user-score-info", reservesHandle.ScoreInfo)
		// 获取预约时间段
		groupRouter.POST("/order/time-list", reservesHandle.GetReserveTimeList)
		// 预约分配
		groupRouter.POST("/order/assign", reservesHandle.OrderAssign)
		// 顾客受邀单列表
		groupRouter.POST("/order/reserve-order-list", reservesHandle.ReserveOrderList)
		// 导购的邀约排行
		groupRouter.POST("/order/rank-list", reservesHandle.RankList)

		// 分享预约单
		groupRouter.POST("/order/share", reservesHandle.Share)
		// 导购端获取预约单详情
		groupRouter.POST("/order/order-info", reservesHandle.OrderInfo)
		// 完成预约
		groupRouter.POST("/order/finish", reservesHandle.Finish)
		// SW端增加签到按钮
		groupRouter.POST("/order/user-sign-in", reservesHandle.UserSignIn)
		// coach 导购确认订单是否是到店购买
		groupRouter.POST("/order/chang-status", reservesHandle.ChangeStatus)
		// coach 导购确认用户到店
		groupRouter.POST("/order/confirm-store", reservesHandle.ConfirmStore)
		// coach 导购邀约用户评论预约单
		groupRouter.POST("/order/comment-share", reservesHandle.CommentShare)
		// coach 用户编辑评论
		groupRouter.POST("/order/user-edit-score", reservesHandle.UserEditScore)
		// 验证导购是否可以给用户发送邀请评价
		groupRouter.POST("/order/verify-share", reservesHandle.VerifyShare)
		// 导购取消分享小程序卡片，删除创建的邀约单
		groupRouter.POST("/order/delete", reservesHandle.Delete)
		// 导出数据的请求列表数据
		groupRouter.POST("/order/export-lists", reservesHandle.OrderExportList)
		// 针对邀约时间段检验邀约单
		groupRouter.POST("/order/check-reserve-time", reservesHandle.CheckReserveTime)
		// 查看用户的所有邀约单列表数据
		groupRouter.POST("/order/customer-reserve-lists", reservesHandle.CustomerReserveLists)
		// 导购创建分享了邀约单，才算创建邀约
		groupRouter.POST("/order/show", reservesHandle.Show)
		// 顾客受邀单列表下载中心接口
		groupRouter.POST("/order/reserve-order-list-new", reservesHandle.ReserveOrderListNew)

		// 后台预约单列表查询，走CK查询
		groupRouter.POST("/order/lists-ck", reservesHandle.OrderListCk)
		// 管理后台受邀统计的订单明细导出
		groupRouter.POST("/order/export-order-record", reservesHandle.ExportOrderRecord)
	}

	reservesScoreHandle := reserve.NewServiceScoreHandle()
	{
		// 用户评分列表
		groupRouter.POST("/service-score/list", reservesScoreHandle.List)
	}

	reservesServiceHandle := reserve.NewReserveServiceHandle()
	{
		// 新增服务
		groupRouter.POST("/service/create", reservesServiceHandle.Create)
		// 修改
		groupRouter.POST("/service/update", reservesServiceHandle.Update)
		// 删除
		groupRouter.POST("/service/delete", reservesServiceHandle.Delete)
		// 列表
		groupRouter.POST("/service/list", reservesServiceHandle.List)
		// coach 检查是否可以删除改服务
		groupRouter.POST("/service/check-delete", reservesServiceHandle.CheckDelete)
		// 导购端获取服务列表
		groupRouter.POST("/service/member-list", reservesServiceHandle.ListsByMember)
		// 服务详情接口
		groupRouter.POST("/service/info", reservesServiceHandle.Info)
	}

}
