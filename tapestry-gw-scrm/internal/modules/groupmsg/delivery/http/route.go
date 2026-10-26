package groupmsg_http

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger/endpoint"
)

func Endpoints() []*swagger.Endpoint {
	return []*swagger.Endpoint{
		endpoint.Post("/group-msg/list", "群发任务列表", endpoint.Handler(list)),
		endpoint.Post("/group-msg/group-send-add", "新增群发消息", endpoint.Handler(groupSendAdd)),
		endpoint.Post("/group-msg/group-send-edit", "编辑群发消息", endpoint.Handler(groupSendEdit)),
		endpoint.Post("/group-msg/member-list", "任务明细-执行群主下拉列表", endpoint.Handler(memberList)),
		endpoint.Post("/group-msg/task-info", "群发任务详情", endpoint.Handler(TaskInfo)),
		endpoint.Post("/group-msg/task-detail", "群发任务明细列表", endpoint.Handler(TaskDetail)),
		endpoint.Post("/group-msg/task-member-detail", "员工群发任务明细", endpoint.Handler(TaskMemberDetail)),
		endpoint.Post("/group-msg/task-member-send-list", "员工群发发送列表", endpoint.Handler(TaskMemberSendList)),
	}
}
