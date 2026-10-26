package downloadcenter

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
)

func (dc downloadCenter) ExportDetail(ctx ctx.Context, request DetailRequest) (response ExportDetailResponse, err error) {

	err = dc.http.Post(ctx, dc.makeUrl(UrlDetail), request, &response)
	return
}
