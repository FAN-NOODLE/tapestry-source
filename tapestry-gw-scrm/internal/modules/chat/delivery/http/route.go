package chat_http

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger/endpoint"
)

func Endpoints() []*swagger.Endpoint {
	return []*swagger.Endpoint{
		endpoint.Post("/chat/msg-type-list", "消息类型列表", endpoint.Handler(msgTypeList)),
		endpoint.Post("/chat/agree-info", "获取会话存档同意信息", endpoint.Handler(agreeInfo)),
		endpoint.Post("/chat/msg-list", "消息分页列表", endpoint.Handler(msgList)),
		endpoint.Post("/chat/group-msg-list", "群消息分页列表", endpoint.Handler(groupMsgList)),
		endpoint.Post("/chat/msg-list-by-id", "根据消息ID找附近的消息", endpoint.Handler(msgListByID)),
		endpoint.Post("/chat/member-list", "员工列表", endpoint.Handler(memberList)),

		endpoint.Post("/chat/object/to-member-list", "员工与员工聊天分页列表", endpoint.Handler(objectToMemberList)),
		endpoint.Post("/chat/object/to-external-list", "员工与外部联系人聊天分页列表", endpoint.Handler(objectToExternalList)),

		endpoint.Post("/chat/object/to-group-list", "员工的社群列表", endpoint.Handler(objectToGroupList)),
		endpoint.Post("/chat/get-user-info", "获取用户信息", endpoint.Handler(getUserInfo)),
		endpoint.Post("/chat/get-file-content", "获取文件流", endpoint.Handler(getFileContent)),
	}
}
