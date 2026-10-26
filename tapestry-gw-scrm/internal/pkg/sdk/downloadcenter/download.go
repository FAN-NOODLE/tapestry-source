package downloadcenter

import (
	"fmt"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
)

func (dc downloadCenter) Download(ctx ctx.Context, request DownloadRequest) (response DownloadResponse, err error) {
	err = dc.http.Get(ctx, dc.makeUrl(fmt.Sprintf("%s?id=%d&type=%d", UrlDownload, request.FileId, request.DownloadType)), request, &response)
	return
}
