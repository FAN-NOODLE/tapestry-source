package downloadcenter

import "gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"

func (dc downloadCenter) LatestExportDetail(ctx ctx.Context, request LatestExportDetailRequest) (response ExportDetailResponse, err error) {

	err = dc.http.Post(ctx, dc.makeUrl(UrlLatestExportDetail), request, &response)
	return
}
