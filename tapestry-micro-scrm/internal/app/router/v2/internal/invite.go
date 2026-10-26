package internal

import (
	"github.com/gin-gonic/gin"
	"micro-scrm/internal/app/http/handle/invite"
)

// Invite 邀请函
func Invite(inviteGroup *gin.RouterGroup) {

	templateHandle := invite.NewTemplateHandle()
	{
		// 后台管理
		// 创建或保存邀请函模板
		inviteGroup.POST("/template/save", templateHandle.Save)
		// 邀请函模板列表
		inviteGroup.POST("/template/list", templateHandle.List)
		// 邀请函模板详情
		inviteGroup.POST("/template/info", templateHandle.Info)
		// 删除邀请函模板
		inviteGroup.POST("/template/delete", templateHandle.Delete)
		// 邀请函-操作上下架
		inviteGroup.POST("/template/update-status", templateHandle.UpdateStatus)

		// 小程序
		// 导购邀请函模板列表
		inviteGroup.POST("/template/config-list", templateHandle.ConfigList)
		// 导购邀请函模板详情
		inviteGroup.POST("/template/config-detail", templateHandle.ConfigDetail)

	}

	orderHandle := invite.NewOrderHandle()
	{
		// 导购发送邀约单记录
		inviteGroup.POST("/create-record", orderHandle.CreateRecord)

		// 顾客查看邀请函信息
		inviteGroup.POST("/get-record-info", orderHandle.GetRecordInfo)
		// 顾客确认邀请函订单
		inviteGroup.POST("/create-order", orderHandle.CreateInviteOrder)

		// 后台管理 邀请数据列表
		inviteGroup.POST("/record-order-list", orderHandle.RecordOrderList)
		// 导购取消分享小程序卡片，删除邀请函记录
		inviteGroup.POST("/record-delete", orderHandle.RecordDelete)

	}

}
