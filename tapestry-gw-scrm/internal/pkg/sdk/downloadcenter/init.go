package downloadcenter

import (
	"fmt"
	"sync"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/pkg/sdk"
)

type downloader interface {
	CreateExport(ctx.Context, CreateExportRequest) (CreateExportResponse, error)
	ExportList(ctx.Context, ExportListRequest) (ExportListResponse, error)
	Retry(ctx.Context, RetryRequest) (RetryResponse, error)
	Process(ctx.Context, ProcessRequest) (ProcessResponse, error)
	Download(ctx.Context, DownloadRequest) (DownloadResponse, error)
	Delete(ctx.Context, DeleteRequest) (DeleteResponse, error)
	ExportDetail(ctx.Context, DetailRequest) (ExportDetailResponse, error)
	LatestExportDetail(ctx.Context, LatestExportDetailRequest) (ExportDetailResponse, error)
	GetAutoDetail(ctx.Context, OperationReq) (ExportDetailResponse, error)
}

var DC downloader

type downloadCenter struct {
	http    *sdk.Client
	baseUri string `desc:"请求链接的公共path"`
}

var once sync.Once

// 获取下载中心实例
// baseUri string 请求链接的公共url path
func newDownload(baseUri string) downloader {
	once.Do(func() {
		DC = &downloadCenter{
			http:    sdk.NewHttpClient(sdk.ParseJson{}),
			baseUri: baseUri,
		}
	})

	return DC
}

// 下载中心对象初始化
// baseUri baseUri string 请求链接的公共url path
func Init(baseUri string) error {
	_ = newDownload(baseUri)

	return nil
}

func (dc downloadCenter) makeUrl(path string) string {
	return fmt.Sprintf("%s%s", dc.baseUri, path)
}
