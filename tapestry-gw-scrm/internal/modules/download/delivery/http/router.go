package http

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger/endpoint"
	"net/http"
)

func Endpoints() []*swagger.Endpoint {
	download := newDownloadHandle()

	return []*swagger.Endpoint{
		endpoint.Post("/download-center/export-list", "下载中心列表",
			endpoint.Handler(download.exportList),
			endpoint.Response(http.StatusOK, "", "export-list", "接口调用成功")),
		endpoint.Get("/download-center/download", "下载文件",
			endpoint.Handler(download.downloadFile),
			endpoint.Response(http.StatusOK, "", "download", "接口调用成功")),
		endpoint.Post("/download-center/delete", "删除文件",
			endpoint.Handler(download.deleteFile),
			endpoint.Response(http.StatusOK, "", "download", "接口调用成功")),
		endpoint.Post("/download-center/detail", "查询详情",
			endpoint.Handler(download.detail),
			endpoint.Response(http.StatusOK, "", "download", "接口调用成功")),
		endpoint.Post("/download-center/retry", "下载重试",
			endpoint.Handler(download.retry),
			endpoint.Response(http.StatusOK, "", "download", "接口调用成功")),
		endpoint.Post("/download-center/latest-export-detail", "查询最近一次导出详情",
			endpoint.Handler(download.latestExportDetail),
			endpoint.Response(http.StatusOK, "", "download", "接口调用成功")),
		endpoint.Post("/download-center/get-auto-detail", "查询自动下载的资源",
			endpoint.Handler(download.getAutoDetail),
			endpoint.Response(http.StatusOK, "", "download", "接口调用成功")),
		endpoint.Post("/download-center/white-list", "白名单列表",
			endpoint.Handler(download.whiteList),
			endpoint.Response(http.StatusOK, "", "download", "接口调用成功")),
	}
}
