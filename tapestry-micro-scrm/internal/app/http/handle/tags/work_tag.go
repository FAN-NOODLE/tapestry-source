package tags

import (
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/app/http/handle"
	tagsReq "micro-scrm/internal/app/http/request/tags"
	tagsSvc "micro-scrm/internal/app/service/tags"
)

type WorkTagHandle struct {
	handle.BaseHandle
	workTagSvc *tagsSvc.WorkTags
}

func NewWorkTagHandle() WorkTagHandle {
	return WorkTagHandle{
		workTagSvc: tagsSvc.NewWorkTags(),
	}
}

// List 客户分组的企微标签用户
func (h WorkTagHandle) List(context *gin.Context) {
	var req tagsReq.InfoReq
	if err := context.ShouldBind(&req); err != nil {
		h.ResError(context, err)
		return
	}
	var result tagsReq.TagListInfo
	userList, total, errList := h.workTagSvc.GetTagUserList(h.WarpContext(context), req)
	if errList != nil {
		h.ResError(context, errList)
		return
	}
	result.Total = total
	result.UserList = userList
	h.ResSuccess(context, result)
}

// UserTags 客户标签分组和标签列表
func (h WorkTagHandle) UserTags(context *gin.Context) {
	var req tagsReq.UserTagsReq
	if err := context.ShouldBind(&req); err != nil {
		h.ResError(context, err)
		return
	}

	userTagsGroupList, err := h.workTagSvc.UserTagsAndTagGroup(ctx.Wrap(context), req)
	if err != nil {
		h.ResError(context, err)
		return
	}
	h.ResSuccess(context, userTagsGroupList)

	return
}

// Auth 判断是否具有添加标签的权限
func (h WorkTagHandle) Auth(context *gin.Context) {
	var req tagsReq.InfoBaseReq
	if err := context.ShouldBind(&req); err != nil {
		h.ResError(context, err)
		return
	}
	type result struct {
		Auth bool `json:"auth"`
	}
	var res result
	res.Auth = true
	count, err := h.workTagSvc.TagsAuth(h.WarpContext(context), req)
	if err != nil {
		h.ResError(context, err)
		return
	}
	if count > 0 {
		res.Auth = false
	}
	h.ResSuccess(context, res)

	return
}

// GetList 企微标签列表
func (h WorkTagHandle) GetList(context *gin.Context) {
	var req tagsReq.ListReq
	if err := context.ShouldBind(&req); err != nil {
		h.ResError(context, err)
		return
	}
	list, err := h.workTagSvc.GetTagList(h.WarpContext(context), req.EnterpriseID)
	if err != nil {
		h.ResError(context, err)
		return
	}
	h.ResSuccess(context, list)

}
