package microscrm

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/pkg/sdk"
)

type SourceApi struct{}

func NewSourceApi() SourceApi {
	return SourceApi{}
}

func (s SourceApi) TextList(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/source/text-list"), request)
}

func (s SourceApi) ImageList(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/source/image-list"), request)
}

func (s SourceApi) LinkList(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/source/link-list"), request)
}

func (s SourceApi) MiniProgramList(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/source/mini-program-list"), request)
}

func (s SourceApi) VideoList(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/source/video-list"), request)
}

func (s SourceApi) FileList(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/source/file-list"), request)
}

func (s SourceApi) TextDetail(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/source/text/detail"), request)
}

func (s SourceApi) ImageDetail(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/source/image/detail"), request)
}

func (s SourceApi) LinkDetail(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/source/link/detail"), request)
}

func (s SourceApi) MiniProgramDetail(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/source/mini_program/detail"), request)
}

func (s SourceApi) VideoDetail(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/source/video/detail"), request)
}

func (s SourceApi) FileUpdate(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/source/file/update"), request)
}

func (s SourceApi) FriendUpdate(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/source/friend/update"), request)
}

func (s SourceApi) TextUpdate(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/source/text/update"), request)
}

func (s SourceApi) ImageUpdate(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/source/image/update"), request)
}

func (s SourceApi) LinkUpdate(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/source/link/update"), request)
}

func (s SourceApi) MiniProgramUpdate(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/source/mini_program/update"), request)
}

func (s SourceApi) VideoUpdate(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/source/video/update"), request)
}

func (s SourceApi) FileDetail(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/source/file/detail"), request)
}

func (s SourceApi) FriendDetail(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/source/friend/detail"), request)
}

func (s SourceApi) RemoveCategory(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/source/remove-category"), request)
}

func (s SourceApi) UpdateStatus(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/source/updateStatus"), request)
}

func (s SourceApi) ConsoleList(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/source/console-list"), request)
}

func (s SourceApi) GetMediaInfo(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/source/get-media-info"), request)
}

func (s SourceApi) CompressFile(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/source/compress-file"), request)
}

func (s SourceApi) SetTop(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/source/set-top"), request)
}
