package customer

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/app/http/handle"
	requestBase "micro-scrm/internal/app/http/request"
	customerReq "micro-scrm/internal/app/http/request/customer"
	excitationReq "micro-scrm/internal/app/http/request/excitation"
	excitationModels "micro-scrm/internal/app/models/excitation"
	customerSvc "micro-scrm/internal/app/service/customer"
	excitationSvc "micro-scrm/internal/app/service/excitation"
	"strconv"
)

type SowTaskHandle struct {
	handle.BaseHandle
	sowTaskSvc *customerSvc.SowTask
}

func NewSowTaskHandle() SowTaskHandle {
	return SowTaskHandle{
		sowTaskSvc: customerSvc.NewSowTask(),
	}
}

func (h SowTaskHandle) List(c *gin.Context) {
	var request customerReq.SowTaskList
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	result, total, err := h.sowTaskSvc.List(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResList(c, total, request.Page, request.Size, result)
	return
}

func (h SowTaskHandle) Create(c *gin.Context) {
	var request customerReq.SowTaskCreate
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	id, err := h.sowTaskSvc.Create(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, id)
	return
}

func (h SowTaskHandle) Update(c *gin.Context) {
	var request customerReq.SowTaskUpdate
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	err := h.sowTaskSvc.Update(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)
	return
}

func (h SowTaskHandle) Finished(c *gin.Context) {
	var request customerReq.SowTaskFinished
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	err := h.sowTaskSvc.Finish(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	// 积分队列埋点
	excitationSvc.NewPointsProduceService().ProduceQueue(ctx.Wrap(c), excitationReq.PointsProduceSave{
		Mobile:       request.Mobile,
		Type:         excitationModels.RuleCategoryId14,
		Value:        1,
		UniqueHash:   strconv.FormatInt(request.Id, 10) + fmt.Sprintf("-%s-%d", request.Mobile, excitationModels.RuleCategoryId14),
		EnterpriseId: request.EnterpriseID,
	})

	h.ResOk(c)
	return
}

func (h SowTaskHandle) Delete(c *gin.Context) {
	var request customerReq.SowTaskFinished
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	err := h.sowTaskSvc.Delete(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResOk(c)
	return
}

func (h SowTaskHandle) TemplateList(c *gin.Context) {
	var request requestBase.EnterpriseBaseReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	templateList, err := h.sowTaskSvc.TemplateList(h.WarpContext(c), request.EnterpriseId)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, templateList)
	return
}

// NewCreate Coach的新播种模板
func (h SowTaskHandle) NewCreate(c *gin.Context) {
	var request customerReq.SowTaskNewCreate
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	sowTaskSimple, err := h.sowTaskSvc.NewCreate(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, sowTaskSimple)
}
