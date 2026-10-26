package microscrm

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/pkg/sdk"
)

type ApiSource struct{}

const (
	MoveSource        = "%s/micro-scrm/v1/source-group/move"
	RemoveSource      = "%s/micro-scrm/v1/source-group/delete"
	SourceTree        = "%s/micro-scrm/v1/source-group/list"
	AddSource         = "%s/micro-scrm/v1/source-group/add"
	GetGroupName      = "%s/micro-scrm/v1/source/get/group_name"
	UpdateIgnoreIndex = "%s/micro-scrm/v1/source-group/update-ignore-index"
)

func (s ApiSource) AddSource(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(AddSource), request)
}

func (s ApiSource) MoveBatch(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(MoveSource), request)
}

func (s ApiSource) SourceDelete(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(RemoveSource), request)
}

func (s ApiSource) SourceTree(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(SourceTree), request)
}

func (s ApiSource) GetGroupName(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(GetGroupName), request)
}

func (s ApiSource) UpdateIgnoreIndex(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(UpdateIgnoreIndex), request)
}
