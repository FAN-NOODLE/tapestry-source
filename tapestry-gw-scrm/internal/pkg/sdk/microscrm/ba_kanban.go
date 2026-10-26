package microscrm

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/pkg/sdk"
)

type ApiBaKanBan struct{}

const (
	UrlBaKanBanData        = "%s/micro-scrm/v1/ba_kanban/console"
	UrlSalesDataDetail     = "%s/micro-scrm/v1/ba_kanban/console/sales"
	UrlOperationDataDetail = "%s/micro-scrm/v1/ba_kanban/console/operation"
	UrlGetBADeptInfo       = "%s/micro-scrm/v1/ba_kanban/ba_source_depte"
	UrlGetKanBanInfo       = "%s/micro-scrm/v1/ba_kanban/v1/console/sales"
	URLBASALEDetailV1      = "%s/micro-scrm/v1/ba_kanban/v1/console/operation"
	URLGrowthRate          = "%s/micro-scrm/v1/ba_kanban/growth_rate"

	BaKanBanDataUploadURL = "%s/micro-scrm/v1/ba_kanban/data/upload"
)

// GetBaKanBanData 获取BA业绩看板数据
func (api ApiBaKanBan) GetBaKanBanData(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(UrlBaKanBanData), request)
}

// GetSalesDataDetail 获取BA业绩看板数据
func (api ApiBaKanBan) GetSalesDataDetail(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(URLBASALEDetailV1), request)
}

// GetOperationDataDetail 获取BA业绩看板数据
func (api ApiBaKanBan) GetOperationDataDetail(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(UrlOperationDataDetail), request)
}

func (api ApiBaKanBan) GetBAInfo(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(UrlGetBADeptInfo), request)
}

func (api ApiBaKanBan) GetKanBanInfo(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(UrlGetKanBanInfo), request)
}

func (api ApiBaKanBan) BaSaleList(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(URLBASALEDetailV1), request)
}

func (api ApiBaKanBan) GrowthRate(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(URLGrowthRate), request)
}
func (api ApiBaKanBan) GetSalesDataDetailV1(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(URLBASALEDetailV1), request)
}

func (api ApiBaKanBan) BaKanBanDataUpload(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(BaKanBanDataUploadURL), request)
}
