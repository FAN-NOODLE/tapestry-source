package source

import (
	"github.com/gin-gonic/gin"
	"micro-scrm/internal/app/http/handle"
	sourceRequest "micro-scrm/internal/app/http/request/source"
	memberSvc "micro-scrm/internal/app/service/member"
	sourceService "micro-scrm/internal/app/service/source"
)

type CategoryHandle struct {
	handle.BaseHandle
	categoryService sourceService.CategoryService
}

func NewCategoryHandle() CategoryHandle {
	return CategoryHandle{
		categoryService: sourceService.NewCategoryService(),
	}
}

func (h CategoryHandle) List(c *gin.Context) {
	var request sourceRequest.CategoryListRequest
	err := c.ShouldBind(&request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	if request.IsCollect > 0 {
		request.ThirdIds, err = memberSvc.NewMembersService().MembersDepartmentsDao.GetThirdsByMemberIds(h.WarpContext(c), request.EnterpriseID, request.MemberId)
		if err != nil {
			h.ResError(c, err)
			return
		}
	}
	err, result, total := h.categoryService.GetList(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResList(c, total, request.Page, request.Size, result)
	return
}

func (h CategoryHandle) Create(c *gin.Context) {
	var request sourceRequest.CategoryCreateRequest
	err := c.ShouldBind(&request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	err = h.categoryService.Create(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)
	return
}

func (h CategoryHandle) Update(c *gin.Context) {
	var request sourceRequest.CategoryUpdateRequest
	err := c.ShouldBind(&request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	err = h.categoryService.Update(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)
	return
}

func (h CategoryHandle) SourceRelation(c *gin.Context) {
	var request sourceRequest.CategoryRequest
	err := c.ShouldBind(&request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	ids, err := h.categoryService.SourceRelation(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, ids)
	return
}

func (h CategoryHandle) Delete(c *gin.Context) {
	var request sourceRequest.CategoryRequest
	err := c.ShouldBind(&request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	err = h.categoryService.Delete(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)
	return
}

func (h CategoryHandle) RemoveCategory(c *gin.Context) {
	var request sourceRequest.RemoveCategoryRequest
	err := c.ShouldBind(&request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	err = h.categoryService.RemoveCategory(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)
	return
}
