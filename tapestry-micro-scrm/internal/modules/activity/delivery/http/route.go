package http

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger/endpoint"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"micro-scrm/internal/modules/activity"
	activityMysql "micro-scrm/internal/modules/activity/respository/mysql"
	"micro-scrm/internal/modules/plan_task"
	"net/http"
)

// Endpoints 路由配置
func Endpoints() []*swagger.Endpoint {
	activityService = activity.NewActivity(activityMysql.New())
	activityHand := NewActivityHandler()
	return []*swagger.Endpoint{
		endpoint.Post("/activity/sync", "活动数据同步",
			endpoint.Handler(syncActivity),
			endpoint.Tags("活动"),
			endpoint.Response(http.StatusOK, "",
				"syncActivity", "接口调用成功")),

		endpoint.Post(
			"/activity/goods/list",
			"活动列表",
			endpoint.Handler(goodsList),
			endpoint.Tags("活动"),
			endpoint.Response(http.StatusOK, "", "list", "接口调用成功"),
		),

		endpoint.Post("/activity/get-type-list", "获取所有活动",
			endpoint.Handler(activityHand.GetTypeList),
			endpoint.Tags("报名活动"),
			endpoint.Body(plan_task.UserListReq{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: plan_task.UserListReq{},
			}, "planTasksInfo", "接口调用成功")),

		endpoint.Post("/activity/get-user-list", "获取所有活动参与人员",
			endpoint.Handler(activityHand.GetUserListByTpyeId),
			endpoint.Tags("报名活动"),
			endpoint.Body(plan_task.UserListReq{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: plan_task.UserListReq{},
			}, "planTasksInfo", "接口调用成功")),

		endpoint.Post("/activity/get-user-info", "获取活动人员详细信息",
			endpoint.Handler(activityHand.GetUserInfo),
			endpoint.Tags("报名活动"),
			endpoint.Body(plan_task.UserInfoReq{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: plan_task.UserInfoReq{},
			}, "planTasksInfo", "接口调用成功")),

		endpoint.Post("/activity/verify-user", "报名活动用户审核",
			endpoint.Handler(activityHand.VerifyUser),
			endpoint.Tags("报名活动"),
			endpoint.Body(plan_task.UserInfoReq{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: plan_task.UserInfoReq{},
			}, "planTasksInfo", "接口调用成功")),

		endpoint.Post("/activity/export-user-data", "导出报名活动用户",
			endpoint.Handler(activityHand.ExportUserData),
			endpoint.Tags("报名活动"),
			endpoint.Body(plan_task.ExportUsreReq{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: plan_task.ExportUsreReq{},
			}, "planTasksInfo", "接口调用成功")),

		endpoint.Post("/mp/activity/insert-user-data", "插入报名活动用户",
			endpoint.Handler(activityHand.InsertUserData),
			endpoint.Tags("报名活动"),
			endpoint.Body(plan_task.InsertUserReq{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: plan_task.InsertUserReq{},
			}, "planTasksInfo", "接口调用成功")),

		endpoint.Post("/mp/activity/get-activity-id", "获取报名活动id",
			endpoint.Handler(activityHand.GetActivityId),
			endpoint.Tags("报名活动"),
			endpoint.Body(plan_task.GetActivityIdReq{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: plan_task.GetActivityIdReq{},
			}, "planTasksInfo", "接口调用成功")),

		endpoint.Post("/mp/activity/check-activity-status", "检查活动状态",
			endpoint.Handler(activityHand.CheckActivityStatus),
			endpoint.Tags("报名活动"),
			endpoint.Body(plan_task.CheckStatusReq{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: plan_task.CheckStatusReq{},
			}, "planTasksInfo", "接口调用成功")),
	}
}
