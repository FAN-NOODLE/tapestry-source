package download_http

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger/endpoint"
	"micro-download/internal/modules/download"
	"net/http"
)

func Endpoints() []*swagger.Endpoint {
	service = &download.Service{}

	return []*swagger.Endpoint{
		endpoint.Post("/v1/export-list", "下载中心列表",
			endpoint.Handler(exportList),
			endpoint.Response(http.StatusOK, "", "download", "接口调用成功")),
		endpoint.Post("/v1/delete", "下载中心删除",
			endpoint.Handler(deleteFile),
			endpoint.Query("id", "integer", "下载中心主键ID", true),
			endpoint.Response(http.StatusOK, "",
				"download", "接口调用成功")),
		endpoint.Post("/v1/create-export", "下载中心添加",
			endpoint.Handler(createExport),
			endpoint.Response(http.StatusOK, "",
				"download", "接口调用成功")),
		endpoint.Post("/v1/retry", "重新下载",
			endpoint.Handler(retry),
			endpoint.Path("id", "integer", "下载中心主键ID", true),
			endpoint.Response(http.StatusOK, "", "download", "接口调用成功")),
		endpoint.Post("/v1/process", "任务进度",
			endpoint.Handler(process),
			endpoint.Path("id", "integer", "下载中心主键ID", true),
			endpoint.Response(http.StatusOK, "", "download", "接口调用成功")),
		endpoint.Get("/v1/download", "下载",
			endpoint.Handler(downLoad),
			endpoint.Path("id", "integer", "下载中心主键ID", true),
			endpoint.Response(http.StatusOK, "", "download", "接口调用成功")),
		endpoint.Post("/v1/detail", "任务详情",
			endpoint.Handler(detail),
			endpoint.Path("id", "integer", "下载中心主键ID", true),
			endpoint.Response(http.StatusOK, "", "download-detail", "接口调用成功")),
		endpoint.Post("/v1/latest-export", "指定业务的最近一次导出",
			endpoint.Handler(latestExport),
			endpoint.Path("id", "integer", "下载中心主键ID", true),
			endpoint.Response(http.StatusOK, "", "download", "接口调用成功")),
		endpoint.Post("/v1/auto-detail", "自动下载资源",
			endpoint.Handler(autoDetail),
			endpoint.Response(http.StatusOK, "", "download", "接口调用成功")),
	}
}
