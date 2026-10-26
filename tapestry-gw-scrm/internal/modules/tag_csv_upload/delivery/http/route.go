package http

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger/endpoint"
	"net/http"
)

func Endpoints() []*swagger.Endpoint {
	return []*swagger.Endpoint{
		endpoint.Post("/tag/csvFile/upload", "",
			endpoint.Handler(CSVFileUpload),
			endpoint.Tags("tag上传"),
			endpoint.Response(http.StatusOK, "", "upload", "接口调用成功")),
	}
}
