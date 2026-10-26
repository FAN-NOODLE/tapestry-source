package downloadcenter

import (
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/config"
	"strings"
)

func (dc downloadCenter) GetAutoDetail(ctx ctx.Context, request OperationReq) (response ExportDetailResponse, err error) {

	err = dc.http.Post(ctx, dc.makeUrl(UrlAutoDetail), request, &response)
	if response.Data.Id > 0 {
		response.Data.Path = formatDownloadCenterPath(response.Data.Path)
	}
	return
}

func formatDownloadCenterPath(path string) string {
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		return path
	}
	return fmt.Sprintf("%s://%s/%s", config.CfgData.Storage.COS.Schema, config.CfgData.Storage.COS.CdnHost, path)
}
