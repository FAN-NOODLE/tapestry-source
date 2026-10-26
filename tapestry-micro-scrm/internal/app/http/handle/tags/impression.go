package tags

import (
	"github.com/gin-gonic/gin"
	"micro-scrm/internal/app/http/handle"
	tagsReq "micro-scrm/internal/app/http/request/tags"
	tagsSvc "micro-scrm/internal/app/service/tags"
)

type ImpressionHandle struct {
	handle.BaseHandle
	impressionTagSvc *tagsSvc.Impression
}

func NewImpressionHandle() ImpressionHandle {
	return ImpressionHandle{
		impressionTagSvc: tagsSvc.NewImpression(),
	}
}

func (h ImpressionHandle) GroupList(c *gin.Context) {
	var request tagsReq.GroupSearchList
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	result, total, err := h.impressionTagSvc.GroupSearchList(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResList(c, total, request.Page, request.Size, result)
	return
}

func (h ImpressionHandle) TagList(c *gin.Context) {
	var request tagsReq.TagList
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	result, total, err := h.impressionTagSvc.TagSearchList(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResList(c, total, request.Page, request.Size, result)
	return
}

func (h ImpressionHandle) TagCreate(c *gin.Context) {
	h.ResSuccess(c, c.Request.Body)
	return
}

func (h ImpressionHandle) TagUpdate(c *gin.Context) {
	h.ResSuccess(c, c.Request.Body)
	return
}

func (h ImpressionHandle) GroupCreate(c *gin.Context) {
	h.ResSuccess(c, c.Request.Body)
	return
}

// UserTagsList 客户印象标签分组和标签列表
func (h ImpressionHandle) UserTagsList(context *gin.Context) {
	var req tagsReq.UserTagsReq
	if err := context.ShouldBind(&req); err != nil {
		h.ResError(context, err)
		return
	}

	userTagsGroupList, err := h.impressionTagSvc.TagsAndTagGroupList(h.WarpContext(context), req)
	if err != nil {
		h.ResError(context, err)
		return
	}
	h.ResSuccess(context, userTagsGroupList)

	return
}

// GetList 印象标签列表
func (h ImpressionHandle) GetList(c *gin.Context) {
	var request tagsReq.ReqList
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	result, err := h.impressionTagSvc.GetList(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, result)
	return
}
