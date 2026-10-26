package downloadcenter

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
)

func (dc downloadCenter) CreateExport(ctx ctx.Context, request CreateExportRequest) (response CreateExportResponse, err error) {

	err = dc.http.Post(ctx, dc.makeUrl(UrlCreateExport), request, &response)
	return
}
