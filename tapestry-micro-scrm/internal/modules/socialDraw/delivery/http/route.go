package http

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger/endpoint"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"micro-scrm/internal/modules/chat"
	chatRepository "micro-scrm/internal/modules/chat/repository/mysql"
	"micro-scrm/internal/modules/customer"
	userEs "micro-scrm/internal/modules/customer/respository/elastic"
	customerMysql "micro-scrm/internal/modules/customer/respository/mysql"
	"micro-scrm/internal/modules/department"
	departmentSQL "micro-scrm/internal/modules/department/respository/mysql"
	"micro-scrm/internal/modules/drainage"
	drainageSQL "micro-scrm/internal/modules/drainage/respository/mysql"
	"micro-scrm/internal/modules/members"
	memberSQL "micro-scrm/internal/modules/members/respository/mysql"
	"micro-scrm/internal/modules/socialDraw"
	"micro-scrm/internal/modules/socialDraw/respository/ck"
	"micro-scrm/internal/modules/socialDraw/respository/mysql"
)

func Endpoints() []*swagger.Endpoint {
	groupChatService = *socialDraw.NewGroupChat(mysql.NewRepo(), ck.NewRepo())
	chatService = *chat.NewGroupChat(chatRepository.NewWorkGroupChat())
	drainageService = *drainage.NewDrainage(drainageSQL.NewDrainage())
	memberService = *members.New(memberSQL.New())
	departmentService = *department.New(departmentSQL.New())
	customerService = *customer.New(customerMysql.NewCustomer(), customerMysql.NewWorkUser(), userEs.NewEsRepo())
	groupChatMemberService = *chat.NewGroupChatMember(chatRepository.NewWorkGroupChatMember())
	return []*swagger.Endpoint{
		endpoint.Post("/chat-group/list", "社群列表",
			endpoint.Handler(workGroupChatList),
			endpoint.Body(socialDraw.WorkGroupListRequest{}, "请求参数", true),
			endpoint.Response(200, ginplus.APIResponse{
				Data: []*chat.WorkGroupChat{},
			}, "chatList", "接口调用成功"),
		),
		endpoint.Post("/chat-group/info", "社群详情",
			endpoint.Handler(workGroupChatInfo),
			endpoint.Body(socialDraw.WorkGroupChat{}, "请求参数", true),
			endpoint.Response(200, ginplus.APIResponse{
				Data: chat.WorkGroupChat{},
			}, "chatInfo", "接口调用成功"),
		),
		endpoint.Post("/chat-group/setting", "chat group setting",
			endpoint.Handler(workGroupChatSetting),
			endpoint.Body(socialDraw.WorkGroupChatSettingRequest{}, "请求参数", true),
			endpoint.Response(200, ginplus.APIResponse{
				Data: socialDraw.WorkGroupChatType{},
			}, "chatSetting", "接口调用成功"),
		),
		endpoint.Post("/chat-group/members", "社群成员明细-后台",
			endpoint.Handler(getWorkGroupChatMembers),
			endpoint.Body(socialDraw.WorkGroupChatMembersRequest{}, "请求参数", true),
			endpoint.Response(200, ginplus.APIResponse{
				Data: chat.WorkGroupChatMember{},
			}, "chatMembers", "接口调用成功"),
		),
		endpoint.Post("/chat-group/members_info", "社群成员明细",
			endpoint.Handler(getGroupChatMembersInfo),
			endpoint.Body(socialDraw.WorkGroupChatMembersRequest{}, "请求参数", true),
			endpoint.Response(200, ginplus.APIResponse{
				Data: chat.WorkGroupChatMember{},
			}, "chatMembers", "接口调用成功"),
		),
		endpoint.Post("/chat-group/members-guide", "社群会员成员列表",
			endpoint.Handler(getMemberList),
			endpoint.Body(socialDraw.WorkGroupChatMembersRequest{}, "请求参数", true),
			endpoint.Response(200, ginplus.APIResponse{
				Data: chat.WorkGroupChatMember{},
			}, "membersGuide", "接口调用成功"),
		),
		endpoint.Post("/chat-group/types", "社群所属类型",
			endpoint.Handler(workChatGroupTypeList),
			endpoint.Body(socialDraw.WorkGroupChatType{}, "请求参数", true),
			endpoint.Response(200, ginplus.APIResponse{
				Data: []*socialDraw.WorkGroupChatType{},
			}, "chatTypeList", "接口调用成功")),
		endpoint.Post("/chat-group/type-manage", "社群类型管理",
			endpoint.Handler(workChatGroupTypeManage),
			endpoint.Body(socialDraw.WorkGroupChatType{}, "请求参数", true),
			endpoint.Response(200, ginplus.APIResponse{
				Data: socialDraw.WorkGroupChatType{},
			}, "chatTypeManage", "接口调用成功")),
		endpoint.Post("/chat-group/index_list", "社群详情",
			endpoint.Handler(workH5GroupChatInfo),
			endpoint.Body(socialDraw.GroupInfoReq{}, "请求参数", true),
			endpoint.Response(200, ginplus.APIResponse{
				Data: socialDraw.WorkGroupChat{},
			}, "h5chatinfo", "接口调用成功"),
		),
		endpoint.Post("/chat-group/filters", "自定义下拉列表",
			endpoint.Handler(workGroupChatFilters),
			endpoint.Body(socialDraw.WorkGroupFilterRequest{}, "请求参数", true),
			endpoint.Response(200, ginplus.APIResponse{
				Data: socialDraw.WorkGroupFilterResponse{},
			}, "filterChannelList", "接口调用成功")),
		endpoint.Post("/chat-group/data-grid", "群数据统计指标",
			endpoint.Handler(getWorkGroupChatDataGrid),
			endpoint.Body(socialDraw.GroupInfoReq{}, "请求参数", true),
			endpoint.Response(200, ginplus.APIResponse{
				Data: socialDraw.WorkGroupChatDataGridResponse{},
			}, "filterChannelList", "接口调用成功")),
		endpoint.Post("/chat-group/upload-settings", "文件上传转数据", endpoint.Handler(fileUploadExchangeData)),
		endpoint.Post("/chat-group/batch-setting", "文件上传转数据",
			endpoint.Handler(workGroupChatBatchSetting),
			endpoint.Body(socialDraw.WorkGroupChatBatchSettingRequest{}, "请求参数", true)),
		endpoint.Post("/chat-group/data-update-time", "数据更新时间",
			endpoint.Handler(getLastDataUpdateTime),
			endpoint.Body(socialDraw.WorkGroupChatBatchSettingRequest{}, "请求参数", true)),
		endpoint.Post("/chat-group/sync", "同步群信息", endpoint.Handler(BatchSyncGroupChats)),
		endpoint.Post("/chat-group/activity-info", "获取用户活动详情",
			endpoint.Handler(getActivityInfo),
			endpoint.Body(socialDraw.ActivityReq{}, "请求参数", true),
			endpoint.Response(200, ginplus.APIResponse{
				Data: socialDraw.WorkGroupChatDataGridResponse{},
			}, "filterChannelList", "接口调用成功")),
	}
}
