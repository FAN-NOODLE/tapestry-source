package http

import (
	"net/http"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger/endpoint"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"micro-scrm/internal/modules/ba_kanban"
)

//Endpoints 路由定义
func Endpoints() []*swagger.Endpoint {
	minniHandle := NewMiniHandle()
	consoleHandle := NewConsoleHandle()

	return []*swagger.Endpoint{
		// endpoint.Post("/ba_kanban/mini", "小程序BA看板",
		// 	endpoint.Handler(minniHandle.GetBaKanBanData),
		// 	endpoint.Tags("小程序BA看板"),
		// 	endpoint.Body(ba_kanban.MiniBaKanBanDataRequest{}, "请求参数", true),
		// 	endpoint.Response(http.StatusOK, ginplus.APIResponse{
		// 		Data: ba_kanban.BaKanBanDataResponse{},
		// 	}, "getPlanList", "接口调用成功")),
		endpoint.Post("/ba_kanban/mini", "小程序BA看板",
			endpoint.Handler(minniHandle.GetBaKanBanDatav1),
			endpoint.Tags("小程序BA看板"),
			endpoint.Body(ba_kanban.MiNiBaKanBanDataRequest{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: ba_kanban.BaKanBanDataResponse{},
			}, "getPlanList", "接口调用成功")),
		// endpoint.Post("/ba_kanban/console", "管理后台BA看板",
		// 	endpoint.Handler(consoleHandle.GetBaKanBanData),
		// 	endpoint.Tags("管理后台BA看板"),
		// 	endpoint.Body(ba_kanban.MiniBaKanBanDataRequest{}, "请求参数", true),
		// 	endpoint.Response(http.StatusOK, ginplus.APIResponse{
		// 		Data: ba_kanban.BaKanBanDataResponse{},
		// 	}, "getPlanList", "接口调用成功")),
		// endpoint.Post("/ba_kanban/console/sales", "管理后台BA业绩详情",
		// 	endpoint.Handler(consoleHandle.GetSalesDataDetail),
		// 	endpoint.Tags("管理后台BA业绩详情"),
		// 	endpoint.Body(ba_kanban.MiniBaKanBanDataRequest{}, "请求参数", true),
		// 	endpoint.Response(http.StatusOK, ginplus.APIResponse{
		// 		Data: ba_kanban.SalesDataDetailResponse{},
		// 	}, "getPlanList", "接口调用成功")),
		// endpoint.Post("/ba_kanban/console/operation", "管理后台BA运营详情",
		// 	endpoint.Handler(consoleHandle.GetOperationDataDetail),
		// 	endpoint.Tags("管理后台BA运营详情"),
		// 	endpoint.Body(ba_kanban.MiniBaKanBanDataRequest{}, "请求参数", true),
		// 	endpoint.Response(http.StatusOK, ginplus.APIResponse{
		// 		Data: ba_kanban.OperationDataDetailResponse{},
		// 	}, "getPlanList", "接口调用成功")),
		endpoint.Post("/ba-performance/list", "计划列表",
			endpoint.Handler(GetList),
			endpoint.Tags("ba指标"),
			endpoint.Body(ba_kanban.BaPerformanceQueryReq{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: ba_kanban.BaPerformanceQueryReq{},
			}, "getPlanList", "接口调用成功")),

		endpoint.Post("/ba-performance/create", "保存指标",
			endpoint.Handler(Create),
			endpoint.Tags("ba指标"),
			endpoint.Body(ba_kanban.BaPerformanceSavaReq{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: ba_kanban.BaPerformanceSavaReq{},
			}, "getPlanList", "接口调用成功")),
		endpoint.Post("/ba_kanban/ba_source_depte", "获取BA看板部门",
			endpoint.Handler(baSourceDepte),
			endpoint.Tags("ba看板部门"),
			endpoint.Body(ba_kanban.BaSourceDeptRequst{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: ba_kanban.BaPerformanceSavaReq{},
			}, "baSourceDepte", "接口调用成功")),
		endpoint.Post("/ba_kanban/v1/console/sales", "管理后台BA业绩详情V1",
			endpoint.Handler(consoleHandle.GetSalesDataDetailV1),
			endpoint.Tags("管理后台BA业绩详情"),
			endpoint.Body(ba_kanban.BaKanBanDataRequest{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: ba_kanban.SalesDataDetailResponse{},
			}, "GetSalesDataDetailV1", "接口调用成功")),
		endpoint.Post("/ba_kanban/growth_rate", "获取同比增长信息",
			endpoint.Handler(consoleHandle.GrowthRate),
			endpoint.Tags("管理后台BA运营详情"),
			endpoint.Body(ba_kanban.BaKanBanDataRequest{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: ba_kanban.SalesDataDetailResponse{},
			}, "getPlanList", "接口调用成功")),
		endpoint.Post("/ba_kanban/v1/console/operation", "ba业绩列表",
			endpoint.Handler(consoleHandle.GetOperationDataDetailv1),
			endpoint.Tags("ba业绩列表"),
			endpoint.Body(ba_kanban.BaKanBanDataRequest{}, "请求参数", true),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: ba_kanban.SalesDataDetailResponse{},
			}, "GetOperationDataDetailv1", "接口调用成功")),

		endpoint.Post("/ba_kanban/data/upload", "BA看板数据上传",
			endpoint.Handler(consoleHandle.baKanbandateupload),
			endpoint.Tags("ba业绩"),
			endpoint.Body(ba_kanban.BaKanDateUploadRequst{}, "请求参数", true),
		),
	}
}
