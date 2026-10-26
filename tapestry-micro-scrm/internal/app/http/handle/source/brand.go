package source

import (
	"github.com/gin-gonic/gin"
	"micro-scrm/internal/app/http/handle"
	"micro-scrm/internal/app/http/request/source"
	sourceService "micro-scrm/internal/app/service/source"
	"micro-scrm/internal/pkg/helper"
)

type BrandHandle struct {
	handle.BaseHandle
	BrandSvc *sourceService.BrandService
}

func NewBrandHandle() BrandHandle {
	return BrandHandle{
		BrandSvc: sourceService.NewBrandService(),
	}
}

// SearchList 品牌素材库检索
func (h BrandHandle) SearchList(c *gin.Context) {
	var request source.BrandListReq
	err := c.ShouldBind(&request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	list, err := h.BrandSvc.SearchList(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, list)
	return
}

// SearchListGroupSku  品牌素材库检索获取skuItem
func (h BrandHandle) SearchListGroupSku(c *gin.Context) {
	var request source.BrandListReq
	err := c.ShouldBind(&request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	request.IncludeNoneImageProduct = true
	list, err := h.BrandSvc.SearchList(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	var skuItemArr []string
	for _, info := range list {
		if !helper.InArray(info.SkuItem, skuItemArr) {
			skuItemArr = append(skuItemArr, info.SkuItem)
		}
	}
	h.ResSuccess(c, skuItemArr)
	return
}

// CacheInfo    品牌素材库图片缓存查询
func (h BrandHandle) CacheInfo(c *gin.Context) {
	var request source.BrandCacheReq
	err := c.ShouldBind(&request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	url, err := h.BrandSvc.CacheInfo(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, url)
	return
}

func (h BrandHandle) CollectSearchList(c *gin.Context) {
	var request source.BrandListReq
	err := c.ShouldBind(&request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	list, err := h.BrandSvc.CollectSearchList(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, list)
	return
}

// TalkieImageSearch Talkie图片搜索
func (h BrandHandle) TalkieImageSearch(c *gin.Context) {
	var request source.TalkieImageSearchReq
	err := c.ShouldBind(&request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	list, err := h.BrandSvc.TalkieImageSearch(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, list)
	return
}
