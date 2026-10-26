package api

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger"
	download_http "micro-download/internal/modules/download/delivery/http"
)

func endpoints() []*swagger.Endpoint {
	var eps []*swagger.Endpoint

	// 下载中心
	eps = append(eps, download_http.Endpoints()...)
	return eps
}
