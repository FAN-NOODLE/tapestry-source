package http

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger/endpoint"
)

func Endpoints() []*swagger.Endpoint {
	tasks := newTasksHandler()

	return []*swagger.Endpoint{
		endpoint.Post("/transfer_customer/task",
			"获取任务详情",
			endpoint.Handler(tasks.get)),

		endpoint.Post("/transfer_customer/tasks_list",
			"获取任务列表",
			endpoint.Handler(tasks.getList)),

		endpoint.Post("/transfer_customer/tasks_add",
			"新增任务",
			endpoint.Handler(tasks.add)),

		endpoint.Post("/transfer_customer/tasks_update",
			"更新任务",
			endpoint.Handler(tasks.update)),

		endpoint.Post("/transfer_customer/tasks_delete",
			"删除任务",
			endpoint.Handler(tasks.delete)),

		endpoint.Post("/transfer_customer/tasks_report",
			"任务报表",
			endpoint.Handler(tasks.report)),
	}
}