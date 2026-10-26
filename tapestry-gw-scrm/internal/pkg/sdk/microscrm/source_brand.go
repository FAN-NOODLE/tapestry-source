package microscrm

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/pkg/sdk"
)

type SourceBrandApi struct {
}

func NewSourceBrandApi() SourceBrandApi {
	return SourceBrandApi{}
}

// BrandSearchList 品牌搜索列表
func (a SourceBrandApi) BrandSearchList(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/source/brand/search-list"), request)
}

// BrandLinkSave 品牌搜索列表
func (a SourceBrandApi) BrandLinkSave(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/source/link-save"), request)
}

// TalkieImageSearch Ai图片搜索
func (a SourceBrandApi) TalkieImageSearch(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl("%s/micro-scrm/v1/source/brand/talkie-image-search"), request)
}
