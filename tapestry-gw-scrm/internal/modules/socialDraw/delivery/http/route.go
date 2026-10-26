package http

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger/endpoint"
)

func Endpoints() []*swagger.Endpoint {
	return []*swagger.Endpoint{
		endpoint.Post("/chat-group/list", "获取群列表", endpoint.Handler(getWorkGroupChatList)),
		endpoint.Post("/chat-group/info", "获取群详情", endpoint.Handler(getWorkGroupChatInfo)),
		endpoint.Post("/chat-group/setting", "群设置", endpoint.Handler(setWorkGroupChat)),
		endpoint.Post("/chat-group/members", "群设置", endpoint.Handler(getWorkGroupChatMembers)),

		endpoint.Post("/chat-group/types", "社群类型列表", endpoint.Handler(getWorkGroupChatTypeList)),
		endpoint.Post("/chat-group/type-manage", "社群类型管理", endpoint.Handler(manageWorkGroupChatType)),
		endpoint.Post("/chat-group/data-grid", "社群详情数据统计", endpoint.Handler(getWorkGroupChatDataGrid)),
		endpoint.Post("/chat-group/upload-settings", "社群批量设置导入", endpoint.Handler(WorkGroupChatUploadSetting)),
		endpoint.Post("/chat-group/batch-setting", "社群批量设置", endpoint.Handler(WorkGroupChatBatchSetting)),

		endpoint.Post("/chat-group/filters", "获取群详情", endpoint.Handler(getWorkGroupChatFilters)),
		endpoint.Post("/chat-group/data-update-time", "数据更新时间", endpoint.Handler(WorkGroupChatDataUpdateTime)),
		endpoint.Post("/chat-group/sync", "同步群基本信息", endpoint.Handler(WorkGroupChatBatchSync)),
	}
}
