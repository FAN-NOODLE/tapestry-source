package http

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger/endpoint"
)

func Endpoints() []*swagger.Endpoint {
	handle := NewBaKanBan()

	return []*swagger.Endpoint{
		endpoint.Post("/ba_kanban", "BA业绩看板", endpoint.Handler(handle.GetBaKanBanData)),
		endpoint.Post("/ba_kanban/sales", "BA业绩看板", endpoint.Handler(handle.GetSalesDataDetail)),
		endpoint.Post("/ba_kanban/sales-export", "BA业绩销售看板", endpoint.Handler(handle.ExportSalesDataDetail)),
		endpoint.Post("/ba_kanban/operation", "BA业绩看板", endpoint.Handler(handle.GetOperationDataDetail)),
		endpoint.Post("/ba_kanban/operation-export", "BA业绩看板", endpoint.Handler(handle.ExportOperationDataDetail)),
		endpoint.Post("/ba_kanban/ba_source_depte", "获取BA看板部门", endpoint.Handler(handle.GetBADeptInfo)),
		endpoint.Post("/ba_kanban/v1/console/sales", "管理后台BA业绩详情V1", endpoint.Handler(handle.GetSalesDataDetailV1)),
		endpoint.Post("/ba_kanban/v1/operation", "BA业绩明细列表", endpoint.Handler(handle.BaSaleListOperation)),
		endpoint.Post("/ba_kanban/ba/sale_export", "BA业绩运营数据明细导出", endpoint.Handler(handle.ExportbaSale)),
		endpoint.Post("/ba_kanban/growth_rate", "获取同比增长信息", endpoint.Handler(handle.GrowthRate)),
	}
}
