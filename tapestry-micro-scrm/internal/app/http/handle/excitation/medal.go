package excitation

import (
	"github.com/gin-gonic/gin"
	"micro-scrm/internal/app/http/handle"
	excitationReq "micro-scrm/internal/app/http/request/excitation"
	"micro-scrm/internal/app/service/excitation"
)

type MedalHandle struct {
	handle.BaseHandle
	medalSvc *excitationSvc.MedalService
}

func NewMedalHandle() MedalHandle {
	return MedalHandle{
		medalSvc: excitationSvc.NewMedalService(),
	}
}

func (h MedalHandle) List(c *gin.Context) {
	var request excitationReq.MedalList
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	resp, total, err := h.medalSvc.GetMedalList(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResList(c, total, request.Page, request.Size, resp)

	return
}

func (h MedalHandle) Create(c *gin.Context) {
	var request excitationReq.MedalCreate
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	err := h.medalSvc.Create(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, "")
	return
}

func (h MedalHandle) Detail(c *gin.Context) {
	var request excitationReq.MedalDetail
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	resp, err := h.medalSvc.Detail(h.WarpContext(c), request.Id)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, resp)
	return
}

func (h MedalHandle) Update(c *gin.Context) {
	var request excitationReq.MedalUpdate
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	err := h.medalSvc.Update(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, "")
	return
}

func (h MedalHandle) Remove(c *gin.Context) {
	var request excitationReq.MedalDetail
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	err := h.medalSvc.RemoveRule(h.WarpContext(c), request.Id)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, "")
	return
}

func (h MedalHandle) Status(c *gin.Context) {
	var request excitationReq.MedalStatus
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	err := h.medalSvc.ChangeStatus(h.WarpContext(c), request.Id, request.Status)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, "")
	return
}

// MedalWall 奖章墙
func (h MedalHandle) MedalWall(c *gin.Context) {
	var request excitationReq.MedalWall
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	result, err := h.medalSvc.MedalWall(h.WarpContext(c), request.EnterpriseId, request.Mobile)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, result)
	return
}

// MedalWallDetail 奖章墙详情
func (h MedalHandle) MedalWallDetail(c *gin.Context) {
	var request excitationReq.MedalWallDetail
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	result, err := h.medalSvc.MedalWallDetail(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, result)
	return
}
