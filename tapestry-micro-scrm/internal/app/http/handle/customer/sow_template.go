package customer

import (
	"errors"
	"github.com/gin-gonic/gin"
	"micro-scrm/internal/app/http/handle"
	customerReq "micro-scrm/internal/app/http/request/customer"
	customerModels "micro-scrm/internal/app/models/customer"
	customerSvc "micro-scrm/internal/app/service/customer"
)

type SowTemplateHandle struct {
	handle.BaseHandle
	SowTemplateSvc *customerSvc.SowTemplateService
	SowTaskSvc     *customerSvc.SowTask
}

func NewSowTemplateHandle() SowTemplateHandle {
	return SowTemplateHandle{
		SowTemplateSvc: customerSvc.NewSowTemplateService(),
		SowTaskSvc:     customerSvc.NewSowTask(),
	}
}

func (h SowTemplateHandle) List(c *gin.Context) {
	var request customerReq.SowTemplateListReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	result, total, err := h.SowTemplateSvc.List(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResList(c, total, request.Page, request.Size, result)
	return
}

func (h SowTemplateHandle) Save(c *gin.Context) {
	var request customerReq.SowTemplateSaveReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	err := h.SowTemplateSvc.Save(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)
	return
}

func (h SowTemplateHandle) Detail(c *gin.Context) {
	var request customerReq.SowTemplateDetailReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	detail, err := h.SowTemplateSvc.Detail(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, detail)
	return
}

func (h SowTemplateHandle) AllTaskList(c *gin.Context) {
	var request customerReq.SowAllTaskListReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	result, total, err := h.SowTaskSvc.AllTaskList(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResList(c, total, request.Page, request.Size, result)
	return
}

func (h SowTemplateHandle) UpdateStatus(c *gin.Context) {
	var request customerReq.SowTemplateDetailReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	if request.Status != customerModels.SowTemplateIsUp && request.Status != customerModels.SowTemplateIsDown {
		h.ResError(c, errors.New("参数错误"))
		return
	}
	err := h.SowTemplateSvc.UpdateStatus(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)
	return
}
func (h SowTemplateHandle) Delete(c *gin.Context) {
	var request customerReq.SowTemplateDetailReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	err := h.SowTemplateSvc.Delete(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)
	return
}
