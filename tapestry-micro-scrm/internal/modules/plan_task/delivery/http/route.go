package http

import (
	"net/http"

	"micro-scrm/internal/modules/members"
	memberMysql "micro-scrm/internal/modules/members/respository/mysql"
	"micro-scrm/internal/modules/plan_task/delivery/http_entity"
	"micro-scrm/internal/modules/plan_task/respository/mysql"
	"micro-scrm/internal/modules/ws_trace"
	wsTraceMysql "micro-scrm/internal/modules/ws_trace/respository/mysql"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger/endpoint"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"micro-scrm/internal/modules/plan_task"
)

// Endpoints 路由定义
func Endpoints() []*swagger.Endpoint {
	wsTraceService = ws_trace.NewService(wsTraceMysql.NewWsTrace())
	memberService = members.New(memberMysql.New())
	pService = plan_task.NewService(mysql.NewPlanTask())
	return []*swagger.Endpoint{
		endpoint.Post("/plan/list", "计划列表",
			endpoint.Handler(GetPlanList),
			endpoint.Tags("计划列表"),
			endpoint.Body(plan_task.PlanListRequest{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: plan_task.PlanListResponse{},
			}, "getPlanList", "接口调用成功")),

		endpoint.Post("/task/list", "任务列表",
			endpoint.Handler(GetTaskList),
			endpoint.Tags("任务列表"),
			endpoint.Body(plan_task.TaskListRequest{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: plan_task.TaskListResponse{},
			}, "getTaskList", "接口调用成功")),
		endpoint.Post("/plan/front-list", "计划列表",
			endpoint.Handler(getPlanListByFront),
			endpoint.Tags("计划列表"),
			endpoint.Body(plan_task.PlanListRequest{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: plan_task.PlanListResponse{},
			}, "getPlanListByFront", "接口调用成功")),

		endpoint.Post("/plan/front-task-list", "计划子任务列表",
			endpoint.Handler(getPlanTaskListByFront),
			endpoint.Tags("计划子任务列表"),
			endpoint.Body(plan_task.PlanListRequest{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: plan_task.PlanListResponse{},
			}, "getPlanTaskListByFront", "接口调用成功")),

		endpoint.Post("/plan/front-task-info", "子任务详情",
			endpoint.Handler(getPlanTaskInfoByFront),
			endpoint.Tags("计划子任务列表"),
			endpoint.Body(plan_task.PlanListRequest{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: plan_task.PlanListResponse{},
			}, "getPlanTaskListByFront", "接口调用成功")),

		endpoint.Post("/plan/front-update-status", "子任务详情任务状态更新",
			endpoint.Handler(UpdateGroupStatus),
			endpoint.Tags("子任务详情任务状态更新"),
			endpoint.Body(plan_task.UpdateReceive{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: plan_task.PlanListResponse{},
			}, "UpdateGroupStatus", "接口调用成功")),
		endpoint.Post("/plan/front-set-status", "完成状态",
			endpoint.Handler(setPlanTaskStatusByFront),
			endpoint.Tags("计划子任务列表"),
			endpoint.Body(plan_task.PlanListRequest{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: plan_task.PlanListResponse{},
			}, "getPlanTaskListByFront", "接口调用成功")),
		endpoint.Post("/plan_task/add_plan", "新增计划任务",
			endpoint.Handler(AddPlan),
			endpoint.Tags("新增计划任务"),
			endpoint.Response(http.StatusOK, "", "", "计划任务添加成功")),

		endpoint.Post("/plan_task/remove", "计划任务删除",
			endpoint.Handler(RemovePlan),
			endpoint.Tags("计划任务删除"),
			endpoint.Response(http.StatusOK, "", "", "计划任务删除成功")),

		endpoint.Post("/plan_task/plan/edit", "计划任务编辑",
			endpoint.Handler(PlanEdit),
			endpoint.Tags("计划任务编辑"),
			endpoint.Response(http.StatusOK, "", "", "编辑成功")),

		endpoint.Post("/plan_task/plan/copy", "计划任务复制",
			endpoint.Handler(PlanCopy),
			endpoint.Tags("计划任务复制"),
			endpoint.Response(http.StatusOK, "", "", "任务复制成功")),

		endpoint.Post("/plan_task/plan/update_status", "更新状态",
			endpoint.Handler(UpdateStatus),
			endpoint.Tags("更新状态"),
			endpoint.Response(http.StatusOK, "", "", "更新状态")),

		endpoint.Post("/plan_tasks/add", "新增子任务",
			endpoint.Handler(planTasksAdd),
			endpoint.Tags("任务管理"),
			endpoint.Body(plan_task.PlanTasksRequest{}, "请求参数", true),
			endpoint.Response(http.StatusOK, "",
				"planTasksAdd", "接口调用成功"),
		),
		endpoint.Post("/plan_tasks/update", "修改子任务",
			endpoint.Handler(planTasksUpdate),
			endpoint.Tags("任务管理"),
			endpoint.Body(plan_task.PlanTasksRequest{}, "请求参数", true),
			endpoint.Response(http.StatusOK, "",
				"planTasksUpdate", "接口调用成功"),
		),
		endpoint.Post("/plan_tasks/info", "子任务详情",
			endpoint.Handler(planTasksInfo),
			endpoint.Tags("任务管理"),
			endpoint.Body(plan_task.PlanTasksInfoRequst{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: plan_task.PlanTasksResponse{},
			}, "planTasksInfo", "接口调用成功")),
		endpoint.Post("/plan_tasks/delete", "删除任务",
			endpoint.Handler(planTasksDelet),
			endpoint.Tags("任务管理"),
			endpoint.Body(plan_task.PlanTasksInfoRequst{}, "请求参数", true),
			endpoint.Response(http.StatusOK, "",
				"planTasksDelet", "接口调用成功"),
		),
		endpoint.Post("/plan_tasks/conditions/completion", "任务条件完成情况",
			endpoint.Handler(planTasksConditionsCompletion),
			endpoint.Tags("任务管理"),
			endpoint.Body(plan_task.PlanTasksConditionsRequst{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: plan_task.PlanTasksConditionsAllotLogRespone{},
			}, "planTasksConditionsCompletion", "接口调用成功")),
		endpoint.Post("/plan_tasks/conditions/user_completion", "导购任务条件完成情况列表",
			endpoint.Handler(planTasksConditionsUserCompletion),
			endpoint.Tags("任务管理"),
			endpoint.Body(plan_task.PlanTasksConditionsBaRequst{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: plan_task.PlanTasksConditionsBaRespone{},
			}, "planTasksInfo", "接口调用成功")),

		endpoint.Post("/plan/csvFile/upload", "上传标签的csv文件",
			endpoint.Handler(getTagCSVFile),
			endpoint.Tags("上传标签的csv文件"),
			endpoint.Body(plan_task.GetFileInfo{}, "请求参数", true),
			endpoint.Response(http.StatusOK, "", "", "接口调用成功")),

		endpoint.Post("/plan/upload/saveInfo", "保存上传信息",
			endpoint.Handler(DownloadFileInfo),
			endpoint.Tags("保存上传信息"),
			endpoint.Body(plan_task.DownLoadFileInfo{}, "请求参数", true),
			endpoint.Response(http.StatusOK, "", "", "接口调用成功")),

		endpoint.Post("/plan/upload/fileDetection", "验证文件是否存在",
			endpoint.Handler(GetAllFileName),
			endpoint.Tags("验证文件是否存在"),
			endpoint.Body(plan_task.FileDetection{}, "请求参数", true),
			endpoint.Response(http.StatusOK, "", "", "接口调用成功")),

		endpoint.Post("/plan/upload/update/is_extract", "更新完成状态",
			endpoint.Handler(UpdateIsExtract),
			endpoint.Tags("更新完成状态"),
			endpoint.Body(plan_task.UpdateIdInfo{}, "请求参数", true),
			endpoint.Response(http.StatusOK, "", "", "接口调用成功")),

		endpoint.Post("/plan/subtask/customer", "子任务用户画像统计数据",
			endpoint.Handler(customerData),
			endpoint.Tags("子任务用户画像统计数据"),
			endpoint.Body(plan_task.CustomerDataReq{}, "请求参数", true),
			endpoint.Response(http.StatusOK, "", "", "接口调用成功")),

		endpoint.Post("/plan/subtask/customer/list_export", "子任务用户画像数据列表数据导出",
			endpoint.Handler(taskListExport),
			endpoint.Tags("子任务用户画像数据任务数据导出"),
			endpoint.Body(http_entity.ExportTasksRequst{}, "请求参数", true),
			endpoint.Response(http.StatusOK, "", "", "接口调用成功")),
		endpoint.Post("/plan/subtask/group_task", "子任务社群任务统计数据",
			endpoint.Handler(groupTaskData),
			endpoint.Tags("子任务社群任务统计数据"),
			endpoint.Body(plan_task.CustomerDataReq{}, "请求参数", true),
			endpoint.Response(http.StatusOK, "", "", "接口调用成功")),

		endpoint.Post("/plan/subtask/group_task_export", "子任务社群任务数据下载",
			endpoint.Handler(groupTaskDataExport),
			endpoint.Tags("子任务社群任务统计数据"),
			endpoint.Body(plan_task.CustomerDataReq{}, "请求参数", true),
			endpoint.Response(http.StatusOK, "", "", "接口调用成功")),

		endpoint.Post("/plan/subtask/group_task_list", "子任务社群任务列表",
			endpoint.Handler(groupTaskListData),
			endpoint.Tags("子任务社群任务统计数据"),
			endpoint.Body(plan_task.CustomerDataReq{}, "请求参数", true),
			endpoint.Response(http.StatusOK, "", "", "接口调用成功")),

		endpoint.Post("/plan/subtask/customer_task_list", "子任务用户画像任务列表",
			endpoint.Handler(customerTaskListData),
			endpoint.Tags("子任务社群任务统计数据"),
			endpoint.Body(plan_task.CustomerDataReq{}, "请求参数", true),
			endpoint.Response(http.StatusOK, "", "", "接口调用成功")),

		endpoint.Post("/plan/subtask/upload-users", "上传用户信息",
			endpoint.Handler(uploadTaskUsers),
			endpoint.Tags("上传用户信息"),
			endpoint.Response(http.StatusOK, "", "", "接口调用成功")),
	}
}
