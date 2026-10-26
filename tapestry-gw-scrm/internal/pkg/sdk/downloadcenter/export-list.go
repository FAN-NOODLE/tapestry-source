package downloadcenter

import "gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"

func (dc downloadCenter) ExportList(ctx ctx.Context, request ExportListRequest) (response ExportListResponse, err error) {
	err = dc.http.Post(ctx, dc.makeUrl(UrlExportList), request, &response)
	for k := range response.Data.List {
		if response.Data.List[k].Status == 3 {
			response.Data.List[k].Path = formatDownloadCenterPath(response.Data.List[k].Path)
		}
	}
	return
}
