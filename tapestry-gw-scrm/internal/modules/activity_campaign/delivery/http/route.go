package http

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger/endpoint"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	planTask "gw-scrm/internal/modules/plan_task"
	"net/http"
)

func Endpoints() []*swagger.Endpoint {
	return []*swagger.Endpoint{
		endpoint.Post("/activity/sign-in/list", "签到列表",
			endpoint.Handler(activitySignInList),
			endpoint.Tags("签到列表")),
		endpoint.Post("/activity/sign-in/detail", "签到详情",
			endpoint.Handler(activitySignInDetail),
			endpoint.Tags("签到详情")),
		endpoint.Post("/activity/sign-in/export", "签到导出",
			endpoint.Handler(activitySignInExport),
			endpoint.Tags("签到导出")),

		endpoint.Post("/activity/export-user-data",
			"导出参加报名活动的用户",
			endpoint.Handler(ExportUserData),
			endpoint.Body(planTask.PlanTasksConditionsBaRequst{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: planTask.PlanTasksConditionsBaRespone{},
			}, "planTasksConditionsCompletion", "接口调用成功")),

		endpoint.Post("/activity/get-type-list",
			"获取所有的活动",
			endpoint.Handler(GetActivityList),
			endpoint.Body(planTask.PlanTasksConditionsBaRequst{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: planTask.PlanTasksConditionsBaRespone{},
			}, "planTasksConditionsCompletion", "接口调用成功")),

		endpoint.Post("/activity/get-user-list",
			"获取报名活动的所有用户",
			endpoint.Handler(GetUserList),
			endpoint.Body(planTask.PlanTasksConditionsBaRequst{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: planTask.PlanTasksConditionsBaRespone{},
			}, "planTasksConditionsCompletion", "接口调用成功")),

		endpoint.Post("/activity/get-user-info",
			"获取报名活动的详细信息",
			endpoint.Handler(GetUserInfo),
			endpoint.Body(planTask.PlanTasksConditionsBaRequst{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: planTask.PlanTasksConditionsBaRespone{},
			}, "planTasksConditionsCompletion", "接口调用成功")),

		endpoint.Post("/activity/verify-user",
			"报名活动用户审核通过",
			endpoint.Handler(VerifyUserInfo),
			endpoint.Body(planTask.PlanTasksConditionsBaRequst{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: planTask.PlanTasksConditionsBaRespone{},
			}, "planTasksConditionsCompletion", "接口调用成功")),
	}
}
