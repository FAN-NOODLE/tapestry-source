package file_http

import (
	"net/http"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger/endpoint"
)

func Endpoints() []*swagger.Endpoint {
	return []*swagger.Endpoint{
		endpoint.Post("/file/upload", "文件上传",
			endpoint.Handler(upload),
			endpoint.Tags("文件上传"),
			endpoint.Response(http.StatusOK, "", "file-upload", "接口调用成功")),

		endpoint.Post("/file/credential", "生成临时凭证",
			endpoint.Handler(credential),
			endpoint.Tags("生成临时凭证"),
			endpoint.Response(http.StatusOK, "", "file-credential", "接口调用成功")),

		endpoint.Post("/file/private-by-key", "下载私有cos文件",
			endpoint.Handler(privateByKey),
			endpoint.Tags("下载私有cos文件"),
			endpoint.Response(http.StatusOK, "", "file-credential", "接口调用成功")),

		endpoint.Post("/file/private-url", "获取私有cos文件临时下载链接",
			endpoint.Handler(privateURL),
			endpoint.Tags("获取私有cos文件临时下载链接"),
			endpoint.Response(http.StatusOK, "", "file-private-url", "接口调用成功")),

		endpoint.Post("/file/xlsx/ba_data_upload", "BA业绩看板数据上传",
			endpoint.Handler(badatassupload),
			endpoint.Tags("BA业绩看板数据上传"),
			endpoint.Response(http.StatusOK, "", "file-credential", "接口调用成功")),
	}
}
