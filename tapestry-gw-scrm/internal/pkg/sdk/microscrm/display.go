package microscrm

import (
	"gw-scrm/internal/pkg/sdk"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
)

type DisplayApi struct {
}

func NewDisplayApi() DisplayApi {
	return DisplayApi{}
}

func (a DisplayApi) List(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/display/task/list"), request)
}

func (a DisplayApi) Create(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/display/task/add"), request)
}

func (a DisplayApi) Detail(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/display/task/detail"), request)
}

func (a DisplayApi) Update(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/display/task/update"), request)
}

func (a DisplayApi) Status(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/display/task/status"), request)
}

func (a DisplayApi) Copy(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/display/task/copy"), request)
}

func (a DisplayApi) AllotList(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/display/task/allot-list"), request)
}

func (a DisplayApi) AllotNum(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/display/task/allot-num"), request)
}

func (a DisplayApi) NoticeSetting(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/display/task/notice-setting"), request)
}

func (a DisplayApi) AllotInfo(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/display/task/allot-info"), request)
}

func (a DisplayApi) AllotCheck(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/display/task/allot-check"), request)
}

func (a DisplayApi) NoticeSettingInfo(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/display/task/notice-setting-info"), request)
}

func (a DisplayApi) AllotSendNotice(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/display/task/send-notice"), request)
}

func (a DisplayApi) GetImageList(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/display/task/image-list"), request)
}

func (a DisplayApi) AllotNoticeNum(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/display/task/notice-num"), request)
}

func (a DisplayApi) VmList(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/display/task/vm-list"), request)
}

func (a DisplayApi) CosAvatar(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/display/task/cos-avatar"), request)
}

func (a DisplayApi) ReviewFileCheck(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/display/task/review-file-check"), request)
}

func (a DisplayApi) ReviewExport(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/display/task/review-export"), request)
}
