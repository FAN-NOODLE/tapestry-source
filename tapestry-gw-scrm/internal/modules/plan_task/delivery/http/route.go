package http

import (
	"net/http"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger/endpoint"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	planTask "gw-scrm/internal/modules/plan_task"
)

func Endpoints() []*swagger.Endpoint {
	return []*swagger.Endpoint{
		endpoint.Post("/plan_tasks/info", "子任务详情",
			endpoint.Handler(planTaskInfo),
			endpoint.Tags("任务管理"),
			endpoint.Body(planTask.PlanTasksInfoRequst{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: planTask.PlanTasksResponse{},
			}, "planTaskInfo", "接口调用成功")),
		endpoint.Post("/plan_tasks/add", "新增子任务",
			endpoint.Handler(planTasksAdd),
			endpoint.Tags("任务管理"),
			endpoint.Body(planTask.PlanTasksRequst{}, "请求参数", true),
			endpoint.Response(http.StatusOK, "",
				"planTasksAdd", "接口调用成功"),
		),
		endpoint.Post("/plan_tasks/update", "修改子任务",
			endpoint.Handler(planTasksUpdate),
			endpoint.Tags("任务管理"),
			endpoint.Body(planTask.PlanTasksRequst{}, "请求参数", true),
			endpoint.Response(http.StatusOK, "",
				"planTasksUpdate", "接口调用成功"),
		),
		endpoint.Post("/plan_tasks/delete", "删除子任务",
			endpoint.Handler(planTasksDelete),
			endpoint.Tags("任务管理"),
			endpoint.Body(planTask.PlanTasksInfoRequst{}, "请求参数", true),
			endpoint.Response(http.StatusOK, "",
				"planTasksDelete", "接口调用成功"),
		),
		endpoint.Post("/plan/list", "计划列表",
			endpoint.Handler(getPlanList),
			endpoint.Tags("计划列表"),
			endpoint.Body(planTask.PlanListRequest{}, "请求参数", true),
			endpoint.Response(http.StatusOK, "",
				"getPlanList", "接口调用成功"),
		),
		endpoint.Post("/task/list", "任务列表",
			endpoint.Handler(getTaskList),
			endpoint.Tags("任务列表"),
			endpoint.Body(planTask.TaskListRequest{}, "请求参数", true),
			endpoint.Response(http.StatusOK, "",
				"getTaskList", "接口调用成功"),
		),

		endpoint.Post("/plan/add", "计划任务新增",
			endpoint.Handler(planAdd),
			endpoint.Tags("新增计划任务"),
			endpoint.Body(planTask.PlanRequest{}, "请求参数", true),
			endpoint.Response(http.StatusOK, "", "", "计划任务添加成功")),

		endpoint.Post("/plan/edit", "计划任务编辑",
			endpoint.Handler(planEdit),
			endpoint.Tags("计划任务编辑"),
			endpoint.Body(planTask.ReceivePlanModel{}, "请求参数", true),
			endpoint.Response(http.StatusOK, "", "", "计划任务编辑成功")),

		endpoint.Post("/plan/remove", "计划任务删除",
			endpoint.Handler(planDelete),
			endpoint.Tags("计划任务删除"),
			endpoint.Body(planTask.PlanId{}, "请求参数", true),
			endpoint.Response(http.StatusOK, "", "", "计划任务删除成功")),

		endpoint.Post("plan/copy", "计划任务复制",
			endpoint.Handler(planCopy),
			endpoint.Tags("计划任务复制"),
			endpoint.Body(planTask.ReceivePlanModel{}, "请求参数", true),
			endpoint.Response(http.StatusOK, "", "", "计划任务复制成功")),

		endpoint.Post("/plan/update/status", "计划状态更新",
			endpoint.Handler(planStatus),
			endpoint.Tags("计划任务状态更新"),
			endpoint.Body(planTask.PlanStatus{}, "请求参数", true),
			endpoint.Response(http.StatusOK, "", "", "计划任务状态更新")),

		endpoint.Post("/plan_tasks/conditions/completion", "任务条件完成情况",
			endpoint.Handler(planTasksConditionsCompletion),
			endpoint.Body(planTask.PlanTasksConditionsRequst{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: planTask.PlanTasksConditionsAllotLogRespone{},
			}, "planTasksConditionsCompletion", "接口调用成功")),
		endpoint.Post("/plan_tasks/conditions/user_completion", "导购任务条件完成情况列表",
			endpoint.Handler(planTasksConditionsUserCompletion),
			endpoint.Body(planTask.PlanTasksConditionsBaRequst{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: planTask.PlanTasksConditionsBaRespone{},
			}, "planTasksConditionsCompletion", "接口调用成功")),

		endpoint.Post("/plan/subtask/customer", "子任务用户画像统计数据",
			endpoint.Handler(customerData),
			endpoint.Body(planTask.CustomerDataReq{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: planTask.CustomerDataReq{},
			}, "planTasksConditionsCompletion", "接口调用成功")),

		endpoint.Post("/plan/subtask/customer/list_export", "子任务用户画像数据列表数据导出",
			endpoint.Handler(customerlistDataExport),
			endpoint.Body(planTask.ExportTasksRequst{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: planTask.ExportTasksRequst{},
			}, "planTasksConditionsCompletion", "接口调用成功")),

		endpoint.Post("/plan/subtask/customer/detail_export", "子任务用户画像数据明细数据导出",
			endpoint.Handler(customerDetailDataExport),
			endpoint.Body(planTask.ExportTasksRequst{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: planTask.ExportTasksRequst{},
			}, "planTasksConditionsCompletion", "接口调用成功")),

		endpoint.Post("/plan/subtask/group_task", "子任务社群数据统计数据",
			endpoint.Handler(groupChatData),
			endpoint.Body(planTask.CustomerDataReq{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: planTask.CustomerDataReq{},
			}, "planTasksConditionsCompletion", "接口调用成功")),

		endpoint.Post("/plan/subtask/group_task_list", "子任务社群任务列表数据导出",
			endpoint.Handler(groupTaskListData),
			endpoint.Body(planTask.TaskListRequst{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: planTask.TaskListRequst{},
			}, "planTasksConditionsCompletion", "接口调用成功")),

		endpoint.Post("/plan/subtask/customer_task_list", "子任务社群任务列表数据导出",
			endpoint.Handler(customerTaskList),
			endpoint.Body(planTask.TaskListRequst{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: planTask.TaskListRequst{},
			}, "planTasksConditionsCompletion", "接口调用成功")),

		endpoint.Post("/plan/subtask/group_task_export", "子任务社群任务列表数据导出",
			endpoint.Handler(groupChatTaskDataExport),
			endpoint.Body(planTask.CustomerDataReq{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: planTask.CustomerDataReq{},
			}, "planTasksConditionsCompletion", "接口调用成功")),

		// TODO
		endpoint.Post("/plan/subtask/upload-users", "子任务导入员工信息",
			endpoint.Handler(subTaskUploadUsers),
			endpoint.Response(http.StatusOK, "", "planTaskUploadUsers", "接口调用成功")),

		endpoint.Post("/plan/task/data/stat", "主任务统计",
			endpoint.Handler(taskDataStat),
			endpoint.Body(planTask.PlanStatRequest{}, "请求参数", true),
			endpoint.Response(http.StatusOK, "", "taskDataStat", "接口调用成功")),

		endpoint.Post("/plan/task/data/stat/detail", "任务统计明细",
			endpoint.Handler(taskDataStatDetail),
			endpoint.Body(planTask.StatDetailRequest{}, "请求参数", true),
			endpoint.Response(http.StatusOK, "", "taskDataStatDetail", "接口调用成功")),

		endpoint.Post("/task/data/stat/detail/export", "任务统计明细",
			endpoint.Handler(taskDetailStatExport),
			endpoint.Body(planTask.StatDetailRequest{}, "请求参数", true),
			endpoint.Response(http.StatusOK, "", "taskDataStatDetailExport", "接口调用成功")),
		endpoint.Post("/task/csvFile/upload", "导入推送客户模板",
			endpoint.Handler(taskCSVFileUpload),
			endpoint.Tags("导入推送客户模板"),
			endpoint.Response(http.StatusOK, "", "taskCSVFileUpload", "接口调用成功")),

		endpoint.Post("/task/consumer/export", "客户导入计划客户名单数据导出",
			endpoint.Handler(taskConsumerExport),
			endpoint.Tags("客户导入计划客户名单数据导出"),
			endpoint.Response(http.StatusOK, "", "upload", "接口调用成功")),
	}
}
