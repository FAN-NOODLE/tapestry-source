package internal

import (
	"micro-scrm/internal/app/http/handle/customer"

	"github.com/gin-gonic/gin"
)

func Customer(customerGroup *gin.RouterGroup) {
	// 标签路由组
	customerTagsGroup := customerGroup.Group("tags")
	{
		customerTagsHandle := customer.NewTagsHandle()
		// 操作顾客的印象标签
		customerTagsGroup.POST("impression/operate", customerTagsHandle.Operate)
		customerWorkTagsHandle := customer.NewWorkTagsHandle()
		// 操作顾客的企微标签
		customerTagsGroup.POST("work/operate", customerWorkTagsHandle.Operate)
	}

	// 离职继承
	customerInheritGroup := customerGroup.Group("inherit")
	{
		customerInheritHandle := customer.NewInheritHandle()
		// 待继承导购统计列表
		customerInheritGroup.POST("/list", customerInheritHandle.List)
		// 待分配的群列表
		customerInheritGroup.POST("/allot/group-chat-list", customerInheritHandle.GetGroupChatAllotList)
		// 待分配的人员列表
		customerInheritGroup.POST("/allot/external-list", customerInheritHandle.GetExternalUserAllotList)
		// 已分配的导购列表
		customerInheritGroup.POST("/takeover/list", customerInheritHandle.GetTakeOverList)

		// 分配离职成员的客户
		customerInheritGroup.POST("transfer_customer", customerInheritHandle.TransferCustomer)
		// 分配离职成员的客户群
		customerInheritGroup.POST("transfer_customer_group", customerInheritHandle.TransferCustomerGroup)
		// 接替的外部联系人列表
		customerInheritGroup.POST("/takeover/external-list", customerInheritHandle.GetTakeOverExternalList)
		// 接替的群列表
		customerInheritGroup.POST("/takeover/group-list", customerInheritHandle.GetTakeOverGroupList)
		// 更新离职待继承数据（同步企微数据到mysql）
		customerInheritGroup.POST("/sync-customer-data", customerInheritHandle.SyncData)

	}

	// 修改备注
	customerHandle := customer.NewCustomerHandle()
	// 后台用户列表
	customerGroup.POST("user-list", customerHandle.UserList)
	// 后台用户详情
	customerGroup.POST("user-detail", customerHandle.UserDetail)
	// 修改用户的备注
	customerGroup.POST("update-remark", customerHandle.UpdateRemark)
	// 修改用户的描述
	customerGroup.POST("update-description", customerHandle.UpdateDescription)
	// 获取外部联系人的基本信息
	customerGroup.POST("user-simple-info", customerHandle.UserSimpleInfo)
	// 同步unionid的微信openid
	customerGroup.POST("set-info", customerHandle.SetInfo)
	// 获取外部联系人的导购列表
	customerGroup.POST("member-lists", customerHandle.MemberLists)

	// 修改顾客地址、鞋码、生日
	customerGroup.POST("save-info-add", customer.NewInfoAddHandle().SaveInfoAdd)
	// 获取顾客地址、鞋码、生日
	customerGroup.POST("get-info-add", customer.NewInfoAddHandle().GetInfoAdd)

	// 顾客端--我的
	customerGroup.POST("my-center", customerHandle.MyCenter)
	// 顾客端--我的-添加专属顾问
	customerGroup.POST("my-sa", customerHandle.MySa)
	// 导购或店长获取顾客邀约和播种数量
	customerGroup.POST("user-detail-tips-num", customerHandle.UserDetailTipsNum)
	// 顾客和导购的绑定关系
	customerGroup.POST("user-bind-relation", customerHandle.UserBindRelation)

	// 获取外部联系人关系
	customerGroup.POST("get-user-member-relation", customer.NewRelationHandle().GetUserMemberRelation)

	// 外部联系人
	externalUserGroup := customerGroup.Group("external-user")
	{
		externalUserHandle := customer.NewExternalUserHandle()
		// 筛选外部联系人列表的条件
		externalUserGroup.POST("choice-option", externalUserHandle.ChoiceOption)
		// 筛选外部联系人列表
		externalUserGroup.POST("list", externalUserHandle.List)
		externalUserGroup.POST("list-ck", externalUserHandle.ListCk)
		// 外部联系人详情
		externalUserGroup.POST("detail", externalUserHandle.Detail)
		// 导购或管理的导购userid
		externalUserGroup.POST("get-member-userid", externalUserHandle.GetMemberUserIds)
		// 三方接口的会员卡号查询信息
		externalUserGroup.POST("customer-list", externalUserHandle.CustomerList)
		// 客户列表全选的统计数量
		externalUserGroup.POST("select-all-num", externalUserHandle.SelectAllNum)
		// 客户列表全选返回external_user_id列表
		externalUserGroup.POST("select-all-list", externalUserHandle.SelectAllList)

	}
	// 顾客列表（大宽表版）202411月
	externalSearchGroup := customerGroup.Group("external-search")
	{
		externalSearchHandle := customer.NewExternalSearchHandle()
		// 筛选外部联系人列表的条件
		externalSearchGroup.POST("choice-option", externalSearchHandle.SearchOption)
		// 顾客列表或者顾客关系列表
		externalSearchGroup.POST("list", externalSearchHandle.List)
		// 获取全选所有id
		externalSearchGroup.POST("get-open-ids", externalSearchHandle.GetOpenIds)
		// 顾客详情
		externalSearchGroup.POST("detail", externalSearchHandle.Detail)

	}
	// sw顾客列表新版 20250915
	externalSearchSwGroup := customerGroup.Group("external-search-sw")
	{
		externalSearchSwHandle := customer.NewExternalSearchSwHandle()
		// 筛选外部联系人列表的条件
		externalSearchSwGroup.POST("choice-option", externalSearchSwHandle.SearchOption)
		// 顾客列表或者顾客关系列表
		externalSearchSwGroup.POST("list", externalSearchSwHandle.List)
		// 获取全选所有id
		externalSearchSwGroup.POST("get-open-ids", externalSearchSwHandle.GetOpenIds)
	}

	externalCrowdSwGroup := customerGroup.Group("external-crowd-sw")
	{
		externalCrowdSwHandle := customer.NewExternalCrowdSwHandle()
		// 筛选外部联系人列表的条件
		externalCrowdSwGroup.POST("add", externalCrowdSwHandle.Add)
		// 获取外部联系人的内容
		externalCrowdSwGroup.POST("info", externalCrowdSwHandle.Info)
		// 获取白名单openid集合
		externalCrowdSwGroup.POST("get-open-ids", externalCrowdSwHandle.GetOpenIds)
		// 获取导购的人群包列表
		externalCrowdSwGroup.POST("list", externalCrowdSwHandle.List)
		// 编辑修改人去包名称
		externalCrowdSwGroup.POST("edit", externalCrowdSwHandle.Edit)
		// 人群包列表排序调整
		externalCrowdSwGroup.POST("update-sort", externalCrowdSwHandle.UpdateSort)
		// 人群包删除
		externalCrowdSwGroup.POST("del", externalCrowdSwHandle.Del)
		// 人群包联系状态更新
		externalCrowdSwGroup.POST("update-contact", externalCrowdSwHandle.UpdateContact)
		// 人群包联系状态批量更新
		externalCrowdSwGroup.POST("batch-update-contact", externalCrowdSwHandle.BatchUpdateContact)

	}
	externalCrowdGroup := customerGroup.Group("external-crowd")
	{
		externalCrowdHandle := customer.NewExternalCrowdHandle()
		// 筛选外部联系人列表的条件
		externalCrowdGroup.POST("add", externalCrowdHandle.Add)
		// 获取外部联系人的内容
		externalCrowdGroup.POST("info", externalCrowdHandle.Info)
		// 获取白名单openid集合
		externalCrowdGroup.POST("get-open-ids", externalCrowdHandle.GetOpenIds)
		// 获取导购的人群包列表
		externalCrowdGroup.POST("list", externalCrowdHandle.List)
		// 编辑修改人去包名称
		externalCrowdGroup.POST("edit", externalCrowdHandle.Edit)
		// 人群包列表排序调整
		externalCrowdGroup.POST("update-sort", externalCrowdHandle.UpdateSort)
		// 人群包删除
		externalCrowdGroup.POST("del", externalCrowdHandle.Del)
		// 人群包联系状态更新
		externalCrowdGroup.POST("update-contact", externalCrowdHandle.UpdateContact)
	}

	// 播种任务
	sowTakGroup := customerGroup.Group("sow-task")
	{
		sowTakHandle := customer.NewSowTaskHandle()
		sowTakGroup.POST("list", sowTakHandle.List)
		sowTakGroup.POST("create", sowTakHandle.Create)
		sowTakGroup.POST("update", sowTakHandle.Update)
		sowTakGroup.POST("finished", sowTakHandle.Finished)
		sowTakGroup.POST("delete", sowTakHandle.Delete)
		sowTakGroup.POST("template-list", sowTakHandle.TemplateList)
		sowTakGroup.POST("new-create", sowTakHandle.NewCreate)
	}
	// 播种任务模版
	sowTemplateGroup := customerGroup.Group("sow-template")
	{
		sowTemplateHandle := customer.NewSowTemplateHandle()
		sowTemplateGroup.POST("list", sowTemplateHandle.List)
		sowTemplateGroup.POST("save", sowTemplateHandle.Save)
		sowTemplateGroup.POST("detail", sowTemplateHandle.Detail)
		sowTemplateGroup.POST("update-status", sowTemplateHandle.UpdateStatus)
		sowTemplateGroup.POST("delete", sowTemplateHandle.Delete)
		sowTemplateGroup.POST("all-task-list", sowTemplateHandle.AllTaskList)
	}
	// 小程序数据罗盘的统计
	dataCompassGroup := customerGroup.Group("data-compass")
	{
		dataCompassHandle := customer.NewDataCompassHandle()
		// 社群数据统计
		dataCompassGroup.POST("chat-group-summary", dataCompassHandle.ChatGroupSummary)
		// 好友数据统计
		dataCompassGroup.POST("friends-summary", dataCompassHandle.FriendsSummary)
	}

	customerOnInheritGroup := customerGroup.Group("on-inherit")
	{
		customerOnInheritHandle := customer.NewOnInheritHandle()
		// 待继承导购统计列表
		customerOnInheritGroup.POST("/takeover-member-list", customerOnInheritHandle.GetTakeoverMemberList)
		// 获取全选的人员总数
		customerOnInheritGroup.POST("/select-all-num", customerOnInheritHandle.GetSelectAllNum)
		// 全选接替分配的关联关系数据
		customerOnInheritGroup.POST("/select-all-list", customerOnInheritHandle.GetSelectList)

		// 获取全选的可转移人员总数(新版)
		customerOnInheritGroup.POST("/transfer-all-num", customerOnInheritHandle.TransferCount)
		// 获取全选的可转移人员列表(新版)
		customerOnInheritGroup.POST("/transfer-all-list", customerOnInheritHandle.TransferList)

	}
}
